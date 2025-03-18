package google

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type User struct {
	PrimaryEmail string `json:"PrimaryEmail"`
	FullName     string `json:"FullName"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	OrgUnitPath  string `json:"OrgUnitPath"`
}

// GetAllUsers retrieves all users from GSuite (Google Workspace)
func GetAllUsers() ([]User, error) {
	// Context for API calls
	ctx := context.Background()

	fmt.Printf("Retrieving all users from GSuite\n")

	// Create Admin SDK Directory API service with impersonation
	directorySrv, err := createImpersonatedDirectoryService(ctx, "scott@proluxe.com") // Replace with the admin email
	if err != nil {
		return nil, fmt.Errorf("unable to create directory service: %v", err)
	}

	fmt.Printf("Created directory service\n")

	var allUsers []User
	pageToken := ""
	for {
		req := directorySrv.Users.List().Customer("my_customer").MaxResults(200)
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		// Execute the request
		users, err := req.Do()
		if err != nil {
			log.Printf("Error retrieving users: %v", err)
			return nil, err
		}

		fmt.Printf("Retrieved %d users\n", len(users.Users))

		// Create an array of User
		var usersList []User
		for _, u := range users.Users {

			user := User{
				PrimaryEmail: u.PrimaryEmail,
				FullName:     u.Name.FullName,
				FirstName:    u.Name.GivenName,
				LastName:     u.Name.FamilyName,
				OrgUnitPath:  u.OrgUnitPath,
			}
			usersList = append(usersList, user)
		}
		// Append the retrieved users to the list
		allUsers = append(allUsers, usersList...)

		// If there's no next page token, break the loop
		if users.NextPageToken == "" {
			break
		}
		pageToken = users.NextPageToken
	}

	return allUsers, nil
}

// createImpersonatedDirectoryService creates a Google Admin SDK Directory API client with impersonation using the service account credentials
func createImpersonatedDirectoryService(ctx context.Context, adminEmail string) (*admin.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to find default credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds.JSON, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Admin SDK client: %v", err)
	}

	config.Subject = adminEmail

	tokenSource := config.TokenSource(ctx)
	adminSrv, err := admin.NewService(ctx, option.WithTokenSource(tokenSource))
	return adminSrv, err
}
