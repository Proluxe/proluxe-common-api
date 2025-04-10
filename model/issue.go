package model

import (
	"fmt"
	"log"

	services "github.com/Proluxe/proluxe-common-api/services"
	"github.com/scottraio/simpleforce"
)

type Issue struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Closed      bool   `json:"Closed"`

	Files    []SharedFile `json:"Files"`
	Comments []Comment    `json:"Comments"`
}

func FetchIssues(client *simpleforce.Client, whereCondition string) []Issue {
	q := fmt.Sprintf(`
		SELECT Id, Name, Description__c, Closed__c
		FROM Issue__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error fetching issues: ", err)
	}

	var issues []Issue

	for _, record := range result.Records {
		i := Issue{
			Id:          getStringField("Id", record),
			Name:        getStringField("Name", record),
			Description: getStringField("Description__c", record),
			Closed:      getBoolField("Closed__c", record),
		}

		issues = append(issues, i)
	}

	return issues
}

func (i *Issue) Create(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("Issue__c").
		Set("Name", i.Name).
		Set("Description__c", i.Description).
		Set("Closed__c", false).
		Create()
}

func (i *Issue) Close(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("Issue__c").
		Set("Id", i.Id).
		Set("Closed__c", true).
		Update()
}

func (i *Issue) Delete(client *simpleforce.Client) error {
	return client.SObject("Issue__c").
		Set("Id", i.Id).
		Delete()
}

func (i *Issue) Update(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("Issue__c").
		Set("Id", i.Id).
		Set("Name", i.Name).
		Set("Description__c", i.Description).
		Update()
}

func (i *Issue) SendClosedIssueNotification(client *simpleforce.Client) {
	knock := services.Knock{
		WorkFlowId: "closed-issue",
	}

	users := FetchUsers(client, "Issue_Notifications__c = true")

	recipients := make([]string, len(users))
	for i, user := range users {
		recipients[i] = user.Email
	}

	var shortDescription string
	if len(i.Description) > 144 {
		shortDescription = i.Description[:141] + "..."
	} else {
		shortDescription = i.Description
	}

	url := fmt.Sprintf("https://crm.proluxe.com/issues/%s", i.Id)

	knock.Trigger(recipients, map[string]interface{}{
		"Name":             i.Name,
		"Description":      i.Description,
		"ShortDescription": shortDescription,
		"Url":              url,
	})
}

func (i *Issue) SendNewIssueNotification(client *simpleforce.Client) {
	knock := services.Knock{
		WorkFlowId: "new-issue",
	}

	users := FetchUsers(client, "Issue_Notifications__c = true")

	recipients := make([]string, len(users))
	for i, user := range users {
		recipients[i] = user.Email
	}

	var shortDescription string
	if len(i.Description) > 144 {
		shortDescription = i.Description[:141] + "..."
	} else {
		shortDescription = i.Description
	}

	url := fmt.Sprintf("https://crm.proluxe.com/issues/%s", i.Id)

	knock.Trigger(recipients, map[string]interface{}{
		"Name":             i.Name,
		"Description":      i.Description,
		"ShortDescription": shortDescription,
		"Url":              url,
	})
}

func (i *Issue) AttachRelatedObjects(client *simpleforce.Client) {
	i.Files = FetchFiles(client, "Object_Id__c = '"+i.Id+"' AND Object__c = 'Issue__c'")
}
