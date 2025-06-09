package cards

import (
	"context"

	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/davidbyttow/govips/v2/vips"
)

func Parse(inbound *postmark.InboundData) (EmailParams, error) {
	// subject we use whatever

	details, err := parseEmailBody(inbound.TextBody)
	if err != nil {
		return EmailParams{}, err
	}

	details.ID = inbound.MessageID
	details.Subject = inbound.Subject

	if len(inbound.Attachments) > 0 {
		details.Attachment = &inbound.Attachments[0]
	}

	return details, nil
}

func assignUnknownValues(email *EmailParams) {
	if email.Border == enum.BorderUnknown {
		email.Border = enum.BorderStandard
	}
}

type sideOutput struct {
	ascii string
	image *vips.ImageRef
}

func Create(ctx context.Context, email *EmailParams) (*unifiedOutput, error) {
	assignUnknownValues(email)

	// Create the back artwork
	back, backErr := createBack(ctx, email)
	if backErr != nil {
		return nil, backErr
	}
	defer back.image.Close()

	// Create the front content
	// create front should have imageref, text contents
	front, frontErr := createFront(email)
	if frontErr != nil {
		return nil, frontErr
	}
	defer front.image.Close()

	bordered, bordersErr := addBorders(front, back, email.Border, email.Country)
	if bordersErr != nil {
		return nil, bordersErr
	}
	defer bordered.front.image.Close()
	defer bordered.back.image.Close()

	// create unified should return unified image ref, joined ascii art
	unified, unifyErr := unify(bordered)
	if unifyErr != nil {
		return nil, unifyErr
	}
	return unified, nil
}

func processJob(ctx context.Context, email EmailParams) error {
	// fmt.Println("processing:", email)
	unified, createErr := Create(ctx, &email)
	if createErr != nil {
		return createErr
	}
	defer unified.UnifiedImage.Close()

	// create email to send should return postmark email
	preparedEmail, prepareEmailErr := prepareEmail(unified, &email)
	if prepareEmailErr != nil {
		return prepareEmailErr
	}
	// // send email
	_, emailErr := postmark.SendWithTemplate(*preparedEmail)
	if emailErr != nil {
		return emailErr
	}
	// to check success measure emailRes.ErrorCode == 0
	// if success, record job as done, delete queued job record using transactions
	return nil
}
