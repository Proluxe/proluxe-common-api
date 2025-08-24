package model

import (
	"log"

	"github.com/scottraio/simpleforce"
)

type AlgoliaIndex struct {
	ObjectId string `json:"objectID"`
	Name     string `json:"Name"`
}

/*
FetchAlgoliaProducts fetches all active products from Salesforce and returns them in AlgoliaIndex format.
Uses "Id" as the objectID and "Name" as the name.
*/
func FetchAlgoliaProducts(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Name
		FROM rstk__soprod__c
		WHERE rstk__soprod_activeind__c = true
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Id", "Name")
	return records
}

/*
FetchAlgoliaCustomers fetches all active customers from Salesforce and returns them in AlgoliaIndex format.
Uses "Id" as the objectID and "Name" as the name.
*/
func FetchAlgoliaCustomers(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Name
		FROM rstk__socust__c
		WHERE Archived__c = false AND rstk__socust_inactive__c = false
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Id", "Name")
	return records
}

/*
FetchAlgoliaContacts fetches all contacts from Salesforce and returns them in AlgoliaIndex format.
Uses "Id" as the objectID and "Display_Name__c" as the name.
*/
func FetchAlgoliaContacts(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Display_Name__c
		FROM Contact
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Id", "Display_Name__c")
	return records
}

/*
FetchAlgoliaParts fetches all parts from Salesforce and returns them in AlgoliaIndex format.
Uses "Id" as the objectID and "Name" as the name.
*/
func FetchAlgoliaParts(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Name
		FROM rstk__icitem__c
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Id", "Name")
	return records
}

/*
FetchAlgoliaVendors fetches all vendors from Salesforce (rstk__povend__c) and returns them in AlgoliaIndex format.
Uses "Id" as the objectID and "Name" as the name.
*/
func FetchAlgoliaVendors(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Name
		FROM rstk__povend__c
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Id", "Name")
	return records
}

/*
FetchAlgoliaCatalog fetches all catalog models from Salesforce and returns them in AlgoliaIndex format.
Uses "Model_Number__c" as the objectID and "Name" as the name.
*/
func FetchAlgoliaCatalog(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Model_Number__c, Name
		FROM Model__c
		ORDER BY CreatedDate DESC
	`
	records := FetchRecords(client, query, "Model_Number__c", "Name")
	return records
}

/*
FetchRecords fetches records from Salesforce using the provided SOQL query,
and maps them to AlgoliaIndex using the specified idField and nameField.
If not provided, defaults to "Id" and "Name".
*/
func FetchRecords(client *simpleforce.Client, query string, fields ...string) []AlgoliaIndex {
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

	// Determine the id and name fields to use
	idField := "Id"
	nameField := "Name"
	if len(fields) == 1 {
		idField = fields[0]
	} else if len(fields) >= 2 {
		idField = fields[0]
		nameField = fields[1]
	}

	// Convert to AlgoliaIndex format
	return setAlgoliaIndexFromSObjectsWithFields(allRecords, idField, nameField)
}

/*
setAlgoliaIndexFromSObjectsWithFields maps Salesforce SObjects to AlgoliaIndex using the specified id and name fields.
*/
func setAlgoliaIndexFromSObjectsWithFields(records []simpleforce.SObject, idField, nameField string) []AlgoliaIndex {
	var products []AlgoliaIndex
	for _, r := range records {
		idVal, idOk := r[idField].(string)
		nameVal, nameOk := r[nameField].(string)
		if !idOk || !nameOk {
			continue // skip if missing or wrong type
		}
		products = append(products, AlgoliaIndex{
			ObjectId: idVal,
			Name:     nameVal,
		})
	}
	return products
}

// TODO: Unit tests are needed for FetchRecords, setAlgoliaIndexFromSObjectsWithFields, and all FetchAlgolia* functions.
// Tests should cover:
// - Mapping of idField and nameField
// - Handling of missing/invalid fields
// - SOQL query error handling
// - End-to-end mapping from Salesforce SObject to AlgoliaIndex
