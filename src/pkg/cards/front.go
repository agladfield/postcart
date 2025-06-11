package cards

import (
	"fmt"
	"strings"

	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/colors"
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	cardsFrontErrFmtStr     = "cards create front err: %w"
	cardsFrontTextErrFmtStr = "cards front text err: %w"
)

const (
	stampPosX         = 1057
	stampPosY         = 42
	stampResizeFactor = 169.
)

type textPos struct {
	w int
	h int
	x int
	y int
}

var frontTextPositions = struct {
	recipientName  textPos
	recipientEmail textPos
	senderName     textPos
	senderEmail    textPos
	message        textPos
}{
	recipientName: textPos{
		w: 360, h: 45, x: 88, y: 118,
	},
	recipientEmail: textPos{
		w: 360, h: 32, x: 88, y: 169,
	},
	senderName: textPos{
		w: 479, h: 51, x: 685, y: 543,
	},
	senderEmail: textPos{
		w: 479, h: 45, x: 685, y: 648,
	},
	message: textPos{
		w: 534, h: 552, x: 88, y: 227,
	},
}

func createFront(params *Params) (*sideOutput, error) {
	ascii := createFrontASCII(params)

	// create the stamp
	stampImg, stampErr := createStamp(params)
	if stampErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, stampErr)
	}
	defer stampImg.Close()

	contentImg, contentImgErr := pcCache.Obtain("res/postcards/content.png")
	if contentImgErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, contentImgErr)
	}
	defer contentImg.Close()
	fgMask, fgMaskErr := pcCache.Obtain("res/postcards/bg-frame.png")
	if fgMaskErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, fgMaskErr)
	}
	defer fgMask.Close()
	maskedContent, maskedContentErr := img.Mask(contentImg, fgMask)
	if maskedContentErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, maskedContentErr)
	}

	// apply writing
	textErr := addTextToFront(maskedContent, params)
	if textErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, textErr)
	}

	// apply texture
	if params.Textured == enum.TexturedEnabled {
		texture, textureErr := pcCache.Obtain("res/postcards/texture-classic.png")
		if textureErr != nil {
			return nil, fmt.Errorf(cardsFrontErrFmtStr, textureErr)
		}
		defer texture.Close()
		maskedTex, maskTexErr := img.Mask(texture, maskedContent)
		if maskTexErr != nil {
			return nil, fmt.Errorf(cardsFrontErrFmtStr, maskTexErr)
		}
		defer maskedTex.Close()

		var multErr error
		maskedContent, multErr = img.Multiply(maskedContent, maskedTex)
		if multErr != nil {
			return nil, fmt.Errorf(cardsFrontErrFmtStr, multErr)
		}
	}

	// // resize and apply stamp
	sizeRatio := stampResizeFactor / float64(stampImg.Width())
	resizeErr := stampImg.Resize(sizeRatio, vips.KernelAuto)
	if resizeErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, resizeErr)
	}
	compositeErr := maskedContent.Composite(stampImg, vips.BlendModeOver, stampPosX, stampPosY)
	if compositeErr != nil {
		return nil, fmt.Errorf(cardsFrontErrFmtStr, compositeErr)
	}

	frontBuff, _, backBuffErr := maskedContent.ExportPng(&vips.PngExportParams{Quality: 90})
	if backBuffErr != nil {
		return nil, backBuffErr
	}
	defer maskedContent.Close()
	front, frontErr := img.LoadFromBuffer(frontBuff)
	if frontErr != nil {
		return nil, frontErr
	}

	return &sideOutput{
		image: front,
		ascii: ascii,
	}, nil
}

func addTextToFront(content *vips.ImageRef, params *Params) error {
	fontEnum := params.Font
	if fontEnum == enum.FontUnknown {
		fontEnum = enum.FontMarker
	}
	font := getTextFont(fontEnum)

	recipientNameText := img.TextParams{
		Text:    params.To.Name,
		Font:    font,
		Width:   frontTextPositions.recipientName.w,
		Height:  frontTextPositions.recipientName.h,
		OffsetX: frontTextPositions.recipientName.x,
		OffsetY: frontTextPositions.recipientName.y,
		Color:   colors.Color{R: 0, G: 0, B: 0, A: 255},
	}

	recipNameErr := img.AddText(content, &recipientNameText)
	if recipNameErr != nil {
		return fmt.Errorf(cardsFrontTextErrFmtStr, recipNameErr)
	}

	recipientEmailText := img.TextParams{
		Text:    params.To.Email,
		Font:    font,
		Width:   frontTextPositions.recipientEmail.w,
		Height:  frontTextPositions.recipientEmail.h,
		OffsetX: frontTextPositions.recipientEmail.x,
		OffsetY: frontTextPositions.recipientEmail.y,
		Color:   colors.Color{R: 0, G: 0, B: 0, A: 255},
	}

	recipEmailErr := img.AddText(content, &recipientEmailText)
	if recipEmailErr != nil {
		return fmt.Errorf(cardsFrontTextErrFmtStr, recipEmailErr)
	}

	if len(params.From.Name) > 0 {
		senderNameText := img.TextParams{
			Text:    params.From.Name,
			Font:    font,
			Width:   frontTextPositions.senderName.w,
			Height:  frontTextPositions.senderName.h,
			OffsetX: frontTextPositions.senderName.x,
			OffsetY: frontTextPositions.senderName.y,
			Color:   colors.Color{R: 0, G: 0, B: 0, A: 255},
		}

		senderNameErr := img.AddText(content, &senderNameText)
		if senderNameErr != nil {
			return fmt.Errorf(cardsFrontTextErrFmtStr, senderNameErr)
		}
	}

	if len(params.From.Email) > 0 {
		senderEmailText := img.TextParams{
			Text:    params.From.Email,
			Font:    font,
			Width:   frontTextPositions.senderEmail.w,
			Height:  frontTextPositions.senderEmail.h,
			OffsetX: frontTextPositions.senderEmail.x,
			OffsetY: frontTextPositions.senderEmail.y,
			Color:   colors.Color{R: 0, G: 0, B: 0, A: 255},
		}

		senderEmailErr := img.AddText(content, &senderEmailText)
		if senderEmailErr != nil {
			return fmt.Errorf(cardsFrontTextErrFmtStr, senderEmailErr)
		}
	}

	if len(params.Message) > 0 {
		messageText := img.TextParams{
			Text:    messageBreakup(sanitizeMessage(params.Message)),
			Font:    font,
			Width:   frontTextPositions.message.w,
			Height:  frontTextPositions.message.h,
			OffsetX: frontTextPositions.message.x,
			OffsetY: frontTextPositions.message.y,
			Color:   colors.Color{R: 0, G: 0, B: 0, A: 255},
		}

		messageErr := img.AddText(content, &messageText)
		if messageErr != nil {
			return fmt.Errorf(cardsFrontTextErrFmtStr, messageErr)
		}
	}

	return nil
}

func sanitizeMessage(msg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(msg, "<", "‹"), ">", "›")
}

const (
	breakupThreshold = 20
)

// messageBreakup breaks up words in a message with hyphens,
// otherwise too long of words can cause improper message rendering
func messageBreakup(msg string) string {
	lines := strings.Split(msg, "\n")
	adjustedLines := []string{}
	for _, line := range lines {
		adjustedWords := []string{}
		words := strings.Fields(line)
		for _, word := range words {
			if len(word) <= breakupThreshold {
				adjustedWords = append(adjustedWords, word)
				continue
			}
			adjustedWord := ""
			for i := 0; i < len(word); i += breakupThreshold {
				end := i + breakupThreshold
				if end > len(word) {
					end = len(word)
				}
				// Add segment and hyphen if not at the end
				adjustedWord += word[i:end]
				if end < len(word) {
					adjustedWord += "-"
				}
			}
			adjustedWords = append(adjustedWords, adjustedWord)
		}
		adjustedLines = append(adjustedLines, strings.Join(adjustedWords, " "))
	}
	return strings.Join(adjustedLines, "\n")
}

// © Arthur Gladfield
