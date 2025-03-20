package model

import (
	"fmt"
	"log"

	services "github.com/Proluxe/proluxe-common-api/services"
	"github.com/simpleforce/simpleforce"
)

type Comment struct {
	Id             string           `json:"Id"`
	CreatedBy      string           `json:"CreatedBy"`
	CreatedByName  string           `json:"CreatedByName"`
	Message        string           `json:"Message"`
	RecordID       string           `json:"RecordID"`
	RecordType     string           `json:"RecordType"` // This needs to be the SF object name
	RecordName     string           `json:"RecordName"`
	Subject        string           `json:"Subject"`
	Avatar         string           `json:"Avatar"`
	MentionedUsers []CommentMention `json:"MentionedUsers"`
}

// Fetch cases from Salesforce
func FetchComments(client *simpleforce.Client, whereCondition string) []Comment {
	// Construct the query to fetch cases
	q := fmt.Sprintf(`
		SELECT Id, Created_By__c, Message__c, Name, Record_ID__c, Avatar__c, Record_Type__c, Record_Name__c, Created_By_Name__c
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
			Id:            getStringField("Id", record),
			CreatedBy:     getStringField("Created_By__c", record),
			Message:       getStringField("Message__c", record),
			Subject:       getStringField("Name", record),
			RecordID:      getStringField("Record_ID__c", record),
			RecordType:    getStringField("Record_Type__c", record),
			RecordName:    getStringField("Record_Name__c", record),
			Avatar:        getStringField("Avatar__c", record),
			CreatedByName: getStringField("Created_By_Name__c", record),
		}

		comments = append(comments, c)
	}

	return comments
}

func (c *Comment) Create(client *simpleforce.Client) {
	sobj := client.SObject("Comment__c").
		Set("Created_By__c", c.CreatedBy).
		Set("Message__c", c.Message).
		Set("Name", truncateString(c.Message, 50)).
		Set("Record_ID__c", c.RecordID).
		Set("Record_Type__c", c.RecordType).
		Set("Avatar__c", c.Avatar).
		Set("Created_By_Name__c", c.CreatedByName).
		Create().Get()

	// Add mentions
	mentions := []CommentMention{}
	c.MentionedUsers = append([]CommentMention{{
		Email:      c.CreatedBy,
		CommentId:  sobj.ID(),
		RecordID:   c.RecordID,
		RecordType: c.RecordType,
	}}, c.MentionedUsers...)
	for _, user := range c.MentionedUsers {
		mentions = append(mentions, CommentMention{
			Email:      user.Email,
			CommentId:  sobj.ID(),
			RecordID:   c.RecordID,
			RecordType: c.RecordType,
		})
	}

	AddMentions(client, mentions)
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

	users := FetchMentions(client, fmt.Sprintf("Record_ID__c = '%s' AND Record_Type__c = '%s'", c.RecordID, c.RecordType))

	recipients := make([]string, len(users))
	for i, user := range users {
		recipients[i] = user.Email
	}

	err := knock.Trigger(recipients, map[string]any{
		"Name":       name,
		"Message":    c.Message,
		"From":       c.CreatedBy,
		"FromName":   c.CreatedByName,
		"Object":     c.normalizeObjectName(),
		"ObjectName": c.RecordName,
		"Url":        c.LinkToRecord(client),
		"Domain":     "https://crm.proluxe.com",
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

func (c *Comment) normalizeObjectName() string {
	objectNames := map[string]string{
		"Lead":            "Lead",
		"Opportunity":     "Opportunity",
		"Case":            "Case",
		"Account":         "Account",
		"Contact":         "Contact",
		"Event__c":        "Event",
		"Issue":           "Issue",
		"rstk__soprod__c": "Product",
	}

	if name, exists := objectNames[c.RecordType]; exists {
		return name
	}
	return "Record"
}

func (c *Comment) LinkToRecord(client *simpleforce.Client) string {
	recordLinks := map[string]string{
		"Lead":            "leads",
		"Opportunity":     "opportunities",
		"Case":            "cases",
		"Account":         "accounts",
		"Contact":         "contacts",
		"Event__c":        "events",
		"Issue__c":        "issues",
		"rstk__soprod__c": "products",
	}

	if path, exists := recordLinks[c.RecordType]; exists {
		return fmt.Sprintf("/%s/%s", path, c.RecordID)
	}
	return fmt.Sprintf("/record/%s", c.RecordID)
}
