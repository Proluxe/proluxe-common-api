package model

import (
	"fmt"
	"log"
	"strings"

	"github.com/scottraio/simpleforce"
)

type CommentMention struct {
	Id         string `json:"Id"`
	CommentId  string `json:"CommentId"`
	Email      string `json:"Email"`
	RecordID   string `json:"RecordID"`
	RecordType string `json:"RecordType"`
}

// Fetch cases from Salesforce
func FetchMentions(client *simpleforce.Client, whereCondition string) []CommentMention {
	// Construct the query to fetch cases
	q := fmt.Sprintf(`
		SELECT Id, Email__c, Record_ID__c, Record_Type__c
		FROM Comment_Mentioned_User__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error fetching mentions: ", err)
	}

	var mentions []CommentMention

	for _, record := range result.Records {
		// Populate the Case struct
		c := CommentMention{
			Id:         getStringField("Id", record),
			RecordID:   getStringField("Record_ID__c", record),
			RecordType: getStringField("Record_Type__c", record),
			Email:      getStringField("Email__c", record),
			CommentId:  getStringField("Comment__c", record),
		}

		mentions = append(mentions, c)
	}

	return mentions
}

// AddMentions adds multiple mentions, ensuring they don't already exist based on RecordType and RecordID.
func AddMentions(client *simpleforce.Client, mentions []CommentMention) error {
	if len(mentions) == 0 {
		return nil
	}

	// Constructing a WHERE clause for batch checking
	var conditions []string
	for _, mention := range mentions {
		conditions = append(conditions, fmt.Sprintf("(Record_ID__c = '%s' AND Record_Type__c = '%s' AND Email__c = '%s')",
			mention.RecordID, mention.RecordType, mention.Email))
	}

	// Query to check existing mentions in batch
	query := fmt.Sprintf(`
		SELECT Record_ID__c, Record_Type__c, Email__c 
		FROM Comment_Mentioned_User__c
		WHERE %s
	`, strings.Join(conditions, " OR "))

	existingRecords, err := client.Query(query)
	if err != nil {
		return fmt.Errorf("error checking existing mentions: %v", err)
	}

	// Create a map of existing mentions for quick lookup
	existingSet := make(map[string]bool)
	for _, record := range existingRecords.Records {
		key := fmt.Sprintf("%s|%s|%s",
			getStringField("Record_ID__c", record),
			getStringField("Record_Type__c", record),
			getStringField("Email__c", record),
		)
		existingSet[key] = true
	}

	// Insert only the mentions that don't already exist
	for _, mention := range mentions {
		key := fmt.Sprintf("%s|%s|%s", mention.RecordID, mention.RecordType, mention.Email)
		if _, exists := existingSet[key]; exists {
			log.Printf("Skipping existing mention: RecordID %s, RecordType %s, Email %s", mention.RecordID, mention.RecordType, mention.Email)
			continue
		}

		// Create a new mention
		client.SObject("Comment_Mentioned_User__c").
			Set("Record_ID__c", mention.RecordID).
			Set("Record_Type__c", mention.RecordType).
			Set("Email__c", mention.Email).
			Set("Comment__c", mention.CommentId).
			Create()

		log.Printf("Mention added successfully: RecordID %s, RecordType %s, Email %s Comment Id %s", mention.RecordID, mention.RecordType, mention.Email, mention.CommentId)
	}

	return nil
}

func (c *CommentMention) Delete(client *simpleforce.Client) error {
	err := client.SObject("Comment_Mentioned_User__c").Set("Id", c.Id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete comment mention: %v", err)
	}
	return nil
}
