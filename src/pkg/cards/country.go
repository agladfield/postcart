package cards

import (
	"embed"
	"fmt"
	"strings"

	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/colors"
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	squarePathFmtStr = "res/countries/square/%s.png"
	circlePathFmtStr = "res/countries/circle/%s.png"
)

const cardsCountryErrFmtStr = "cards country err: %w"

func getCountryFlagImage(iso2 string, shape enum.StampShapeEnum) (*vips.ImageRef, error) {
	var embeddedFS *embed.FS
	var path string
	if shape.IsCircular() {
		embeddedFS = &circleFlags
		path = fmt.Sprintf(circlePathFmtStr, strings.ToLower(iso2))
	} else {
		embeddedFS = &squareFlags
		path = fmt.Sprintf(squarePathFmtStr, strings.ToLower(iso2))
	}

	countryImg, loadErr := img.LoadFromEmbed(embeddedFS, path)
	if loadErr != nil {
		return nil, fmt.Errorf(cardsCountryErrFmtStr, loadErr)
	}

	return countryImg, nil
}

func getCountryColors(iso2 string) ([]colors.Color, error) {
	countryFlag, countryFlagErr := squareFlags.Open(fmt.Sprintf(squarePathFmtStr, strings.ToLower(iso2)))
	if countryFlagErr != nil {
		return nil, fmt.Errorf(cardsCountryErrFmtStr, countryFlagErr)
	}
	defer countryFlag.Close()

	colors, colorErr := colors.ExtractColors(countryFlag)
	if colorErr != nil {
		return nil, fmt.Errorf(cardsCountryErrFmtStr, colorErr)
	}

	return colors, nil
}
