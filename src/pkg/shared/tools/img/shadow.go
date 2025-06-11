package img

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const imgShadowErrFmtStr = "img shadow err: %w"

func Shadow(target *vips.ImageRef, opacity int) (*vips.ImageRef, error) {
	if opacity > 1 {
		opacity = 1
	}
	if opacity < 0 {
		opacity = 0
	}

	padded, padErr := Pad(target, [2]int{5, 5}, true)
	if padErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, padErr)
	}

	tc, copyErr := padded.Copy()
	if copyErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, copyErr)
	}

	if !tc.HasAlpha() {
		alphaErr := tc.AddAlpha()
		if alphaErr != nil {
			return nil, fmt.Errorf(imgShadowErrFmtStr, alphaErr)
		}
	}

	shadowAlpha, shadowAlphaErr := tc.ExtractBandToImage(3, 1)
	if shadowAlphaErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, shadowAlphaErr)
	}

	shadowRadius := 10
	// Create a black shadow with the target's alpha channel
	floatSizeW := float64(shadowAlpha.Width()) * 1.041
	floatSizeH := float64(shadowAlpha.Height()) * 1.041

	embedErr := shadowAlpha.Embed(shadowRadius, shadowRadius, int(floatSizeW), int(floatSizeH), vips.ExtendBackground)
	if embedErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, embedErr)
	}

	blurErr := shadowAlpha.GaussianBlur(float64(shadowRadius))
	if blurErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, blurErr)
	}

	shadowRGB, shadowRGBErr := New(shadowAlpha.Width(), shadowAlpha.Height(), false)
	if shadowRGBErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, shadowRGBErr)
	}

	shadowErr := shadowRGB.BandJoin(shadowAlpha)
	if shadowErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, shadowErr)
	}

	csErr := shadowRGB.ToColorSpace(vips.InterpretationSRGB)
	if csErr != nil {
		return nil, fmt.Errorf(imgShadowErrFmtStr, csErr)
	}

	return shadowRGB, nil
}

// © Arthur Gladfield
