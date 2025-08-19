package model

import (
	"fmt"
	"log"

	"github.com/scottraio/simpleforce"
)

type User struct {
	Id                         string `json:"Id"`
	Name                       string `json:"Name"`
	Email                      string `json:"Email"`
	Phone                      string `json:"Phone"`
	Teams                      string `json:"Teams"`
	IsAdmin                    bool   `json:"IsAdmin"`
	IssueNotifications         bool   `json:"IssueNotifications"`
	NewLeadNotification        bool   `json:"NewLeadNotification"`
	NewOpportunityNotification bool   `json:"NewOpportunityNotification"`
	Notes                      string `json:"Notes"`
}

func FetchAppUsers(client *simpleforce.Client, whereCondition string) []User {
	q := fmt.Sprintf(`
		SELECT Id, Name, Email__c, Phone__c, Last_Login__c, Teams__c, Is_Admin__c, Issue_Notifications__c, New_Lead_Notification__c, New_Opportunity_Notification__c
		FROM App_User__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var users []User
	for _, record := range result.Records {
		u := User{
			Id:                         GetStringField("Id", record),
			Name:                       GetStringField("Name", record),
			Email:                      GetStringField("Email__c", record),
			Phone:                      GetStringField("Phone__c", record),
			Teams:                      GetStringField("Teams__c", record),
			IssueNotifications:         getBoolField("Issue_Notifications__c", record),
			NewLeadNotification:        getBoolField("New_Lead_Notification__c", record),
			NewOpportunityNotification: getBoolField("New_Opportunity_Notification__c", record),
			IsAdmin:                    getBoolField("Is_Admin__c", record),
		}

		users = append(users, u)
	}

	return users
}

type MentionableUser struct {
	Id      string `json:"id"`
	Display string `json:"display"`
}

func FetchUsers(client *simpleforce.Client, whereCondition string) []User {
	q := fmt.Sprintf(`
		SELECT Id, Name, Email__c, Phone__c, Last_Login__c, Teams__c, Is_Admin__c, Issue_Notifications__c, New_Lead_Notification__c, Notes__c, New_Opportunity_Notification__c
		FROM App_User__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var users []User
	for _, record := range result.Records {
		u := User{
			Id:                         GetStringField("Id", record),
			Name:                       GetStringField("Name", record),
			Email:                      GetStringField("Email__c", record),
			Phone:                      GetStringField("Phone__c", record),
			Teams:                      GetStringField("Teams__c", record),
			IssueNotifications:         getBoolField("Issue_Notifications__c", record),
			NewLeadNotification:        getBoolField("New_Lead_Notification__c", record),
			NewOpportunityNotification: getBoolField("New_Opportunity_Notification__c", record),
			IsAdmin:                    getBoolField("Is_Admin__c", record),
			Notes:                      GetStringField("Notes__c", record),
		}

		users = append(users, u)
	}

	return users
}

func (u *User) SaveNotes(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("App_User__c").
		Set("Id", u.Id).
		Set("Notes__c", u.Notes).
		Update()
}

func (u *User) Update(client *simpleforce.Client) *simpleforce.SObject {
	return client.SObject("App_User__c").
		Set("Id", u.Id).
		Set("Phone__c", u.Phone).
		Set("Issue_Notifications__c", u.IssueNotifications).
		Set("New_Lead_Notification__c", u.NewLeadNotification).
		Set("New_Opportunities_Notification__c", u.NewOpportunityNotification).
		Update()
}

func MentionableUsers(client *simpleforce.Client) []MentionableUser {
	users := FetchUsers(client, "Id != null")

	var mentionableUsers []MentionableUser
	for _, user := range users {
		mentionableUsers = append(mentionableUsers, MentionableUser{
			Id:      user.Email,
			Display: user.Name,
		})
	}
	return mentionableUsers
}
