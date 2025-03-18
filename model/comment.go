package model

import (
	"fmt"
	"log"

	services "github.com/Proluxe/proluxe-common-api/services"
	"github.com/simpleforce/simpleforce"
)

type Comment struct {
	Id         string `json:"Id"`
	CreatedBy  string `json:"CreatedBy"`
	Message    string `json:"Message"`
	RecordID   string `json:"RecordID"`
	RecordType string `json:"RecordType"` // This needs to be the SF object name
	Subject    string `json:"Subject"`
	Avatar     string `json:"Avatar"`
}

// Fetch cases from Salesforce
func FetchComments(client *simpleforce.Client, whereCondition string) []Comment {
	// Construct the query to fetch cases
	q := fmt.Sprintf(`
		SELECT Id, Created_By__c, Message__c, Name, Record_ID__c, Avatar__c, Record_Type__c
		FROM Comment__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error fetching comments: ", err)
	}

	var comments []Comment

	for _, record := range result.Records {
		// Populate the Case struct
		c := Comment{
			Id:         getStringField("Id", record),
			CreatedBy:  getStringField("Created_By__c", record),
			Message:    getStringField("Message__c", record),
			Subject:    getStringField("Name", record),
			RecordID:   getStringField("Record_ID__c", record),
			RecordType: getStringField("Record_Type__c", record),
			Avatar:     getStringField("Avatar__c", record),
		}

		comments = append(comments, c)
	}

	return comments
}

func (c *Comment) Create(client *simpleforce.Client) {
	client.SObject("Comment__c").
		Set("Created_By__c", c.CreatedBy).
		Set("Message__c", c.Message).
		Set("Name", truncateString(c.Message, 50)).
		Set("Record_ID__c", c.RecordID).
		Set("Record_Type__c", c.RecordType).
		Set("Avatar__c", c.Avatar).
		Create()
}

func truncateString(str string, num int) string {
	if len(str) > num {
		return str[:num]
	}
	return str
}

func (c *Comment) SendNotificationEmail(client *simpleforce.Client) {
	name := c.GetRecordName(client)

	knock := services.Knock{
		WorkFlowId: "new-comment",
		Email:      c.CreatedBy,
	}

	fmt.Printf("Identifying user %s\n", c.CreatedBy)

	knock.Identify()

	var whereCondition string
	switch c.RecordType {
	case "Opportunity":
		whereCondition = "New_Opportunity_Notification__c = true"
	case "Lead":
		whereCondition = "New_Lead_Notification__c = true"
	default:
		whereCondition = "Id != null"
	}

	users := FetchUsers(client, whereCondition)

	recipients := make([]string, len(users))
	for i, user := range users {
		recipients[i] = user.Email
	}

	err := knock.Trigger(recipients, map[string]any{
		"Name":    name,
		"Message": c.Message,
		"From":    c.CreatedBy,
		"Object":  c.RecordType,
		"Url":     c.LinkToRecord(client),
		"Domain":  "https://crm.proluxe.com",
	})

	if err != nil {
		log.Println("Error sending email: ", err)
	}
}

func (c *Comment) Delete(client *simpleforce.Client) error {
	err := client.SObject("Comment__c").Set("Id", c.Id).Delete()
	return err
}

func (c *Comment) GetRecordName(client *simpleforce.Client) string {
	q := fmt.Sprintf(`
		SELECT Name
		FROM %s
		WHERE Id = '%s'
	`, c.RecordType, c.RecordID)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error fetching record name: ", err)
	}
	return getStringField("Name", result.Records[0])
}

func (c *Comment) LinkToRecord(client *simpleforce.Client) string {
	var link string

	switch c.RecordType {
	case "Opportunity":
		link = fmt.Sprintf("/opportunities/%s", c.RecordID)
	default:
		link = fmt.Sprintf("/record/%s", c.RecordID)
	}

	return link
}
