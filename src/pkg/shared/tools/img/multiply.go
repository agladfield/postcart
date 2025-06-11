package img

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/colors"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	imgMultiplyErrFmtStr = "img multiply err: %w"
	imgColorErrFmtStr    = "img color err: %w"
)

func Multiply(target *vips.ImageRef, multiplier *vips.ImageRef) (*vips.ImageRef, error) {
	toMult, copyErr := target.Copy()
	if copyErr != nil {
		return nil, fmt.Errorf(imgMultiplyErrFmtStr, copyErr)
	}
	multErr := toMult.Composite(multiplier, vips.BlendModeMultiply, 0, 0)
	if multErr != nil {
		return nil, fmt.Errorf(imgMultiplyErrFmtStr, multErr)
	}

	return toMult, nil
}

func Color(target *vips.ImageRef, color colors.Color) (*vips.ImageRef, error) {
	// Copy the target image to avoid modifying the original
	targetCopy, err := target.Copy()
	if err != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, fmt.Errorf("failed to copy target image: %w", err))
	}
	defer targetCopy.Close()

	// Get RGBA values from the input color (assuming 0-255 range)
	r, g, b, a := color.RGBA()
	// Normalize to [0, 255] if RGBA returns 0-65535 (adjust if needed)

	// Create a black image with the same dimensions as the target
	black, err := vips.Black(target.Width(), target.Height())
	if err != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, fmt.Errorf("failed to create black image: %w", err))
	}
	defer black.Close()

	// Copy the black image to ensure it’s modifiable
	colorImage, err := black.Copy()
	if err != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, fmt.Errorf("failed to copy black image: %w", err))
	}
	defer colorImage.Close()
	colorSpaceErr := colorImage.ToColorSpace(vips.InterpretationSRGB)
	if colorSpaceErr != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, colorSpaceErr)
	}

	// Set the pixel color using DrawRect
	vipsColor := vips.ColorRGBA{R: r, G: g, B: b, A: a}
	drawErr := colorImage.DrawRect(vipsColor, 0, 0, target.Width(), target.Height(), true)
	if drawErr != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, drawErr)
	}

	maskedColor, maskErr := Mask(colorImage, targetCopy)
	if maskErr != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, maskErr)
	}
	defer maskedColor.Close()

	multiplied, multErr := Multiply(targetCopy, maskedColor)
	if multErr != nil {
		return nil, fmt.Errorf(imgColorErrFmtStr, multErr)
	}

	return multiplied, nil
}

// © Arthur Gladfield
