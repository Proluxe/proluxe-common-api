package services

import (
	"context"

	"github.com/knocklabs/knock-go/knock"
	u "github.com/scottraio/go-utils"
)

type Knock struct {
	User       *knock.User
	WorkFlowId string
	Email      string
}

func (k *Knock) client() (context.Context, *knock.Client) {
	token := u.GetDotEnvVariable("KNOCK_API_KEY")

	ctx := context.Background()

	// create a new Knock API client with the given access token
	client, _ := knock.NewClient(
		knock.WithAccessToken(token),
	)

	return ctx, client
}

func (k *Knock) Identify() *knock.User {
	ctx, client := k.client()

	user, _ := client.Users.Identify(ctx, &knock.IdentifyUserRequest{
		ID: k.Email,
	})

	k.User = user

	return user
}

func (k *Knock) Trigger(recipients []string, payload map[string]interface{}) error {
	ctx, client := k.client()

	req := &knock.TriggerWorkflowRequest{
		Workflow: k.WorkFlowId,
		Data:     payload,
		Actor:    k.User,
	}

	for _, recipient := range recipients {
		req.AddRecipientByID(recipient)
	}

	_, err := client.Workflows.Trigger(ctx, req, nil)

	return err
}
