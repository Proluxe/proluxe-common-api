package model

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/scottraio/simpleforce"
)

type Event struct {
	Id               string    `json:"Id"`
	Name             string    `json:"Name"`
	AddressLine1     string    `json:"AddressLine1"`
	AddressLine2     string    `json:"AddressLine2"`
	City             string    `json:"City"`
	StateProvince    string    `json:"StateProvince"`
	ZipCode          string    `json:"ZipCode"`
	Description      string    `json:"Description"`
	StartDateTime    time.Time `json:"StartDateTime"`
	EndDateTime      time.Time `json:"EndDateTime"`
	Type             string    `json:"Type"`
	CreatedById      string    `json:"CreatedById"`
	LastModifiedById string    `json:"LastModifiedById"`
	OwnerId          string    `json:"OwnerId"`
	Traveling        bool      `json:"Traveling"`
	Informational    bool      `json:"Informational"`
	TextAddress      string    `json:"TextAddress"`
}

func FetchEvents(client *simpleforce.Client, whereCondition string) []Event {
	// Construct the query to fetch events
	q := fmt.Sprintf(`
		SELECT Id, Name, Address_Line_1__c, Address_Line_2__c, City__c, State_Province__c, Zip_Code__c, Description__c, 
		Start_Date_Time__c, End_Date_Time__c, Type__c, CreatedById, LastModifiedById, OwnerId, Travel__c, Informational__c, TextAddress__c
		FROM Events__c
		WHERE %s
	`, whereCondition)

	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var events []Event

	for _, record := range result.Records {
		// Populate the Event struct
		event := Event{
			Id:               getStringField("Id", record),
			Name:             getStringField("Name", record),
			AddressLine1:     getStringField("Address_Line_1__c", record),
			AddressLine2:     getStringField("Address_Line_2__c", record),
			City:             getStringField("City__c", record),
			StateProvince:    getStringField("State_Province__c", record),
			ZipCode:          getStringField("Zip_Code__c", record),
			Description:      getStringField("Description__c", record),
			StartDateTime:    convertToTime(record["Start_Date_Time__c"].(string)),
			EndDateTime:      convertToTime(record["End_Date_Time__c"].(string)),
			Type:             getStringField("Type__c", record),
			CreatedById:      getStringField("CreatedById", record),
			LastModifiedById: getStringField("LastModifiedById", record),
			OwnerId:          getStringField("OwnerId", record),
			Traveling:        getBoolField("Travel__c", record),
			Informational:    getBoolField("Informational__c", record),
			TextAddress:      getStringField("TextAddress__c", record),
		}

		events = append(events, event)
	}

	return events
}

func CreateEvent(client *simpleforce.Client, newEvent Event) *simpleforce.SObject {
	fmt.Printf("Creating new event: %v\n", newEvent)
	return SetEvent(client, newEvent).
		Create()
}

func UpdateEvent(client *simpleforce.Client, event Event) *simpleforce.SObject {
	fmt.Printf("Updating event: %v\n", event)
	return SetEvent(client, event).
		Set("Id", event.Id).
		Update()
}

func SetEvent(client *simpleforce.Client, event Event) *simpleforce.SObject {
	return client.SObject("Events__c").
		Set("Name", event.Name).
		Set("Address_Line_1__c", event.AddressLine1).
		Set("Address_Line_2__c", event.AddressLine2).
		Set("City__c", event.City).
		Set("State_Province__c", event.StateProvince).
		Set("Zip_Code__c", event.ZipCode).
		Set("Description__c", event.Description).
		Set("Start_Date_Time__c", event.StartDateTime).
		Set("End_Date_Time__c", event.EndDateTime).
		Set("Type__c", event.Type).
		Set("Travel__c", event.Traveling)
}

func SortedEvents(events []Event) []Event {

	// Separate events into future and past
	var futureEvents, pastEvents []Event
	now := time.Now()

	for _, event := range events {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		nowPT := now.In(loc)
		if event.StartDateTime.After(nowPT) {
			futureEvents = append(futureEvents, event)
		} else {
			pastEvents = append(pastEvents, event)
		}
	}

	// Sort past events by most recent first
	sort.Slice(pastEvents, func(i, j int) bool {
		return pastEvents[i].StartDateTime.After(pastEvents[j].StartDateTime)
	})

	// Combine future and past events, with future events first
	events = append(futureEvents, pastEvents...)
	return events
}
