package cards

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/agladfield/postcart/pkg/shared/tools/upload"
)

const cardsPrepareEmailErrFmtStr = "cards prepare email err: %w"

type PostcardTemplateArguments struct {
	ImageURL  string `json:"image_url"`
	ASCIIText string `json:"ascii_text"`
	Subject   string `json:"subject"`
}

func prepareEmail(unified *unifiedOutput, params *Params) (*postmark.NewEmailFromTemplate[PostcardTemplateArguments], error) {
	bytes, _, nativeErr := unified.UnifiedImage.ExportNative()
	if nativeErr != nil {
		return nil, fmt.Errorf(cardsPrepareEmailErrFmtStr, nativeErr)
	}

	uploadedURL, urlErr := upload.UploadImage(bytes)
	if urlErr != nil {
		return nil, fmt.Errorf(cardsPrepareEmailErrFmtStr, urlErr)
	}

	var subjectString string
	if params.From.Name != "Anonymous" && params.From.Name != "" {
		subjectString = fmt.Sprintf("%s sent you a postcard: %s", params.From.Name, params.Subject)
	} else {
		subjectString = fmt.Sprintf("You have a new postcard: %s", params.Subject)
	}

	newEmail := postmark.NewEmailFromTemplate[PostcardTemplateArguments]{
		From:          fmt.Sprintf("Postcards <deliveries@%s>", env.PostmarkEmailDomain()),
		To:            params.To.Email,
		TemplateAlias: deliveriesTemplateAlias,
		TemplateModel: PostcardTemplateArguments{
			ImageURL:  uploadedURL,
			ASCIIText: unified.UnifiedText,
			Subject:   subjectString,
		},
		Metadata: map[string]string{
			"sender_email": params.From.Email,
		},
	}

	return &newEmail, nil
}

// © Arthur Gladfield
