package img

import (
	"embed"

	"github.com/davidbyttow/govips/v2/vips"
)

func LoadFromBuffer(buf []byte) (*vips.ImageRef, error) {
	return vips.NewImageFromBuffer(buf)
}

func LoadFromEmbed(fs *embed.FS, path string) (*vips.ImageRef, error) {
	embededBytes, embedReadErr := fs.ReadFile(path)
	if embedReadErr != nil {
		return nil, embedReadErr
	}

	return vips.NewImageFromBuffer(embededBytes)
}

func newTransparent(width, height int) (*vips.ImageRef, error) {
	// Create a new black image with RGB colorspace explicitly
	image, err := vips.Black(width, height)
	if err != nil {
		return nil, err
	}

	// Ensure the image is in sRGB colorspace
	err = image.ToColorSpace(vips.InterpretationSRGB)
	if err != nil {
		return nil, err
	}

	// Add alpha channel (a single band for transparency)
	err = image.AddAlpha()
	if err != nil {
		return nil, err
	}

	// Set alpha channel to 0 (transparent)
	// This creates a fully transparent image with 4 bands (RGB + Alpha)
	err = image.ExtractBand(3, 1)
	if err != nil {
		return nil, err
	}

	err = image.Linear([]float64{0}, []float64{0})
	if err != nil {
		return nil, err
	}

	// Join back with the RGB image
	rgbImage, err := vips.Black(width, height)
	if err != nil {
		return nil, err
	}

	err = rgbImage.ToColorSpace(vips.InterpretationSRGB)
	if err != nil {
		return nil, err
	}

	err = rgbImage.BandJoin(image)
	if err != nil {
		return nil, err
	}

	return rgbImage, nil
}

func New(width, height int, transparent bool) (*vips.ImageRef, error) {
	if transparent {
		return newTransparent(width, height)
	} else {
		return vips.Black(width, height)
	}
}
