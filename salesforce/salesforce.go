package salesforce

import (
	"fmt"
	"log"

	u "github.com/scottraio/go-utils"
	simpleforce "github.com/simpleforce/simpleforce"
)

type SF struct {
	Client     *simpleforce.Client
	QueryCount int
}

func NewSF() *SF {
	client, err := CreateClient()
	if err != nil {
		log.Println("Error creating Salesforce client: ", err)
		return nil
	}

	return &SF{Client: client}
}

func CreateClient() (*simpleforce.Client, error) {
	sfURL := u.GetDotEnvVariable("SF_URL")
	sfUser := u.GetDotEnvVariable("SF_USER")
	sfPassword := u.GetDotEnvVariable("SF_PASSWORD")
	sfToken := u.GetDotEnvVariable("SF_TOKEN")

	client := simpleforce.NewClient(sfURL, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)

	if client == nil {
		// handle the error
		log.Println("Something went wrong with creating the client.")
		return nil, fmt.Errorf("error creating client")
	}

	err := client.LoginPassword(sfUser, sfPassword, sfToken)

	// Do some other stuff with the client instance if needed.
	return client, err
}

func (sf *SF) Query(query string) (*simpleforce.QueryResult, error) {
	result, err := sf.Client.Query(query)
	if err != nil {
		log.Println("Error querying Salesforce: ", err)
		return nil, err
	}

	sf.QueryCount++

	return result, nil
}
