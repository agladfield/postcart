package cards

import (
	"errors"
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/colors"
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

type borders struct {
	frontBorder *vips.ImageRef
	backBorder  *vips.ImageRef
}

type borderedOutput struct {
	front sideOutput
	back  sideOutput
}

const (
	cardsBorderErrFmtStr           = "cards border err: %w"
	cardsBorderCreateErrFmtStr     = "cards border create err: %w"
	cardsBorderStripesErrFmtStr    = "cards border stripes err: %w"
	cardsBorderLinesErrFmtStr      = "cards border lines err: %w"
	cardsBorderCubesErrFmtStr      = "cards border cubes err: %w"
	cardsBorderPhotoFrameErrFmtStr = "cards border photo frame err: %w"
)

func addBorders(front, back *sideOutput, border enum.BorderEnum, textured enum.TexturedEnum, country string) (*borderedOutput, error) {
	countryColors, countryColorsErr := getCountryColors(country)
	if countryColorsErr != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, countryColorsErr)
	}
	// create the borders
	borders, borderErrs := createBorders(border, textured, countryColors, back.image)
	if borderErrs != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, borderErrs)
	}

	// add them
	// flatten for front, multiply back content
	frontFlattened, flattenErr := img.Flatten(borders.frontBorder, front.image)
	if flattenErr != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, flattenErr)
	}

	// pad the back
	// we have to pad the back either way
	var updatedArtImage *vips.ImageRef
	if border != enum.BorderPhoto {
		paddedArtwork, padErr := img.Pad(back.image, [2]int{24, 36}, false)
		if padErr != nil {
			return nil, fmt.Errorf(cardsBorderErrFmtStr, padErr)
		}
		bgMask, bgMaskLoadErr := pcCache.Obtain("res/postcards/bg-frame.png")
		if bgMaskLoadErr != nil {
			return nil, fmt.Errorf(cardsBorderErrFmtStr, bgMaskLoadErr)
		}
		var bgMaskErr error
		updatedArtImage, bgMaskErr = img.Mask(paddedArtwork, bgMask)
		if bgMaskErr != nil {
			return nil, fmt.Errorf(cardsBorderErrFmtStr, bgMaskErr)
		}
	} else {
		var emptyErr error
		updatedArtImage, emptyErr = img.New(cardWidth, cardHeight, true)
		if emptyErr != nil {
			return nil, fmt.Errorf(cardsBorderErrFmtStr, emptyErr)
		}
		compositeErr := updatedArtImage.Composite(back.image, vips.BlendModeOver, 36, 36)
		if compositeErr != nil {
			return nil, fmt.Errorf(cardsBorderErrFmtStr, compositeErr)
		}
	}
	backMultiplied, multErr := img.Multiply(borders.backBorder, updatedArtImage)
	if multErr != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, multErr)
	}

	frontBuff, _, frontBuffErr := frontFlattened.ExportPng(&vips.PngExportParams{Quality: 90})
	if frontBuffErr != nil {
		return nil, frontBuffErr
	}
	defer frontFlattened.Close()
	frontImg, frontErr := img.LoadFromBuffer(frontBuff)
	if frontErr != nil {
		return nil, frontErr
	}

	backBuff, _, backBuffErr := backMultiplied.ExportPng(&vips.PngExportParams{Quality: 90})
	if backBuffErr != nil {
		return nil, backBuffErr
	}
	defer backMultiplied.Close()
	backImg, backErr := img.LoadFromBuffer(backBuff)
	if backErr != nil {
		return nil, backErr
	}

	frontSide := sideOutput{
		ascii: front.ascii,
		image: frontImg,
	}

	backSide := sideOutput{
		ascii: back.ascii,
		image: backImg,
	}

	return &borderedOutput{
		front: frontSide,
		back:  backSide,
	}, nil
}

func createBorders(border enum.BorderEnum, textured enum.TexturedEnum, countryColors []colors.Color, artImg *vips.ImageRef) (*borders, error) {
	baseImg, baseImgErr := pcCache.Obtain("res/postcards/rect.png")
	if baseImgErr != nil {
		return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, baseImgErr)
	}

	var borderErr error

	// create the front side border
	switch border {
	case enum.BorderStripes:
		borderErr = createFrameStripes(baseImg, countryColors)
	case enum.BorderPhoto:
		borderErr = createPhotoFrame(baseImg, artImg)
	case enum.BorderCubes:
		borderErr = createFrameCubes(baseImg, countryColors)
	case enum.BorderLines:
		borderErr = createFrameLines(baseImg, countryColors)
	default:
	}

	if borderErr != nil {
		return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, borderErr)
	}

	var frontBorder *vips.ImageRef
	var backBorder *vips.ImageRef

	frontBorder = baseImg

	if border != enum.BorderPhoto {
		frontCopy, frontCopyErr := frontBorder.Copy()
		if frontCopyErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, frontCopyErr)
		}
		flipErr := frontCopy.Flip(vips.DirectionHorizontal)
		if flipErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, flipErr)
		}
		backBorder = frontCopy
	} else {
		flatImg, flatImgErr := pcCache.Obtain("res/postcards/rect.png")
		if flatImgErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, flatImgErr)
		}
		backBorder = flatImg
	}

	// at this stage we apply textures to the borders if required
	if textured == enum.TexturedEnabled {
		texture, textureErr := pcCache.Obtain("res/postcards/texture-classic.png")
		if textureErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, textureErr)
		}
		maskedTex, maskTexErr := img.Mask(texture, frontBorder)
		if maskTexErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, maskTexErr)
		}
		var oldFrontBorder = frontBorder
		var multErr error
		frontBorder, multErr = img.Multiply(frontBorder, maskedTex)
		if multErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, multErr)
		}
		defer oldFrontBorder.Close()

		var oldBackBorder = backBorder
		flipErr := maskedTex.Flip(vips.DirectionHorizontal)
		if flipErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, flipErr)
		}
		backBorder, multErr = img.Multiply(backBorder, maskedTex)
		if multErr != nil {
			return nil, fmt.Errorf(cardsBorderCreateErrFmtStr, multErr)
		}
		defer oldBackBorder.Close()
	}

	_, toBytesErr := frontBorder.ToBytes()
	if toBytesErr != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, toBytesErr)
	}

	_, toBytesErr = backBorder.ToBytes()
	if toBytesErr != nil {
		return nil, fmt.Errorf(cardsBorderErrFmtStr, toBytesErr)
	}

	return &borders{
		frontBorder: frontBorder,
		backBorder:  backBorder,
	}, nil
}

func createFrameStripes(baseImg *vips.ImageRef, countryColors []colors.Color) error {
	stripesPrimary, stripesPrimErr := pcCache.Obtain("res/postcards/stripe1-primary.png")
	if stripesPrimErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, stripesPrimErr)
	}
	defer stripesPrimary.Close()

	stripesSecondary, stripesSecErr := pcCache.Obtain("res/postcards/stripe1-secondary.png")
	if stripesSecErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, stripesSecErr)
	}
	defer stripesSecondary.Close()

	primColor, secColor, colorsErr := colors.GetDesiredColors(countryColors)
	if colorsErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, colorsErr)
	}

	primColorImg, primColorErr := img.Color(stripesPrimary, primColor)
	if primColorErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, primColorErr)
	}

	secColorImg, secColorErr := img.Color(stripesSecondary, secColor)
	if secColorErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, secColorErr)
	}

	primCompErr := baseImg.Composite(primColorImg, vips.BlendModeOver, 0, 0)
	if primCompErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, primCompErr)
	}

	secCompErr := baseImg.Composite(secColorImg, vips.BlendModeOver, 0, 0)
	if secCompErr != nil {
		return fmt.Errorf(cardsBorderStripesErrFmtStr, secCompErr)
	}

	return nil
}

func createFrameLines(baseImg *vips.ImageRef, countryColors []colors.Color) error {
	linesPrimary, linesPrimErr := pcCache.Obtain("res/postcards/border-lines1.png")
	if linesPrimErr != nil {
		return fmt.Errorf(cardsBorderLinesErrFmtStr, linesPrimErr)
	}
	defer linesPrimary.Close()

	primColor, secColor, colorsErr := colors.GetDesiredColors(countryColors)
	if colorsErr != nil {
		return fmt.Errorf(cardsBorderLinesErrFmtStr, colorsErr)
	}

	primColorImg, primColorErr := img.Color(linesPrimary, primColor)
	if primColorErr != nil {
		return fmt.Errorf(cardsBorderLinesErrFmtStr, primColorErr)
	}

	primCompErr := baseImg.Composite(primColorImg, vips.BlendModeOver, 0, 0)
	if primCompErr != nil {
		return fmt.Errorf(cardsBorderLinesErrFmtStr, primCompErr)
	}

	if len(countryColors) > 1 {
		linesSecondary, linesSecErr := pcCache.Obtain("res/postcards/border-lines2.png")
		if linesSecErr != nil {
			return fmt.Errorf(cardsBorderLinesErrFmtStr, linesSecErr)
		}
		defer linesSecondary.Close()

		secColorImg, secColorErr := img.Color(linesSecondary, secColor)
		if secColorErr != nil {
			return fmt.Errorf(cardsBorderLinesErrFmtStr, secColorErr)
		}

		secCompErr := baseImg.Composite(secColorImg, vips.BlendModeOver, 0, 0)
		if secCompErr != nil {
			return fmt.Errorf(cardsBorderLinesErrFmtStr, secCompErr)
		}
	}

	return nil
}

func createFrameCubes(baseImg *vips.ImageRef, countryColors []colors.Color) error {
	cubesPrimary, cubesPrimErr := pcCache.Obtain("res/postcards/border-cubes1.png")
	if cubesPrimErr != nil {
		return fmt.Errorf(cardsBorderCubesErrFmtStr, cubesPrimErr)
	}
	defer cubesPrimary.Close()

	primColor, secColor, colorsErr := colors.GetDesiredColors(countryColors)
	if colorsErr != nil {
		return fmt.Errorf(cardsBorderCubesErrFmtStr, colorsErr)
	}

	primColorImg, primColorErr := img.Color(cubesPrimary, primColor)
	if primColorErr != nil {
		return fmt.Errorf(cardsBorderCubesErrFmtStr, primColorErr)
	}

	primCompErr := baseImg.Composite(primColorImg, vips.BlendModeOver, 0, 0)
	if primCompErr != nil {
		return fmt.Errorf(cardsBorderCubesErrFmtStr, primCompErr)
	}

	if len(countryColors) > 1 {
		cubesSecondary, cubesSecErr := pcCache.Obtain("res/postcards/border-cubes2.png")
		if cubesSecErr != nil {
			return fmt.Errorf(cardsBorderCubesErrFmtStr, cubesSecErr)
		}
		defer cubesSecondary.Close()

		secColorImg, secColorErr := img.Color(cubesSecondary, secColor)
		if secColorErr != nil {
			return fmt.Errorf(cardsBorderCubesErrFmtStr, secColorErr)
		}

		secCompErr := baseImg.Composite(secColorImg, vips.BlendModeOver, 0, 0)
		if secCompErr != nil {
			return fmt.Errorf(cardsBorderCubesErrFmtStr, secCompErr)
		}
	}

	return nil
}

var errBorderPhotoFrameNilArtwork = errors.New("cannot have a nil artwork image passed to borders")

func createPhotoFrame(baseImg, artImg *vips.ImageRef) error {
	if artImg == nil {
		return fmt.Errorf(cardsBorderPhotoFrameErrFmtStr, errBorderPhotoFrameNilArtwork)
	}

	bgFrameMask, bgFrameMaskErr := pcCache.Obtain("res/postcards/bg-frame-inverse.png")
	if bgFrameMaskErr != nil {
		return fmt.Errorf(cardsBorderPhotoFrameErrFmtStr, bgFrameMaskErr)
	}

	maskedFrame, frameArtErr := img.Mask(artImg, bgFrameMask)
	if frameArtErr != nil {
		return fmt.Errorf(cardsBorderPhotoFrameErrFmtStr, frameArtErr)
	}

	compositeErr := baseImg.Composite(maskedFrame, vips.BlendModeOver, 0, 0)
	if compositeErr != nil {
		return fmt.Errorf(cardsBorderPhotoFrameErrFmtStr, compositeErr)
	}

	cropWidth := artImg.Width() - 36*2
	cropHeight := artImg.Height() - 36*2

	cropErr := artImg.ExtractArea(36, 36, cropWidth, cropHeight)
	if cropErr != nil {
		return fmt.Errorf(cardsBorderPhotoFrameErrFmtStr, cropErr)
	}

	return nil
}

// © Arthur Gladfield
