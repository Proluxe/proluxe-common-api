package model

import (
	"fmt"
	"log"

	"github.com/scottraio/simpleforce"
)

type Meta struct {
	Terms        []Term        `json:"Terms"`
	Carriers     []Carrier     `json:"Carriers"`
	FOBCodes     []FOBCode     `json:"FOBCodes"`
	ShipMethods  []ShipMethod  `json:"ShipMethods"`
	FreightTerms []FreightTerm `json:"FreightTerms"`
	Reps         []Rep         `json:"Reps"`
}

// Exported types
type Carrier struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type FOBCode struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type ShipMethod struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	CarrierId string `json:"CarrierId"` // This is the lookup field to Carrier
}

type FreightTerm struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type Term struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type Rep struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

// Exported fetch functions
func FetchCarriers(client *simpleforce.Client) []Carrier {
	q := `
        SELECT Id, Name
        FROM rstk__sycarrier__c
    `
	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error fetching carriers: ", err)
	}
	var carriers []Carrier
	for _, record := range result.Records {
		carriers = append(carriers, Carrier{
			Id:   GetStringField("Id", record),
			Name: GetStringField("Name", record),
		})
	}
	return carriers
}

func FetchFOBCodes(client *simpleforce.Client, whereCondition string) []FOBCode {
	q := fmt.Sprintf(`
		SELECT Id, Name
		FROM rstk__syfob__c 
		WHERE %s
	`, whereCondition)
	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	var codes []FOBCode
	for _, record := range result.Records {
		codes = append(codes, FOBCode{
			Id:   GetStringField("Id", record),
			Name: GetStringField("Name", record),
		})
	}
	return codes
}

func FetchShipMethods(client *simpleforce.Client) []ShipMethod {
	q := `
		SELECT rstk__socarriervia_shipvia__c, Name, rstk__socarriervia_carrier__c
		FROM rstk__socarriervia__c
	`

	result, err := client.Query(q)
	if err != nil {
		log.Fatalf("Error fetching ship methods: %v", err)
	}
	var methods []ShipMethod
	for _, record := range result.Records {
		// Use the lookup field value as the Ship Method id
		methods = append(methods, ShipMethod{
			Id:        GetStringField("rstk__socarriervia_shipvia__c", record),
			Name:      GetStringField("Name", record),
			CarrierId: GetStringField("rstk__socarriervia_carrier__c", record),
		})
	}
	return methods
}

func FetchFreightTerms(client *simpleforce.Client, whereCondition string) []FreightTerm {
	q := fmt.Sprintf(`
		SELECT Id, Name
		FROM rstk__syfrghtrm__c 
		WHERE %s
	`, whereCondition)
	result, err := client.Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	var terms []FreightTerm
	for _, record := range result.Records {
		terms = append(terms, FreightTerm{
			Id:   GetStringField("Id", record),
			Name: GetStringField("Name", record),
		})
	}
	return terms
}

func FetchTerms(client *simpleforce.Client) []Term {
	q := `
		SELECT Id, Name
		FROM rstk__syterms__c
		ORDER BY Name ASC
	`

	result, err := client.Query(q)
	if err != nil {
		panic(err)
	}

	var terms []Term
	for _, record := range result.Records {
		t := Term{
			Id:   record["Id"].(string),
			Name: record["Name"].(string),
		}
		terms = append(terms, t)
	}

	return terms
}

func FetchReps(client *simpleforce.Client) []Rep {
	// Construct the query to fetch products by product ID
	q := fmt.Sprintf(`
		SELECT Id, Name, rstk__soregion_region__c, Customer_Master__r.Name
		FROM rstk__soregion__c
		WHERE Active__c = TRUE
		ORDER BY rstk__soregion_region__c ASC
	`)

	result, err := client.Query(q) // Note: for Tooling API, use client.Tooling().Query(q)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var reps []Rep
	for _, record := range result.Records {
		p := Rep{
			Id:   record["Id"].(string),
			Name: record["Name"].(string),
		}
		reps = append(reps, p)
	}

	return reps
}
