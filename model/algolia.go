package model

import (
	"log"

	"github.com/scottraio/simpleforce"
)

type AlgoliaIndex struct {
	ObjectId string `json:"objectID"`
	Name     string `json:"Name"`
}

func FetchAlgoliaProducts(client *simpleforce.Client) []AlgoliaIndex {
	// Initial query
	query := `
		SELECT Id, Name
		FROM rstk__soprod__c
		WHERE rstk__soprod_activeind__c = true
		ORDER BY CreatedDate DESC
	`

	// Fetch records from Salesforce
	records := FetchRecords(client, query)

	// Return the records in AlgoliaIndex format
	return records
}

func FetchAlgoliaCustomers(client *simpleforce.Client) []AlgoliaIndex {
	// Initial query
	query := `
		SELECT Id, Name
		FROM rstk__socust__c
		ORDER BY CreatedDate DESC
	`
	// Fetch records from Salesforce
	records := FetchRecords(client, query)

	// Return the records in AlgoliaIndex format
	return records
}

func FetchAlgoliaContacts(client *simpleforce.Client) []AlgoliaIndex {
	// Initial query
	query := `
		SELECT Id, Name
		FROM Contact
		ORDER BY CreatedDate DESC
	`

	// Fetch records from Salesforce
	records := FetchRecords(client, query)

	// Return the records in AlgoliaIndex format
	return records
}

func FetchRecords(client *simpleforce.Client, query string) []AlgoliaIndex {
	var allRecords []simpleforce.SObject

	result, err := client.Query(query)
	if err != nil {
		log.Fatalf("Error executing initial query: %v", err)
	}
	allRecords = append(allRecords, result.Records...)

	// Keep fetching more records if available
	for !result.Done && result.NextRecordsURL != "" {
		result, err = client.QueryMore(result.NextRecordsURL)
		if err != nil {
			log.Fatalf("Error during QueryMore: %v", err)
		}
		allRecords = append(allRecords, result.Records...)
	}

	// Convert to AlgoliaIndex format
	return setAlgoliaIndexFromSObjects(allRecords)
}

func setAlgoliaIndexFromSObjects(records []simpleforce.SObject) []AlgoliaIndex {
	var products []AlgoliaIndex
	for _, r := range records {
		products = append(products, AlgoliaIndex{
			ObjectId: r["Id"].(string),
			Name:     r["Name"].(string),
		})
	}
	return products
}
