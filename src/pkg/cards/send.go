package cards

import (
	"os"

	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/tools/upload"
	"github.com/davidbyttow/govips/v2/vips"
)

type PostcardTemplateArguments struct {
	EncodedImage string `json:"encoded_image"`
	ASCIIText    string `json:"ascii_text"`
	Subject      string `json:"subject"`
}

func prepareEmail(unified *unifiedOutput, email *EmailParams) (*postmark.NewEmailFromTemplate[PostcardTemplateArguments], error) {
	// convert unified image ref to buffer to base64 encoding
	bytes, _, bytesErr := unified.UnifiedImage.ExportJpeg(&vips.JpegExportParams{
		Quality: 75,
	})
	if bytesErr != nil {
		return nil, bytesErr
	}

	uploadedURL, urlErr := upload.UploadImage(bytes)
	if urlErr != nil {
		return nil, urlErr
	}

	// fmt.Println(unified.unifiedText)

	os.WriteFile("./encoded.txt", []byte(uploadedURL), 0600)

	newEmail := postmark.NewEmailFromTemplate[PostcardTemplateArguments]{
		From:          "Postc.art Postcards <deliveries@postc.art>",
		To:            email.To.Email,
		TemplateAlias: "postcart-test-template",
		TemplateModel: PostcardTemplateArguments{
			EncodedImage: uploadedURL,
			ASCIIText:    unified.UnifiedText,
			Subject:      email.Subject,
		},
	}

	return &newEmail, nil
}

// func sendPostcard(){
// 	//
// }
