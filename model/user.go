package model

import (
	"fmt"
	"log"

	"github.com/simpleforce/simpleforce"
)

type User struct {
	Id                         string `json:"Id"`
	Name                       string `json:"Name"`
	Email                      string `json:"Email"`
	Phone                      string `json:"Phone"`
	IssueNotifications         bool   `json:"IssueNotifications"`
	NewLeadNotification        bool   `json:"NewLeadNotification"`
	NewOpportunityNotification bool   `json:"NewOpportunityNotification"`
}

func FetchUsers(client *simpleforce.Client, whereCondition string) []User {
	q := fmt.Sprintf(`
		SELECT Id, Name, rstk__syusr_empl_email__c, rstk__syusr_phone__c, Issue_Notifications__c, New_Lead_Notification__c, New_Opportunities_Notification__c
		FROM rstk__syusr__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var users []User
	for _, record := range result.Records {
		u := User{
			Id:                         getStringField("Id", record),
			Name:                       getStringField("Name", record),
			Email:                      getStringField("rstk__syusr_empl_email__c", record),
			Phone:                      getStringField("rstk__syusr_phone__c", record),
			IssueNotifications:         getBoolField("Issue_Notifications__c", record),
			NewLeadNotification:        getBoolField("New_Lead_Notification__c", record),
			NewOpportunityNotification: getBoolField("New_Opportunities_Notification__c", record),
		}

		users = append(users, u)
	}

	return users
}

func (u *User) Update(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("rstk__syusr__c").
		Set("Id", u.Id).
		Set("rstk__syusr_phone__c", u.Phone).
		Set("Issue_Notifications__c", u.IssueNotifications).
		Set("New_Lead_Notification__c", u.NewLeadNotification).
		Set("New_Opportunities_Notification__c", u.NewOpportunityNotification).
		Update()
}
