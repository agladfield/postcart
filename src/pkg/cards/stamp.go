package cards

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/agladfield/postcart/pkg/shared/tools/random"
	"github.com/davidbyttow/govips/v2/vips"
)

var stampShapePadding = map[enum.StampShapeEnum][2]int{
	enum.StampShapeRect:          {38, 38},
	enum.StampShapeRectClassic:   {40, 40},
	enum.StampShapeCircle:        {34, 34},
	enum.StampShapeCircleClassic: {46, 46},
}

var stampShapePaths = map[enum.StampShapeEnum]string{
	enum.StampShapeRect:          "stamp_rect.png",
	enum.StampShapeRectClassic:   "stamp_classic.png",
	enum.StampShapeCircle:        "stamp_circle.png",
	enum.StampShapeCircleClassic: "stamp_circle_classic.png",
}

var acceptableStampShapes = []enum.StampShapeEnum{
	enum.StampShapeRect,
	enum.StampShapeRectClassic,
	enum.StampShapeCircle,
	enum.StampShapeCircleClassic,
}

const (
	cardsStampErrFmtStr         = "cards create stamp err: %w"
	cardsStampGetShapeErrFmtStr = "cards stamp get shape err: %w"
)

func getStampShape(stampShape enum.StampShapeEnum) (*vips.ImageRef, error) {
	stampShapePath := fmt.Sprintf("res/stamps/%s", stampShapePaths[stampShape])
	loadedShape, loadErr := img.LoadFromEmbed(&stampResources, stampShapePath)
	if loadErr != nil {
		return nil, fmt.Errorf(cardsStampGetShapeErrFmtStr, loadErr)
	}

	return loadedShape, nil
}

func createStamp(email *EmailParams) (*vips.ImageRef, error) {
	shapeEnum := email.StampShape
	if shapeEnum == enum.StampShapeUnknown {
		shapeEnum = random.FromSlice(acceptableStampShapes)
	}

	stampShape, stampShapeErr := getStampShape(shapeEnum)
	if stampShapeErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, stampShapeErr)
	}
	defer stampShape.Close()

	stampTex, stampTexErr := img.LoadFromEmbed(&stampResources, "res/stamps/stamp-tex-1.jpg")
	if stampTexErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, stampTexErr)
	}
	maskedTex, maskTexErr := img.Mask(stampTex, stampShape)
	if maskTexErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, maskTexErr)
	}
	var stampMultErr error
	stampShape, stampMultErr = img.Multiply(stampShape, maskedTex)
	if stampMultErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, stampMultErr)
	}

	countryImg, countryErr := getCountryFlagImage(email.Country, shapeEnum)
	if countryErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, countryErr)
	}
	defer countryImg.Close()

	// get the country image
	shapePadding := stampShapePadding[shapeEnum]
	paddedCountry, padErr := img.Pad(countryImg, shapePadding, false)
	if padErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, padErr)
	}
	defer paddedCountry.Close()

	flattenedStamp, flattenErr := img.Flatten(stampShape, paddedCountry)
	if flattenErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, flattenErr)
	}

	stampShadow, shadowErr := img.Shadow(flattenedStamp, 200)
	if shadowErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, shadowErr)
	}

	paddedFlattened, padFlatErr := img.Pad(flattenedStamp, [2]int{18, 18}, false)
	if padFlatErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, padFlatErr)
	}

	compositeErr := stampShadow.Composite(paddedFlattened, vips.BlendModeOver, 9, 9)
	if compositeErr != nil {
		return nil, fmt.Errorf(cardsStampErrFmtStr, compositeErr)
	}

	return stampShadow, nil
}
