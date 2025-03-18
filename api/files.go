package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"

	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"

	model "github.com/Proluxe/proluxe-common-api/model"
)

type BucketPayload struct {
	Contents    []map[string]any   `json:"Contents"`
	SharedFiles []model.SharedFile `json:"SharedFiles"`
}

// GET_BUCKET_CONTENTS lists the contents of the specified path in the 'common_production' bucket.
func GET_BUCKET_CONTENTS(c *gin.Context, App *util.App) {
	bucketName := "common_production"
	path := c.Query("path") // Path within the bucket

	if path == "root/" {
		path = "/"
	}

	file := model.New(c, App)

	contents, err := file.ListBucketContents(bucketName, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bucket contents"})
		return
	}
	bucket := &BucketPayload{
		Contents:    contents,
		SharedFiles: file.FetchFilesForPath(path),
	}

	c.JSON(http.StatusOK, bucket)
}

func POST_SEND_FILES(c *gin.Context, App *util.App) {
	var payload struct {
		Email string `json:"Email"`
		Path  string `json:"Path"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	file := model.New(c, App)
	url, err := file.MarkPublic("common_production", payload.Path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark file as public"})
		return
	}

	if url != "" {
		currentUserEmail, _, _ := GetCurrentUser(c)

		err = file.SendFile(currentUserEmail, payload.Email, payload.Path, url)
		if err != nil {
			fmt.Printf("Failed to send email: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func SERVE_FILE_FROM_BUCKET(c *gin.Context) {
	bucketName := "common_production"
	objectName := c.Param("path")

	// Check and clean up the object path
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	// Create storage client
	client, err := storage.NewClient(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage client: " + err.Error()})
	}
	defer client.Close()

	// Open a reader to the object
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or unable to read: " + err.Error()})
	}
	defer rc.Close()

	// Use http.ServeContent for better performance and file handling
	contentType := rc.ContentType()
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(objectName)))

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, rc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content: " + err.Error()})
	}
	reader := bytes.NewReader(buf.Bytes())

	c.DataFromReader(http.StatusOK, rc.Attrs.Size, rc.ContentType(), reader, nil)
}

func UPLOAD_FILE_TO_BUCKET(c *gin.Context, App *util.App) {
	bucketName := "common_production"
	path := c.Query("path") // Path within the bucket

	file := model.New(c, App)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
	}

	// Open the uploaded file
	rawFile, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
	}
	defer rawFile.Close()

	// Upload the file to the bucket
	objectName := path + fileHeader.Filename
	wc := file.Client.Bucket(bucketName).Object(objectName).NewWriter(context.Background())

	if _, err := io.Copy(wc, rawFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
	}

	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize file upload"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded to bucket %s at %s", fileHeader.Filename, bucketName, objectName)})
}

// CREATE_FOLDER_IN_BUCKET creates a folder in the 'common_production' bucket.
func CREATE_FOLDER_IN_BUCKET(c *gin.Context, App *util.App) {
	bucketName := "common_production"

	file := model.New(c, App)

	var json struct {
		Path string `json:"Path"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}

	folderName := json.Path // Path within the bucket

	err := file.CreateFolder(bucketName, folderName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Folder %s created in bucket %s", folderName, bucketName)})
}

// DELETE_FILE_FROM_BUCKET deletes a file from the specified path in the 'common_production' bucket.
func DELETE_FILE_FROM_BUCKET(c *gin.Context, App *util.App) {
	bucketName := "common_production"
	objectName := c.Query("path") // Path within the bucket

	file := model.New(c, App)

	err := file.DeleteFile(bucketName, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s deleted from bucket %s", objectName, bucketName)})
}

func POST_SHARE_FILES(c *gin.Context, App *util.App) {
	var payload struct {
		SharedItems []model.SharedItem `json:"SharedItems"`
		Path        string             `json:"Path"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}

	file := model.New(c, App)

	for _, item := range payload.SharedItems {
		externalID := fmt.Sprintf("%s-%s-%s", item.Object, item.ObjectId, payload.Path)
		encodedExternalID := base64.StdEncoding.EncodeToString([]byte(externalID))

		file.SF.SObject("common_File__c").
			Set("ExternalIDField", "ExternalId__c").
			Set("ExternalId__c", encodedExternalID).
			Set("Object__c", item.Object).
			Set("Object_Id__c", item.ObjectId).
			Set("Object_Name__c", item.ObjectName).
			Set("Path__c", payload.Path).
			Upsert()
	}

	currentUserEmail, _, _ := GetCurrentUser(c)
	err := file.SendFileShareConfirmationEmail(currentUserEmail, payload.Path, payload.SharedItems)

	if err != nil {
		fmt.Printf("Failed to attach file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files shared successfully"})
}

func POST_MAKE_PUBLIC(c *gin.Context, App *util.App) {
	path := c.Query("path") // Path within the bucket

	file := model.New(c, App)
	url, err := file.MarkPublic("common_production", path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark file as public"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func POST_MAKE_PRIVATE(c *gin.Context, App *util.App) {
	path := c.Query("path") // Path within the bucket

	file := model.New(c, App)
	err := file.MarkPrivate("common_production", path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark file as private"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File marked as private"})
}
