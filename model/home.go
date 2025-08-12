package model

import (
	"fmt"
	"log"

	"github.com/scottraio/simpleforce"
)

type HomeRecord struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	CreatedDate string `json:"CreatedDate"`

	// Leads
	Address string `json:"Address,omitempty"`
	Company string `json:"Company,omitempty"`
	Email   string `json:"Email,omitempty"`
	Phone   string `json:"Phone,omitempty"`
	Product string `json:"Product,omitempty"`
	Status  string `json:"Status,omitempty"`
	Summary string `json:"Summary,omitempty"`

	// Opportunities
	Amount      string `json:"Amount,omitempty"`
	StageName   string `json:"StageName,omitempty"`
	Description string `json:"Description,omitempty"`

	// Tasks
	AssignedTo string `json:"AssignedTo,omitempty"`
	Project    string `json:"Project,omitempty"`
	ProjectId  string `json:"ProjectId,omitempty"`
	ParentCard string `json:"ParentCard,omitempty"`
	Tags       string `json:"Tags,omitempty"`
	DueDate    string `json:"DueDate,omitempty"`
}

func FetchHomeLeads(client *simpleforce.Client) []HomeRecord {
	// Initial query
	query := `
		SELECT Id, Name, Address, Company, CreatedDate, Email, Phone, Product__c, Status, Summary__c
		FROM Lead
		WHERE Status != 'Closed' 
		ORDER BY LastModifiedDate DESC 
		LIMIT 10
	`

	// Fetch records from Salesforce
	records := FetchHomeRecords(client, query)

	var allRecords []HomeRecord
	for _, record := range records {
		// Add additional fields specific to leads
		lead := HomeRecord{
			Id:          GetStringField("Id", record),
			Name:        GetStringField("Name", record),
			CreatedDate: GetStringField("CreatedDate", record),
			Address:     GetStringField("Address", record),
			Company:     GetStringField("Company", record),
			Email:       GetStringField("Email", record),
			Phone:       GetStringField("Phone", record),
			Product:     GetStringField("Product__c", record),
			Status:      GetStringField("Status", record),
			Summary:     GetStringField("Summary__c", record),
		}

		allRecords = append(allRecords, lead)
	}

	// Return the records in AlgoliaIndex format
	return allRecords
}

func FetchHomeOpportunities(client *simpleforce.Client) []HomeRecord {
	// Initial query
	query := `
		SELECT Id, Name, Amount, StageName, CreatedDate, Description
		FROM Opportunity
		WHERE StageName != 'Closed Won' AND StageName != 'Closed Lost' 
		ORDER BY LastModifiedDate DESC 
		LIMIT 10
	`
	// Fetch records from Salesforce
	records := FetchHomeRecords(client, query)

	var allRecords []HomeRecord
	for _, record := range records {

		// Populate the HomeRecord struct
		opportunity := HomeRecord{
			Id:          GetStringField("Id", record),
			Name:        GetStringField("Name", record),
			CreatedDate: GetStringField("CreatedDate", record),
			Amount:      GetStringField("Amount", record),
			StageName:   GetStringField("StageName", record),
			Description: GetStringField("Description", record),
		}

		allRecords = append(allRecords, opportunity)
	}

	// Return the records in AlgoliaIndex format
	return allRecords
}

func FetchHomeTasks(client *simpleforce.Client, email string) []HomeRecord {
	// Initial query
	query := fmt.Sprintf(`
		SELECT Id, Name, Project__c, Project__r.Name, Description__c, Status__c, Due_Date__c, Assigned_To__c, Assigned_To__r.Name, Assigned_To__r.Email__c, Parent_Card__c, Tags__c
		FROM Project_Card__c
		WHERE Assigned_To__r.Email__c = '%s'
		ORDER BY Name
	`, email)

	// Fetch records from Salesforce
	records := FetchHomeRecords(client, query)

	var allRecords []HomeRecord
	for _, record := range records {
		// Populate the HomeRecord struct
		task := HomeRecord{
			Id:          GetStringField("Id", record),
			Name:        GetStringField("Name", record),
			Project:     GetStringField("Project__r.Name", record),
			ProjectId:   GetStringField("Project__c", record),
			Description: GetStringField("Description__c", record),
			Status:      GetStringField("Status__c", record),
			DueDate:     GetStringField("Due_Date__c", record),
			AssignedTo:  GetStringField("Assigned_To__r.Name", record),
			ParentCard:  GetStringField("Parent_Card__c", record),
			Tags:        GetStringField("Tags__c", record),
		}

		allRecords = append(allRecords, task)
	}

	// Return the records in AlgoliaIndex format
	return allRecords
}

func FetchHomeRecords(client *simpleforce.Client, query string) []simpleforce.SObject {

	result, err := client.Query(query)
	if err != nil {
		log.Fatalf("Error executing initial query: %v", err)
	}

	return result.Records
}
