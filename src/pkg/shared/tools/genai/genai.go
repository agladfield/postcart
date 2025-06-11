// Package genai wraps different generative AI providers to create images
package genai

import (
	"context"
	"errors"
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/enum"
)

const genAIErrFmtStr = "genai err: %w"

func GenerateImage(ctx context.Context, genAIProvider enum.GenAIProviderEnum, prompt string) ([]byte, error) {
	var result []byte
	var genAIErr error

	switch genAIProvider {
	case enum.GenAIGoogleImagen4:
		result, genAIErr = newImagen4Image(ctx, prompt)
	default:
		genAIErr = errors.New("unknown generative ai image provider")
	}

	if genAIErr != nil {
		return nil, fmt.Errorf(genAIErrFmtStr, genAIErr)
	}

	return result, nil
}

// © Arthur Gladfield
