package cards

import (
	"context"
	"fmt"

	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	cardsJobErrFmtStr    = "cards job err: %w"
	cardsCreateErrFmtStr = "cards create err: %w"
)

type sideOutput struct {
	ascii string
	image *vips.ImageRef
}

func Create(ctx context.Context, params *Params) (*unifiedOutput, error) {
	assignUnknownValues(params)

	// Create the back artwork
	back, backErr := createBack(ctx, params)
	if backErr != nil {
		return nil, fmt.Errorf(cardsCreateErrFmtStr, backErr)
	}
	defer back.image.Close()

	// Create the front content
	// create front should have imageref, text contents
	front, frontErr := createFront(params)
	if frontErr != nil {
		return nil, fmt.Errorf(cardsCreateErrFmtStr, frontErr)
	}
	defer front.image.Close()

	bordered, bordersErr := addBorders(front, back, params.Border, params.Textured, params.Country)
	if bordersErr != nil {
		return nil, fmt.Errorf(cardsCreateErrFmtStr, bordersErr)
	}
	defer bordered.front.image.Close()
	defer bordered.back.image.Close()

	// create unified should return unified image ref, joined ascii art
	unified, unifyErr := unify(bordered)
	if unifyErr != nil {
		return nil, fmt.Errorf(cardsCreateErrFmtStr, unifyErr)
	}
	return unified, nil
}

func processJob(ctx context.Context, email Params) error {
	unified, createErr := Create(ctx, &email)
	if createErr != nil {
		return fmt.Errorf(cardsJobErrFmtStr, createErr)
	}
	defer unified.UnifiedImage.Close()

	// prepare email contents to send with postmark (returns postmark email with template)
	preparedEmail, prepareEmailErr := prepareEmail(unified, &email)
	if prepareEmailErr != nil {
		return fmt.Errorf(cardsJobErrFmtStr, prepareEmailErr)
	}
	// send email
	_, emailErr := postmark.SendWithTemplate(*preparedEmail)
	if emailErr != nil {
		return fmt.Errorf(cardsJobErrFmtStr, emailErr)
	}
	jdb.RecordSent()
	jdb.RemoveJobFromRecords(email.ID)
	if email.Attachment != nil {
		remErr := jdb.RemoveAttachmentForJob(email.ID)
		if remErr != nil {
			return fmt.Errorf(cardsJobErrFmtStr, remErr)
		}
	}

	return nil
}

// © Arthur Gladfield
