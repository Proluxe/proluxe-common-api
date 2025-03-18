package google

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// MessageWithDate holds a Gmail message along with its internal date for sorting
type Message struct {
	Snippet      string
	InternalDate int64
	MessageId    string
	UserEmail    string
	Id           string
	From         string
	To           string
	Body         string
	Message      *gmail.Message
}

func GetEmails(email string) ([]Message, error) {
	leadEmail := email

	userEmails := []string{
		"mcole@proluxe.com",
		"scott@proluxe.com",
		"bbauer@proluxe.com",
		"sraio@proluxe.com",
		"emartinez@proluxe.com",
		"bmunguia@proluxe.com",
	}

	// Map to track unique message IDs and avoid duplicates
	uniqueMessages := make(map[string]Message)

	// Mutex to protect access to uniqueMessages when multiple goroutines update it
	var mu sync.Mutex

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Context for API calls
	ctx := context.Background()

	// For each user, impersonate them and search their Gmail
	for _, userEmail := range userEmails {
		wg.Add(1)

		// Launch a goroutine for each user
		go func(userEmail string) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			fmt.Printf("Searching emails for user: %s\n", userEmail)

			gmailSrv, err := createImpersonatedGmailService(ctx, userEmail)
			if err != nil {
				log.Printf("Unable to create Gmail client for user %s: %v", userEmail, err)
				return
			}

			// Query for emails involving the lead email
			query := fmt.Sprintf("\"%s\" OR from:%s OR to:%s OR cc:%s", leadEmail, leadEmail, leadEmail, leadEmail)
			messageList, err := gmailSrv.Users.Messages.List(userEmail).Q(query).MaxResults(100).Do()
			if err != nil {
				log.Printf("Error searching for lead email (%s) in %s's mailbox: %v", leadEmail, userEmail, err)
				return
			}

			if len(messageList.Messages) == 0 {
				log.Printf("No messages found for user: %s with query: %s", userEmail, query)
				return
			}

			for _, msg := range messageList.Messages {
				// Fetch the message details
				message, err := gmailSrv.Users.Messages.Get(userEmail, msg.Id).Do()
				if err != nil {
					log.Printf("Error fetching message details for message ID %s: %v", msg.Id, err)
					continue
				}

				mu.Lock() // Protect access to uniqueMessages with the mutex
				if _, exists := uniqueMessages[message.Id]; !exists {
					messageId := getHeader(message.Payload.Headers, "Message-ID")

					uniqueMessages[messageId] = Message{
						Snippet:      message.Snippet,
						InternalDate: message.InternalDate,
						MessageId:    messageId,
						UserEmail:    userEmail,
						From:         getHeader(message.Payload.Headers, "From"),
						To:           getHeader(message.Payload.Headers, "To"),
						Id:           message.Id,
						Body:         getBody(message),
					}
				}
				mu.Unlock() // Unlock after updating the map
			}
		}(userEmail)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Convert the uniqueMessages map to a slice and sort by InternalDate
	var sortedMessages []Message
	for _, msgWithDate := range uniqueMessages {
		sortedMessages = append(sortedMessages, msgWithDate)
	}

	// Sort the messages by InternalDate
	sort.Slice(sortedMessages, func(i, j int) bool {
		return sortedMessages[i].InternalDate < sortedMessages[j].InternalDate
	})

	return sortedMessages, nil
}

// GetEmailByID retrieves the details of a single email based on the message ID and user email
func GetEmailByID(userEmail string, messageId string) (*gmail.Message, error) {
	// Context for API calls
	ctx := context.Background()

	// Create a Gmail service with impersonation
	gmailSrv, err := createImpersonatedGmailService(ctx, userEmail)
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail client for user %s: %v", userEmail, err)
	}

	// Retrieve the email by message ID
	message, err := gmailSrv.Users.Messages.Get(userEmail, messageId).Do()
	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("error fetching message with ID %s: %v", messageId, err)
	}

	return message, nil
}

func getHeader(headers []*gmail.MessagePartHeader, key string) string {
	for _, header := range headers {
		if header.Name == key {
			return header.Value
		}
	}
	return ""
}

// createImpersonatedGmailService creates a Gmail API client with impersonation using the service account credentials
func createImpersonatedGmailService(ctx context.Context, userEmail string) (*gmail.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, gmail.MailGoogleComScope)
	if err != nil {
		return nil, fmt.Errorf("unable to find default credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds.JSON, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Admin SDK client: %v", err)
	}

	config.Subject = userEmail

	tokenSource := config.TokenSource(ctx)
	gmailSrv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	return gmailSrv, err
}

func getBody(message *gmail.Message) string {
	if message.Payload.Body.Data != "" {
		return message.Payload.Body.Data
	}

	if len(message.Payload.Parts) > 0 {
		for _, part := range message.Payload.Parts {
			if part.MimeType == "text/plain" {
				return part.Body.Data
			}
			if part.MimeType == "multipart/alternative" {
				for _, subPart := range part.Parts {
					if subPart.MimeType == "text/plain" {
						return subPart.Body.Data
					}
				}
			}
		}
	}

	return ""
}
