// Package splash is splash artwork tool I wrote for having ascii text show up at the launch of a program
package splash

import (
	"strings"
)

type stringDimensions struct {
	height int
	width  int
}

func splashSW(s string) int {
	return len(s)
}

func getStringPrintWidthHeight(s string) stringDimensions {
	//
	theSplit := strings.Split(s, "\n")
	height := len(theSplit)

	width := 0

	for i := 0; i < len(theSplit); i++ {
		lw := splashSW(theSplit[i])
		if lw > width {
			width = lw
		}
	}

	return stringDimensions{width: width, height: height}
}

func isEven(n int) bool {
	return n%2 == 0
}

func centerTextBlob(
	txtBlob string,
	txtBlobWidth int,
	width int,
	_frameWidth int,
	_frameText string,
	_space string,
	_replaceWhitespace bool,
) string {
	if width <= txtBlobWidth {
		return txtBlob
	}

	frameWidth := 1
	if _frameWidth != 0 {
		frameWidth = _frameWidth
	}

	frameText := "#"
	if _frameText != "" {
		frameText = _frameText
	}

	space := " "
	if _space != "" {
		space = _space
	}

	replaceWhitespace := false
	if _replaceWhitespace {
		replaceWhitespace = _replaceWhitespace
	}

	startXSpace := int(float64(width/2 - txtBlobWidth/2))
	startX := startXSpace - frameWidth

	adjustedBlob := ""

	blobSplit := strings.Split(txtBlob, "\n")
	for i := 0; i < len(blobSplit); i++ {
		line := blobSplit[i]
		adjustedBlobLine := line

		frameBuff := strings.Repeat(frameText, frameWidth)

		adjustedBlobLine = frameBuff + strings.Repeat(space, startX) + adjustedBlobLine

		if frameText != "" {
			charLeft := width - splashSW(adjustedBlobLine)
			spaceLeft := charLeft - frameWidth*splashSW(frameText)
			adjustedBlobLine = adjustedBlobLine + strings.Repeat(space, spaceLeft) + frameBuff
		}

		if splashSW(adjustedBlob) != 0 {
			adjustedBlobLine = "\n" + adjustedBlobLine
		}

		if replaceWhitespace {
			adjustedBlobLine = strings.ReplaceAll(adjustedBlobLine, " ", space)
		}

		adjustedBlob += adjustedBlobLine
	}

	return adjustedBlob
}

func frameNewEmptyLine(
	width int,
	space string,
	_frameRepeat int,
	_frameText string,
	_ bool,
) string {
	//
	if _frameText == "" {
		return "\n"
	}

	gap := int(float64(width - _frameRepeat*2*splashSW(_frameText)))
	fr := strings.Repeat(_frameText, _frameRepeat)
	nl := fr + strings.Repeat(space, gap) + fr

	return nl
}

func frameNewLine(
	width int,
	frameText string,
	end bool,
) string {
	if frameText == "" {
		if end {
			return "\n"
		} else {
			return ""
		}
	}

	nl := strings.Repeat(frameText, int(float64(width/splashSW(frameText))))
	if end {
		nl += "\n"
	}

	return nl
}

type splashUsed struct {
	header bool
	center bool
	footer bool
}

type SplashContentOptions struct {
	Header string
	Center string
	Footer string
}

type SplashConfigOptions struct {
	BorderChar string
	RowThick   int
	ColThick   int
	SpaceChar  string
	Width      int
	Height     int
}

func Splash(content SplashContentOptions, config SplashConfigOptions) string {
	termWidth := 80
	termHeight := 24

	replaceWhitespace := false

	rowThick := 1
	colThick := 2
	borderChar := "#"
	spaceChar := " "

	if config.BorderChar != "" {
		borderChar = config.BorderChar
	}
	if config.RowThick != 0 {
		rowThick = config.RowThick
	}
	if config.ColThick != 0 {
		colThick = config.ColThick
	}
	if config.SpaceChar != "" {
		spaceChar = config.SpaceChar
	}
	if config.Width != 0 {
		termWidth = config.Width
	}
	if config.Height != 0 {
		termHeight = config.Height
	}

	header := content.Header
	center := content.Center
	footer := content.Footer

	headerSize := getStringPrintWidthHeight(header)
	headerBlob := centerTextBlob(header, headerSize.width, termWidth, colThick, borderChar, spaceChar, replaceWhitespace)
	centerSize := getStringPrintWidthHeight(center)
	centerBlob := centerTextBlob(center, centerSize.width, termWidth, colThick, borderChar, spaceChar, replaceWhitespace)
	footerSize := getStringPrintWidthHeight(footer)
	footerBlob := centerTextBlob(footer, footerSize.width, termWidth, colThick, borderChar, spaceChar, replaceWhitespace)

	toSubtractFramepad := rowThick * 2
	toSubtractFrameSpacing := rowThick * 2
	spacesAvailable := termHeight - toSubtractFramepad - toSubtractFrameSpacing

	isUsed := splashUsed{header: false, center: false, footer: false}

	if headerSize.width > 0 {
		spacesAvailable -= headerSize.height
		isUsed.header = true
	}

	if centerSize.width > 0 {
		spacesAvailable -= centerSize.height
		isUsed.center = true
	}

	if footerSize.width > 0 {
		spacesAvailable -= footerSize.height
		isUsed.footer = true
	}

	var toDistributeSpacesTo []string

	if isUsed.header && isUsed.center {
		toDistributeSpacesTo = append(toDistributeSpacesTo, "headerCenter")
	}
	if isUsed.center && isUsed.footer {
		toDistributeSpacesTo = append(toDistributeSpacesTo, "centerFooter")
	}
	if isUsed.header && !isUsed.center && isUsed.footer {
		toDistributeSpacesTo = append(toDistributeSpacesTo, "headerFooter")
	}
	if len(toDistributeSpacesTo) == 0 {
		toDistributeSpacesTo = append(toDistributeSpacesTo, "solo")
	}

	spaceCounter := map[string]int{
		"headerCenter": 0,
		"centerFooter": 0,
		"headerFooter": 0,
		"solo":         0,
	}

	spaceTargIndex := 0
	for space := spacesAvailable; space > 0; space-- {
		if len(toDistributeSpacesTo) == 0 {
			break
		}

		spaceCounter[toDistributeSpacesTo[spaceTargIndex]] += 1

		if spaceTargIndex == (len(toDistributeSpacesTo) - 1) {
			spaceTargIndex = 0
		} else {
			spaceTargIndex += 1
		}
	}

	nl := frameNewLine(termWidth, borderChar, false)
	enl := frameNewEmptyLine(termWidth, spaceChar, colThick, borderChar, replaceWhitespace)

	var printLines []string

	// Top Frame Padding
	for f := 0; f < rowThick; f++ {
		printLines = append(printLines, nl)
	}

	// Top Frame Spacing
	for s := 0; s < toSubtractFramepad/2; s++ {
		printLines = append(printLines, enl)
	}
	//
	if toDistributeSpacesTo[0] == "solo" {
		lineN := 0
		if isEven(spaceCounter["solo"]) {
			lineN = spaceCounter["solo"] / 2
		} else {
			lineN = int(float64(spaceCounter["solo"] / 2))
		}

		for s := 0; s < lineN; s++ {
			printLines = append(printLines, enl)
		}
	}

	if isUsed.header {
		if isUsed.header {
			printLines = append(printLines, headerBlob)
		}
		for hc := 0; hc < spaceCounter["headerCenter"]; hc++ {
			printLines = append(printLines, enl)
		}
	}

	// Between footer and header
	for hf := 0; hf < spaceCounter["headerFooter"]; hf++ {
		printLines = append(printLines, enl)
	}
	if isUsed.center {
		if isUsed.center {
			printLines = append(printLines, centerBlob)
		}
		for cf := 0; cf < spaceCounter["centerFooter"]; cf++ {
			printLines = append(printLines, enl)
		}
	}

	if isUsed.footer {
		printLines = append(printLines, footerBlob)
	}

	if toDistributeSpacesTo[0] == "solo" {
		lineN := 0
		if isEven(spaceCounter["solo"]) {
			lineN = spaceCounter["solo"] / 2
		} else {
			lineN = int(float64(spaceCounter["solo"]/2)) + 1
		}

		for s := 0; s < lineN; s++ {
			printLines = append(printLines, enl)
		}
	}

	//
	// Bottom Frame Spacing
	for s := 0; s < toSubtractFramepad/2; s++ {
		printLines = append(printLines, enl)
	}

	// Bottom Frame Padding
	for f := 0; f < rowThick; f++ {
		printLines = append(printLines, nl)
	}

	splashString := strings.Join(printLines, "\n")

	return splashString
}

// © Arthur Gladfield
