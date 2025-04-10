package model

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/Proluxe/proluxe-common-api/services"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
	"github.com/scottraio/simpleforce"
)

type File struct {
	Client     *storage.Client
	Context    context.Context
	SF         *simpleforce.Client
	GinContext *gin.Context
}

type SharedFile struct {
	Object     string `json:"Object"`
	ObjectId   string `json:"ObjectId"`
	ObjectName string `json:"ObjectName"`
	Path       string `json:"Path"`
}

type SharedItem struct {
	Object     string `json:"Object"`
	ObjectId   string `json:"ObjectId"`
	ObjectName string `json:"ObjectName"`
}

func New(c *gin.Context, App *util.App) *File {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Printf("Failed to create client: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to storage client"})
		return nil
	}

	defer client.Close()

	return &File{
		Client:     client,
		Context:    ctx,
		SF:         App.SF.Client,
		GinContext: c,
	}
}

// SF Functions

func FetchFiles(client *simpleforce.Client, whereCondition string) []SharedFile {
	q := fmt.Sprintf(`
		SELECT Object__c, Object_Id__c, Path__c, Object_Name__c
		FROM CXP_File__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var files []SharedFile
	for _, record := range result.Records {
		f := SharedFile{
			Object:     getStringField("Object__c", record),
			ObjectId:   getStringField("Object_Id__c", record),
			ObjectName: getStringField("Object_Name__c", record),
			Path:       getStringField("Path__c", record),
		}

		files = append(files, f)
	}

	return files
}

// Instance Functions

func (f *File) DeleteFile(bucketName, objectName string) error {
	o := f.Client.Bucket(bucketName).Object(objectName)
	if err := o.Delete(f.Context); err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	err := f.deleteSharedFile(objectName)
	if err != nil {
		log.Printf("Error Deleting from Salesforce: %v\n", err)
		f.GinContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete files from Salesforce"})
	}

	return nil
}

func (f *File) deleteSharedFile(path string) error {
	// Delete the corresponding files from Salesforce
	q := fmt.Sprintf(`
		SELECT Id
		FROM CXP_File__c
		WHERE Path__c = '%s'
	`, path)

	result, err := f.SF.Query(q)
	if err != nil {
		return err
	}

	for _, record := range result.Records {
		if err := f.SF.SObject("CXP_File__c").Set("Id", record["Id"].(string)).Delete(); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions
func (f *File) ListBucketContents(bucketName, path string) ([]map[string]interface{}, error) {
	it := f.Client.Bucket(bucketName).Objects(f.Context, &storage.Query{Prefix: path, Delimiter: "/"})

	var contents []map[string]interface{}
	dirMap := make(map[string]struct{}) // To keep track of unique directories

	for {
		objAttrs, err := it.Next()

		var isPublic bool
		if objAttrs != nil {
			for _, acl := range objAttrs.ACL {
				if acl.Entity == "allUsers" {
					isPublic = true
					break
				}
			}
		}

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to iterate bucket contents: %v", err)
		}

		// Add files to the list
		if objAttrs.Name != "" && objAttrs.Name != path {
			content := map[string]interface{}{
				"name":     objAttrs.Name,
				"size":     objAttrs.Size,
				"updated":  objAttrs.Updated,
				"isPublic": isPublic,
			}
			contents = append(contents, content)
		}

		// Capture directories by checking for prefixes within object names
		if len(objAttrs.Prefix) > 0 {
			if _, exists := dirMap[objAttrs.Prefix]; !exists {
				dirMap[objAttrs.Prefix] = struct{}{}
				contents = append(contents, map[string]interface{}{
					"name":    objAttrs.Prefix,
					"size":    nil, // Directories don't have a size
					"updated": nil, // Directories don't have an updated timestamp
				})
			}
		}
	}

	return contents, nil
}

func (f *File) CreateFolder(bucketName, folderName string) error {
	wc := f.Client.Bucket(bucketName).Object(folderName).NewWriter(f.Context)
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	return nil
}

func (f *File) FetchFilesForPath(folderPath string) []SharedFile {
	escapedFolderPath := url.PathEscape(folderPath)

	q := fmt.Sprintf(`
		SELECT Object__c, Object_Id__c, Path__c, Object_Name__c
		FROM CXP_File__c
		WHERE Path__c LIKE '%%%s%%'
	`, escapedFolderPath)

	result, err := f.SF.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var files []SharedFile
	for _, record := range result.Records {
		f := SharedFile{
			Object:     getStringField("Object__c", record),
			ObjectId:   getStringField("Object_Id__c", record),
			ObjectName: getStringField("Object_Name__c", record),
			Path:       getStringField("Path__c", record),
		}

		files = append(files, f)
	}

	return files
}

func (f *File) SendFileShareConfirmationEmail(from, path string, sharedItems []SharedItem) error {
	fileName := filepath.Base(path)
	folderName := filepath.Dir(path)

	knock := services.Knock{
		WorkFlowId: "attach-file",
		Email:      from,
	}

	knock.Identify()
	recipients := []string{"cxp@proluxe.com"}

	return knock.Trigger(recipients, map[string]interface{}{
		"Name":        fileName,
		"SharedItems": sharedItems,
		"Domain":      "https://cxp.proluxe.com/files/",
		"FolderPath":  folderName,
	})
}

func (f *File) SendFile(from, address, path, url string) error {
	fileName := filepath.Base(path)

	knock := services.Knock{
		WorkFlowId: "send-file",
		Email:      from,
	}

	knock.Identify()
	recipients := []string{address}

	return knock.Trigger(recipients, map[string]interface{}{
		"Name": fileName,
		"To":   address,
		"URL":  url,
	})
}

func (f *File) MarkPublic(bucketName, objectName string) (string, error) {
	// Get a reference to the object
	obj := f.Client.Bucket(bucketName).Object(objectName)

	// Set the ACL (Access Control List) to make the object publicly readable
	if err := obj.ACL().Set(f.Context, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set object as public: %v", err)
	}

	// Construct and return the public URL for the object
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	return publicURL, nil

}

func (f *File) MarkPrivate(bucketName, objectName string) error {
	// Get a reference to the object
	obj := f.Client.Bucket(bucketName).Object(objectName)

	// Set the ACL (Access Control List) to make the object private
	if err := obj.ACL().Delete(f.Context, storage.AllUsers); err != nil {
		return fmt.Errorf("failed to set object as private: %v", err)
	}

	return nil
}
