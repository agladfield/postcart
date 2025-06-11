package postmark

import (
	"fmt"
	"net/http"
)

const emailWithTemplatePath = "/email/withTemplate"

type NewEmailFromTemplate[T any] struct {
	// From format should be "Johnny Appleseed <johnny@apples.com>"
	From          string            `json:"From"`
	To            string            `json:"To"`
	TemplateID    int               `json:"TemplateId,omitempty"`
	TemplateAlias string            `json:"TemplateAlias,omitempty"` // can be used as alternative to TemplateID
	TemplateModel T                 `json:"TemplateModel"`
	Cc            string            `json:"Cc,omitempty"`
	Bcc           string            `json:"Bcc,omitempty"`
	Tag           string            `json:"Tag,omitempty"`
	ReplyTo       string            `json:"ReplyTo,omitempty"`
	Headers       *emailHeader      `json:"Headers,omitempty"`
	TrackOpens    bool              `json:"TrackOpens,omitempty"`
	Attachments   []EmailAttachment `json:"Attachments,omitempty"`
	Metadata      map[string]string `json:"Metadata,omitempty"`
	MessageStream string            `json:"MessageStream,omitempty"`
}

type EmailResponse struct {
	To          string `json:"To"`
	SubmittedAt string `json:"SubmittedAt"`
	MessageID   string `json:"MessageID"`
	ErrorCode   int    `json:"ErrorCode,omitempty"`
	Message     string `json:"Message"`
}

func SendWithTemplate[T any](email NewEmailFromTemplate[T]) (*EmailResponse, error) {
	body, bodyErr := EncodeToStruct(email)
	if bodyErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, bodyErr)
	}

	var emailRes EmailResponse

	reqErr := request(api, emailWithTemplatePath, http.MethodPost, body, &emailRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}
	if emailRes.ErrorCode > 0 {
		return nil, errorWithMessage(emailRes.ErrorCode, emailRes.Message)
	}

	return nil, nil
}
