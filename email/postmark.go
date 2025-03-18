package email

import (
	"net/http"

	"github.com/mattevans/postmark-go"
	u "github.com/scottraio/go-utils"
)

type Postmark struct {
	Client *postmark.Client
}

func New() Postmark {
	return Postmark{
		Client: postmark.NewClient(
			postmark.WithClient(&http.Client{
				Transport: &postmark.AuthTransport{Token: u.GetDotEnvVariable("POSTMARK_API_TOKEN")},
			}),
		)}
}

func (p *Postmark) SendEmail(from, to string, templateId int, payload map[string]any) error {
	emailReq := &postmark.Email{
		From:          from,
		To:            to,
		TemplateID:    templateId,
		TemplateModel: payload,
		TrackOpens:    true,
	}

	_, _, err := p.Client.Email.Send(emailReq)
	if err != nil {
		return err
	}

	return nil
}
