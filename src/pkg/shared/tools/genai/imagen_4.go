package genai

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/auth/credentials"
	"github.com/agladfield/postcart/pkg/shared/env"
	"google.golang.org/genai"
)

const (
	imagen4Model = "imagen-4.0-generate-preview-05-20"
	negPrompt    = "text, words, letters, characters, border, frame, country, color, hexcode, #, letter, word, countries, name, color code, signature, watermark, username, people, face, children, adults, kids, boat, borders, side, frames, borderline, boundary, perimeter"
)

const (
	imagen4ErrFmtStr = "imagen4 err: %w"
)

var imageGen4Client *genai.Client

func newImagen4Image(ctx context.Context, prompt string) ([]byte, error) {
	if imageGen4Client == nil {
		// Load service account credentials
		data, readKeyErr := os.ReadFile(env.GCPCredsPath())
		if readKeyErr != nil {
			return nil, fmt.Errorf(imagen4ErrFmtStr, readKeyErr)
		}
		creds, credErr := credentials.DetectDefault(&credentials.DetectOptions{
			CredentialsJSON: data,
			Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		})
		if credErr != nil {
			return nil, fmt.Errorf(imagen4ErrFmtStr, credErr)
		}

		var clientErr error
		imageGen4Client, clientErr = genai.NewClient(ctx, &genai.ClientConfig{
			Location:    "us-central1",
			Project:     env.GCPProject(),
			Backend:     genai.BackendVertexAI,
			Credentials: creds,
		})
		if clientErr != nil {
			return nil, fmt.Errorf(imagen4ErrFmtStr, clientErr)
		}
	}

	res, genErr := imageGen4Client.Models.GenerateImages(
		ctx, imagen4Model,
		prompt,
		&genai.GenerateImagesConfig{
			IncludeRAIReason:        true,
			IncludeSafetyAttributes: true,
			OutputMIMEType:          "image/jpeg",
			AspectRatio:             "4:3",
			NumberOfImages:          1,
			NegativePrompt:          negPrompt,
			PersonGeneration:        genai.PersonGenerationAllowAll,
			EnhancePrompt:           false,
		},
	)
	if genErr != nil {
		return nil, fmt.Errorf(imagen4ErrFmtStr, genErr)
	}

	return res.GeneratedImages[0].Image.ImageBytes, nil
}

// © Arthur Gladfield
