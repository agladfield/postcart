// Package img wraps govips to perform image and buffer operations
package img

import (
	"embed"
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const (
	imgLoadBuffErrFmtStr  = "img load from buffer err: %w"
	imgLoadEmbedErrFmtStr = "img load from embed err: %w"
	imgNewErrFmtStr       = "img new err: %w"
)

func LoadFromBuffer(buf []byte) (*vips.ImageRef, error) {
	buffRef, buffErr := vips.NewImageFromBuffer(buf)
	if buffErr != nil {
		return nil, fmt.Errorf(imgLoadBuffErrFmtStr, buffErr)
	}
	return buffRef, nil
}

func LoadFromEmbed(fs *embed.FS, path string) (*vips.ImageRef, error) {
	embededBytes, embedReadErr := fs.ReadFile(path)
	if embedReadErr != nil {
		return nil, fmt.Errorf(imgLoadEmbedErrFmtStr, embedReadErr)
	}

	embedRef, embedRefErr := vips.NewImageFromBuffer(embededBytes)
	if embedRefErr != nil {
		return nil, fmt.Errorf(imgLoadEmbedErrFmtStr, embedRefErr)
	}
	return embedRef, nil
}

func newTransparent(width, height int) (*vips.ImageRef, error) {
	// Create a new black image with RGB colorspace explicitly
	image, blackErr := vips.Black(width, height)
	if blackErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, blackErr)
	}

	// Ensure the image is in sRGB colorspace
	csErr := image.ToColorSpace(vips.InterpretationSRGB)
	if csErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, csErr)
	}

	// Add alpha channel (a single band for transparency)
	alphaErr := image.AddAlpha()
	if alphaErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, alphaErr)
	}

	// Set alpha channel to 0 (transparent)
	// This creates a fully transparent image with 4 bands (RGB + Alpha)
	extractErr := image.ExtractBand(3, 1)
	if extractErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, extractErr)
	}

	linearErr := image.Linear([]float64{0}, []float64{0})
	if linearErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, linearErr)
	}

	rgbImage, rgbBlackErr := vips.Black(width, height)
	if rgbBlackErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, rgbBlackErr)
	}

	rgbCSErr := rgbImage.ToColorSpace(vips.InterpretationSRGB)
	if rgbCSErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, rgbCSErr)
	}

	rgbJoinErr := rgbImage.BandJoin(image)
	if rgbJoinErr != nil {
		return nil, fmt.Errorf(imgNewErrFmtStr, rgbJoinErr)
	}

	return rgbImage, nil
}

func New(width, height int, transparent bool) (*vips.ImageRef, error) {
	if transparent {
		return newTransparent(width, height)
	} else {
		black, blackErr := vips.Black(width, height)
		if blackErr != nil {
			return nil, fmt.Errorf(imgNewErrFmtStr, blackErr)
		}
		return black, nil
	}
}

// © Arthur Gladfield
