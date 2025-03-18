package model

import (
	"log"
	"regexp"
	"time"

	"cloud.google.com/go/civil"
)

type Address struct {
	City          string `json:"City"`
	Country       string `json:"Country"`
	State         string `json:"State"`
	Street        string `json:"Street"`
	Street2       string `json:"Street2"`
	Zip           string `json:"Zip"`
	TaxLocation   string `json:"TaxLocation"`
	Email         string `json:"Email"`
	DefaultShipTo bool   `json:"DefaultShipTo"`
	DefaultBillTo bool   `json:"DefaultBillTo"`
}

func (a *Address) fullAddress() string {
	var street string

	if a.Street == "" {
		street = ""
	} else {
		street = a.Street + ", "
	}

	return street + a.City + ", " + a.State + " " + a.Zip + ", " + a.Country
}

func getStringField(field string, record map[string]interface{}) string {
	if record == nil {
		return ""
	}

	if value, ok := record[field]; ok && value != nil {
		return value.(string)
	}

	return ""
}

// Helper function to get an optional string field from a nested map
func getOptionalStringField(field string, nestedMap interface{}) string {
	if nestedMap == nil {
		return ""
	}
	if nestedMapMap, ok := nestedMap.(map[string]interface{}); ok {
		if value, exists := nestedMapMap[field]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}
	return ""
}

func getOptionalFloatField(field string, nestedMap interface{}) float64 {
	if nestedMap == nil {
		return 0
	}
	if nestedMapMap, ok := nestedMap.(map[string]interface{}); ok {
		if value, exists := nestedMapMap[field]; exists {
			if floatValue, ok := value.(float64); ok {
				return floatValue
			}
		}
	}
	return 0
}

// Helper function to extract float fields from the record
func getFloatField(fieldName string, record map[string]interface{}) float64 {
	if value, ok := record[fieldName]; ok && value != nil {
		return value.(float64)
	}
	return 0
}

func getIntField(fieldName string, record map[string]interface{}) int {
	if value, ok := record[fieldName]; ok && value != nil {
		return int(value.(float64))
	}
	return 0
}

func getBoolField(fieldName string, record map[string]interface{}) bool {
	if value, ok := record[fieldName]; ok && value != nil {
		return value.(bool)
	}
	return false
}

func getDateField(fieldName string, record map[string]interface{}) civil.Date {
	if value, ok := record[fieldName]; ok && value != nil {
		if dateStr, ok := value.(string); ok {
			// Try parsing with the new layout
			layout := "2006-01-02T15:04:05.000-0700"
			if parsedTime, err := time.Parse(layout, dateStr); err == nil {
				return civil.DateOf(parsedTime)
			}
			date, err := civil.ParseDate(dateStr)
			if err == nil {
				return date
			}
		}
	}

	// Return zero value if parsing fails
	return civil.Date{}
}

func convertToTime(date string) time.Time {
	layout := "2006-01-02T15:04:05.000-0700"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return t
}

func isSalesforceID(id string) bool {
	// Check if the ID is 15 or 18 characters long
	if len(id) != 15 && len(id) != 18 {
		return false
	}

	// Salesforce IDs only contain alphanumeric characters
	validIDPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !validIDPattern.MatchString(id) {
		return false
	}

	// If no prefix matches, it could still be a valid Salesforce ID, but without checking specific prefixes
	return true
}

func parseFloat(value interface{}) float64 {
	if value == nil {
		return 0
	}
	if floatValue, ok := value.(float64); ok {
		return floatValue
	}
	return 0
}

func convertToDate(date string) time.Time {
	if date == "" {
		return time.Time{}
	}

	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return t
}
