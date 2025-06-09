package img

import (
	"github.com/davidbyttow/govips/v2/vips"
)

func Shadow(target *vips.ImageRef, opacity int) (*vips.ImageRef, error) {
	if opacity > 1 {
		opacity = 1
	}
	if opacity < 0 {
		opacity = 0
	}

	padded, padErr := Pad(target, [2]int{5, 5}, true)
	if padErr != nil {
		return nil, padErr
	}

	tc, copyErr := padded.Copy()
	if copyErr != nil {
		return nil, copyErr
	}

	if !tc.HasAlpha() {
		alphaErr := tc.AddAlpha()
		if alphaErr != nil {
			return nil, alphaErr
		}
	}

	shadowAlpha, shadowAlphaErr := tc.ExtractBandToImage(3, 1)
	if shadowAlphaErr != nil {
		return nil, shadowAlphaErr
	}

	shadowRadius := 10
	// Create a black shadow with the target's alpha channel
	floatSizeW := float64(shadowAlpha.Width()) * 1.041
	floatSizeH := float64(shadowAlpha.Height()) * 1.041

	embedErr := shadowAlpha.Embed(shadowRadius, shadowRadius, int(floatSizeW), int(floatSizeH), vips.ExtendBackground)
	if embedErr != nil {
		return nil, embedErr
	}

	blurErr := shadowAlpha.GaussianBlur(float64(shadowRadius))
	if blurErr != nil {
		return nil, blurErr
	}

	shadowRGB, shadowRGBErr := New(shadowAlpha.Width(), shadowAlpha.Height(), false)
	if shadowRGBErr != nil {
		return nil, shadowRGBErr
	}

	shadowErr := shadowRGB.BandJoin(shadowAlpha)
	if shadowErr != nil {
		return nil, shadowErr
	}

	csErr := shadowRGB.ToColorSpace(vips.InterpretationSRGB)
	if csErr != nil {
		return nil, csErr
	}

	DebugSave(shadowRGB, "./shadow.png")

	// shadowColor := colors.Color{R: 0, G: 0, B: 0, A: uint8(255 * opacity)}

	// shadowImg, err := Color(tc, shadowColor)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create shadow color: %v", err)
	// }

	// // Apply Gaussian blur to soften the shadow
	// sigma := 100.0 // Adjust blur radius as needed (smaller than 1000 for subtle effect)
	// err = shadowImg.GaussianBlur(sigma, 100)
	// if err != nil {
	// 	shadowImg.Close()
	// 	return nil, fmt.Errorf("failed to apply Gaussian blur: %v", err)
	// }

	// // Optionally offset the shadow (e.g., 5 pixels right and down)
	// err = shadowImg.Embed(5, 5, target.Width()+10, target.Height()+10, vips.ExtendBackground)
	// if err != nil {
	// 	shadowImg.Close()
	// 	return nil, fmt.Errorf("failed to offset shadow: %v", err)
	// }

	return shadowRGB, nil
}
