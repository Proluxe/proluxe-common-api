package model

import (
	"fmt"
	"log"

	"github.com/scottraio/simpleforce"
)

// PinnedLink represents a saved link for a user
type PinnedLink struct {
	Id    string `json:"Id"`
	App   string `json:"App"`
	Email string `json:"Email"`
	Name  string `json:"Name"`
	Path  string `json:"Path"`
}

// CreatePinnedLink saves a pinned link in Salesforce
func CreatePinnedLink(client *simpleforce.Client, link PinnedLink) (PinnedLink, error) {
	// 🛑 Debugging: Log incoming request
	log.Println("📌 Creating Pinned Link in Salesforce:", link)

	newLink := client.SObject("Pinned_Link__c").
		Set("Email__c", link.Email).
		Set("Name", link.Name).
		Set("App__c", link.App).
		Set("Path__c", link.Path)

	// 🔍 Execute Create() and capture the returned object
	createdLink := newLink.Create()

	// 🚨 Handle failure case properly
	if createdLink == nil || createdLink.ID() == "" {
		log.Println("❌ Failed to create pinned link: No ID returned from Salesforce")
		return link, fmt.Errorf("failed to create pinned link: No ID returned from Salesforce")
	}

	log.Printf("✅ Pinned Link Created in Salesforce with ID: %s\n", createdLink.ID())

	link.Id = createdLink.ID()

	return link, nil
}

// FetchPinnedLinks retrieves all pinned links for a given user from Salesforce
func FetchPinnedLinks(client *simpleforce.Client, email string, app ...string) ([]PinnedLink, error) {
	query := fmt.Sprintf(`
		SELECT Id, Name, Path__c, Email__c, App__c
		FROM Pinned_Link__c 
		WHERE Email__c = '%s'`, email)

	// Add app filter if provided
	if len(app) > 0 && app[0] != "" {
		query += fmt.Sprintf(" AND App__c = '%s'", app[0])
	}

	log.Println("🔍 Fetching pinned links with query:", query)

	result, err := client.Query(query)
	if err != nil {
		log.Println("❌ Salesforce Query Error:", err)
		return nil, err
	}

	var links []PinnedLink
	for _, record := range result.Records {
		link := PinnedLink{
			Id:    GetStringField("Id", record),
			Name:  GetStringField("Name", record),
			Path:  GetStringField("Path__c", record),
			Email: GetStringField("Email__c", record),
			App:   GetStringField("App__c", record),
		}
		links = append(links, link)
	}

	log.Println("✅ Retrieved pinned links:", links)
	return links, nil
}

func DeletePinnedLink(client *simpleforce.Client, id string) error {
	log.Println("🚫 Deleting Pinned Link with ID:", id)

	err := client.SObject("Pinned_Link__c").Set("Id", id).Delete()
	if err != nil {
		log.Println("❌ Failed to delete pinned link:", err)
		return err
	}

	log.Println("✅ Pinned Link Deleted Successfully")
	return nil
}
