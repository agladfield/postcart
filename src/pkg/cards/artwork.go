package cards

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/agladfield/postcart/pkg/shared/tools/genai"
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	cardsArtworkErrFmtStr                 = "cards artwork err: %w"
	cardsArtworkPlaceholderErrFmtStr      = "cards artwork placeholder err: %w"
	cardsArtworkDecodeAttachmentErrFmtStr = "cards artwork decode attachment err: %w"
	cardsArtworkGenAIErrFmtStr            = "cards artwork gen ai err: %w"
	cardsArtworkPromptErrFmtStr           = "cards artwork prompt err: %w"
)

func getArtwork(ctx context.Context, artwork enum.ArtworkEnum, style enum.StyleEnum, attachment *postmark.EmailAttachment) (*vips.ImageRef, error) {
	var artworkBytes []byte
	var artworkErr error

	if artwork == enum.ArtworkAttachment {
		artworkBytes, artworkErr = decodeAttachment(attachment)
	} else {
		if env.UseAI() {
			artworkBytes, artworkErr = getGenAIArtwork(ctx, artwork, style)
		} else {
			artworkBytes, artworkErr = getPlaceholderArtwork(artwork)
		}
	}
	if artworkErr != nil {
		return nil, fmt.Errorf(cardsArtworkErrFmtStr, artworkErr)
	}

	artImg, loadErr := img.LoadFromBuffer(artworkBytes)
	if loadErr != nil {
		return nil, fmt.Errorf(cardsArtworkErrFmtStr, loadErr)
	}

	return artImg, nil
}

func validateContentType(contentType string) error {
	switch strings.ToLower(contentType) {
	case "image/jpeg":
		fallthrough
	case "image/png":
		fallthrough
	case "image/webp":
		return nil
	default:
		return fmt.Errorf("attachment content type unsupported: %s", contentType)
	}
}

func decodeAttachment(attachment *postmark.EmailAttachment) ([]byte, error) {
	// if jpeg, png, webp, or gif
	contentTypeInvalidErr := validateContentType(attachment.ContentType)
	if contentTypeInvalidErr != nil {
		return nil, fmt.Errorf(cardsArtworkDecodeAttachmentErrFmtStr, contentTypeInvalidErr)
	}
	bytes, decodeErr := base64.RawStdEncoding.DecodeString(attachment.Content)
	if decodeErr != nil {
		return nil, fmt.Errorf(cardsArtworkDecodeAttachmentErrFmtStr, decodeErr)
	}
	return bytes, nil
}

func getPlaceholderArtwork(artwork enum.ArtworkEnum) ([]byte, error) {
	var placeholderPath string
	switch artwork {
	case enum.ArtworkCity:
		placeholderPath = "res/artwork/city.png"
	case enum.ArtworkIslands:
		placeholderPath = "res/artwork/islands.png"
	case enum.ArtworkLakeside:
		placeholderPath = "res/artwork/lakeside.png"
	case enum.ArtworkMountains:
		fallthrough
	default:
		placeholderPath = "res/artwork/mountains.png"
	}

	bytes, readErr := placeholderArtwork.ReadFile(placeholderPath)
	if readErr != nil {
		return nil, fmt.Errorf(cardsArtworkPlaceholderErrFmtStr, readErr)
	}

	return bytes, nil
}

func getGenAIArtwork(ctx context.Context, artwork enum.ArtworkEnum, style enum.StyleEnum) ([]byte, error) {
	// get the prompt
	prompt, promptErr := createPrompt(artwork, style)
	if promptErr != nil {
		return nil, fmt.Errorf(cardsArtworkGenAIErrFmtStr, promptErr)
	}
	artworkBytes, artworkErr := genai.GenerateImage(ctx, enum.GenAIGoogleImagen4, prompt)
	if artworkErr != nil {
		return nil, fmt.Errorf(cardsArtworkGenAIErrFmtStr, artworkErr)
	}
	return artworkBytes, nil
}

var (
	errArtPromptNotFound   = errors.New("no prompt found")
	errStylePromptNotFound = errors.New("no style prompt found")
)

func createPrompt(artwork enum.ArtworkEnum, style enum.StyleEnum) (string, error) {
	artText, artTextExists := artworkPrompts[artwork]
	if !artTextExists {
		return "", fmt.Errorf(cardsArtworkPromptErrFmtStr, errArtPromptNotFound)
	}

	styleText, styleTextExists := stylePrompts[style]
	if !styleTextExists {
		return "", fmt.Errorf(cardsArtworkPromptErrFmtStr, errStylePromptNotFound)
	}

	prompt := fmt.Sprintf(basePrompt, styleText, artText)
	return prompt, nil
}

var artworkPrompts = map[enum.ArtworkEnum]string{
	enum.ArtworkLakeside:  "A large lake on the outskirts of a town, with either mountains or forest in the background. A few buildings sprinkled in the distance.",
	enum.ArtworkIslands:   "Tropical Islands on a nice sunny day. Water should be a focal point but we should have the perspective from a beach. We should also see the beach and some other tropical islands in the distance.",
	enum.ArtworkCity:      "A large size city with skyscrapers and a downtown on an average day. The city should depict a variety of different industries.",
	enum.ArtworkMountains: "A nice mountain range between seasons. There should be only a tiny amount of human intervention in the mountains, maybe one or two houses. There should be maybe a tiny a bit of snow.",
}

var stylePrompts = map[enum.StyleEnum]string{
	enum.StyleIllustrated:  "an illustrated post-war mid-century",
	enum.StylePhotograph:   "an ultra photo-realistic",
	enum.StyleVintagePhoto: "a mid 1960s - 1980s vintage photograph",
	enum.StylePainting:     "an oil painting with thick, heavy brush strokes",
}

const basePrompt = "You are to create a landscape strictly adhering to %s style that conforms to the following parameters: The scene should be of a %s. The focal content of the generated image should take up the full width and height of the frame."

// © Arthur Gladfield
