package img

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/colors"
	"github.com/davidbyttow/govips/v2/vips"
)

const imgAddTextErrFmtStr = "img add text err: %w"

type TextParams struct {
	Text    string
	Font    string
	Width   int
	Height  int
	OffsetX int
	OffsetY int
	Color   colors.Color
}

func AddText(target *vips.ImageRef, text *TextParams) error {
	var alpha *vips.ImageRef
	if target.Bands() == 4 {
		var alphaExtractErr error
		alpha, alphaExtractErr = target.ExtractBandToImage(3, 1)
		if alphaExtractErr != nil {
			return fmt.Errorf(imgAddTextErrFmtStr, alphaExtractErr)
		}
		extractErr := target.ExtractBand(0, 3)
		if extractErr != nil {
			return fmt.Errorf(imgAddTextErrFmtStr, extractErr)
		}
	}

	relWidth := float64(text.Width) / float64(target.Width())
	relHeight := float64(text.Height) / float64(target.Height())
	relX := float64(text.OffsetX) / float64(target.Width())
	relY := float64(text.OffsetY) / float64(target.Height())

	label := vips.LabelParams{
		Text: text.Text,
		Font: text.Font,
		Width: vips.Scalar{
			Value:    relWidth,
			Relative: true,
		},
		Height: vips.Scalar{
			Value:    relHeight,
			Relative: true,
		},
		OffsetX: vips.Scalar{
			Value:    relX,
			Relative: true,
		},
		OffsetY: vips.Scalar{
			Value:    relY,
			Relative: true,
		},
		Opacity: float32(text.Color.A) / 255,
		Color: vips.Color{
			R: text.Color.R,
			G: text.Color.G,
			B: text.Color.B,
		},
		Alignment: vips.AlignLow,
	}
	labelErr := target.Label(&label)
	if labelErr != nil {
		return fmt.Errorf(imgAddTextErrFmtStr, labelErr)
	}

	if alpha != nil {
		bandJoinErr := target.BandJoin(alpha)
		if bandJoinErr != nil {
			return fmt.Errorf(imgAddTextErrFmtStr, bandJoinErr)
		}
	}

	return nil
}

// © Arthur Gladfield
