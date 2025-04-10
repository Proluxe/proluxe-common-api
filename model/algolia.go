package model

import (
	"log"

	"github.com/simpleforce/simpleforce"
)

type AlgoliaIndex struct {
	ObjectId string `json:"objectID"`
	Name     string `json:"Name"`
}

func FetchAlgoliaProducts(client *simpleforce.Client) []AlgoliaIndex {
	query := `
		SELECT Id, Name
		FROM rstk__soprod__c
		WHERE rstk__soprod_activeind__c = true
		ORDER BY CreatedDate DESC
	`

	// Execute the initial query
	result, err := client.Query(query)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	// Set the Algolia index
	products := setAlgoliaIndex(result)

	return products
}

func setAlgoliaIndex(result *simpleforce.QueryResult) []AlgoliaIndex {
	var index []AlgoliaIndex

	for _, record := range result.Records {
		index = append(index, AlgoliaIndex{
			ObjectId: record["Id"].(string),
			Name:     record["Name"].(string),
		})
	}

	return index
}
