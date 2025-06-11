package cards

import (
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/agladfield/postcart/pkg/shared/tools/random"
)

var acceptableBorders = []enum.BorderEnum{
	enum.BorderStandard,
	enum.BorderStripes,
	enum.BorderLines,
	enum.BorderCubes,
	enum.BorderPhoto,
}

var acceptableArtwork = []enum.ArtworkEnum{
	enum.ArtworkCity,
	enum.ArtworkLakeside,
	enum.ArtworkIslands,
	enum.ArtworkMountains,
}

var acceptableStampShapes = []enum.StampShapeEnum{
	enum.StampShapeRect,
	enum.StampShapeRectClassic,
	enum.StampShapeCircle,
	enum.StampShapeCircleClassic,
}

var acceptableStyles = []enum.StyleEnum{
	enum.StylePainting,
	enum.StylePhotograph,
	enum.StyleVintagePhoto,
	enum.StyleIllustrated,
}

var acceptableFonts = []enum.FontEnum{
	enum.FontMarker,
	enum.FontPolite,
	enum.FontMidCentury,
	enum.FontTypewriter,
}

var acceptableTextured = []enum.TexturedEnum{
	enum.TexturedDisabled,
	enum.TexturedEnabled,
}

func assignUnknownValues(params *Params) {
	if params.Border == enum.BorderUnknown {
		params.Border = random.FromSlice(acceptableBorders)
	}
	if params.StampShape == enum.StampShapeUnknown {
		params.StampShape = random.FromSlice(acceptableStampShapes)
	}
	if params.Artwork == enum.ArtworkUnknown {
		params.Artwork = random.FromSlice(acceptableArtwork)
	}
	if params.Artwork == enum.ArtworkAttachment && !env.AllowAttachments() {
		params.Artwork = random.FromSlice(acceptableArtwork)
	}
	if params.Style == enum.StyleUnknown {
		params.Style = random.FromSlice(acceptableStyles)
	}
	if params.Font == enum.FontUnknown {
		params.Font = random.FromSlice(acceptableFonts)
	}
	if params.Textured == enum.TexturedUnknown {
		params.Textured = random.FromSlice(acceptableTextured)
	}
}

// © Arthur Gladfield
