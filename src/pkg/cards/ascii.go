package cards

import (
	"fmt"
	"strings"

	"github.com/agladfield/postcart/pkg/shared/enum"
)

func createBackASCII(artwork enum.ArtworkEnum, border enum.BorderEnum) string {
	asciiArtwork := getBackASCIIArtwork(artwork)
	borderASCII := getASCIIBorder(border)
	return addASCIIToFrame(asciiArtwork, borderASCII)
}

func createFrontASCII(email *EmailParams) string {
	to := email.To
	from := email.From
	message := strings.ReplaceAll(email.Message, "\n", " ")
	country := email.Country
	border := email.Border
	shape := email.StampShape

	borderASCII := getASCIIBorder(border)
	leftSide := createFrontASCIILeftSide(from, message, len(borderASCII.left))
	rightSide := createFrontASCIIRightSide(to, shape, country, len(borderASCII.right))

	sideASCII := joinFrontASCIISides(leftSide, rightSide, borderASCII)
	top := strings.Repeat(borderASCII.top, asciiBorderWidth/len(borderASCII.top))
	bottom := strings.Repeat(borderASCII.bottom, asciiBorderWidth/len(borderASCII.bottom))

	// add message unsmushed
	msg := breakupMessageForASCII(message)

	return strings.Join([]string{top, sideASCII, bottom, "", msg}, "\n")
}

func getBackASCIIArtwork(artwork enum.ArtworkEnum) string {
	switch artwork {
	case enum.ArtworkIslands:
		return islandASCII
	case enum.ArtworkLakeside:
		return lakesideASCII
	case enum.ArtworkAttachment:
		return attachmentASCII
	case enum.ArtworkCity:
		return cityASCII
	case enum.ArtworkMountains:
		fallthrough
	default:
		return mountainsASCII
	}
}

const lakesideASCII = `     .      .     __
    / \    / \   /  \
  / /  \ /  \  \   /  \
/_______. ~   | .______\
-_~_--___\_____/--___---
\__-__--_--_--__--_----/
 \--__--_--___~~_--_~_/
  --~--#---~----~--#--
   \~--~--#--~-#--~-/
    \__##_____#___/`

const islandASCII = `
      /\/\/\/\
     //\/\/\/\\
      '_\V/_'
         #
         #
         #
         #.a@@a.
       .aa@@@@@@@@@a
    .a@@@@@@@@@@@@@@@@@@aa.
    ~~~~~~~~~~~~~~~~~~~~~~~`

const mountainsASCII = `
           /\
          /  \/\
         /   /  \
        /\_/\_/\/\  /\
       /          \/ -\
      /  /     \   \ \ \
     /    /       \ \   \
    /-___--_-___-_-__\###\
    #######################`

const attachmentASCII = `                       ___
        _______________   /   \
       |         /\    | |  __
       |    __  /__\   | | |  \
       |   |__|   /\   | | |  |
       |          \/   | | |  |
       |_______________| | |  |
                          \__/
    `

const cityASCII = `                           __|__
               ~          /  O  \
             ~   ______  |  | |  |
  __|__     ~   | City | | # # # |
  |+|+|  _/|___ |______| |  | |  |
  |+|+| |~~~[]~|   ||    | # _ # |
  ||_|| |~H~~~~|   []    |__|_|__|
  --------------------------------
  - -- -- -- -- -- -- -- -- -- --
  --------------------------------`

type asciiBorder struct {
	top    string
	left   string
	right  string
	bottom string
}

var linesBorder = asciiBorder{
	left:   "||",
	right:  "-",
	bottom: "-",
	top:    "||",
}

var jaggedBorder = asciiBorder{
	left:   "<",
	right:  ">",
	bottom: "\\/",
	top:    "/\\",
}
var stripedBorder = asciiBorder{
	left:   "\\",
	right:  "\\",
	bottom: "\\ \\ \\",
	top:    "\\ \\ \\",
}
var photoBorder = asciiBorder{
	left:   "[]",
	right:  "[]",
	bottom: "[]",
	top:    "[]",
}
var defaultBorder = asciiBorder{
	left:   "#",
	right:  "#",
	bottom: "#",
	top:    "#",
}

const (
	asciiBorderWidth  = 40
	asciiBorderHeight = 15
)

func addASCIIToFrame(ascii string, border asciiBorder) string {
	lines := strings.Split(ascii, "\n")
	height := len(lines)

	maxWidth := 0

	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	topBorderSeq := border.top
	leftBorderSeq := border.left
	rightBorderSeq := border.right
	bottomBorderSeq := border.bottom

	xPadSpaces := (38 - maxWidth) / 2
	yPadSpaces := (asciiBorderHeight - 2 - height) / 2

	borderedASCIILines := []string{}
	spaceYCount := 0
	asciiPos := 0
	asciiDone := false

	leftPad := strings.Repeat(" ", xPadSpaces-len(leftBorderSeq))

	for i := 0; i < asciiBorderHeight; i++ {
		if i == 0 || i == asciiBorderHeight-1 {
			if i == 0 {
				borderedASCIILines = append(borderedASCIILines, strings.Repeat(topBorderSeq, asciiBorderWidth/len(topBorderSeq)))
			} else if i == asciiBorderHeight-1 {
				borderedASCIILines = append(borderedASCIILines, strings.Repeat(bottomBorderSeq, asciiBorderWidth/len(bottomBorderSeq)))
			}
		} else {
			if spaceYCount < yPadSpaces || asciiDone {
				spaceYCount++
				toAddSpaces := strings.Repeat(" ", asciiBorderWidth-(len(leftBorderSeq)+len(rightBorderSeq)))
				borderedASCIILines = append(borderedASCIILines, leftBorderSeq+toAddSpaces+rightBorderSeq)
				continue
			}
			if asciiPos < len(lines) {
				asciiLineArt := lines[asciiPos]
				artLine := leftBorderSeq + leftPad + asciiLineArt
				toAddSpaces := strings.Repeat(" ", asciiBorderWidth-(len(artLine)+len(rightBorderSeq)))
				borderedASCIILines = append(borderedASCIILines, artLine+toAddSpaces+rightBorderSeq)
				asciiPos++
			} else {
				toAddSpaces := strings.Repeat(" ", asciiBorderWidth-(len(leftBorderSeq)+len(rightBorderSeq)))
				borderedASCIILines = append(borderedASCIILines, leftBorderSeq+toAddSpaces+rightBorderSeq)
				asciiDone = true
			}
		}
	}

	return strings.Join(borderedASCIILines, "\n")
}

const (
	asciiRectStampShape = ` ____
| %s |
|____|`
	asciiRectClassicStampShape = `######
# %s #
######`
	asciiCircleStampShape = `  __
 /  \
| %s |
 \__/`
	asciiCircleClassicStampShape = ` ####
# %s #
 ####`
)

func createASCIIStamp(shape enum.StampShapeEnum, country string) string {
	switch shape {
	case enum.StampShapeRect:
		return fmt.Sprintf(asciiRectStampShape, country)
	case enum.StampShapeCircle:
		return fmt.Sprintf(asciiCircleStampShape, country)
	case enum.StampShapeCircleClassic:
		return fmt.Sprintf(asciiCircleClassicStampShape, country)
	case enum.StampShapeRectClassic:
		fallthrough
	default:
		return fmt.Sprintf(asciiRectClassicStampShape, country)
	}
}

func getASCIIBorder(border enum.BorderEnum) asciiBorder {
	switch border {
	case enum.BorderCubes:
		fallthrough
	case enum.BorderPhoto:
		return photoBorder
	case enum.BorderLines:
		return linesBorder
	case enum.BorderStripes:
		return stripedBorder
	case enum.BorderStandard:
		fallthrough
	default:
		return defaultBorder
	}
}

func joinFrontASCIISides(left, right string, border asciiBorder) string {
	leftRows := strings.Split(left, "\n")
	rightRows := strings.Split(right, "\n")
	joined := []string{}

	for r := range leftRows {
		leftRow := leftRows[r]
		rightRow := rightRows[r]
		row := border.left + leftRow + "|" + rightRow + border.right
		joined = append(joined, row)
	}

	return strings.Join(joined, "\n")
}

func createFrontASCIILeftSide(from Person, message string, bCharSize int) string {
	height := asciiBorderHeight - 2
	sideLength := 20 - bCharSize
	senderASCII := createPersonASCII("From", from, sideLength)
	// stuff should fit in the side
	spaceLine := strings.Repeat(" ", sideLength)
	sideLines := strings.Split(senderASCII, "\n")
	sideLines = append(sideLines, spaceLine)

	maxMsgLines := height - 1 - len(sideLines)

	messageLines := createASCIIContentLines(message, sideLength-2, maxMsgLines)
	sideLines = append(sideLines, messageLines...)
	sideLines = append(sideLines, spaceLine)

	return strings.Join(sideLines, "\n")
}

func createFrontASCIIRightSide(to Person, shape enum.StampShapeEnum, country string, bCharSize int) string {
	height := asciiBorderHeight - 2
	sideLength := 19 - bCharSize
	spaceLine := strings.Repeat(" ", sideLength)
	sideLines := []string{spaceLine}

	// stamp ascii is one line from top
	stampASCII := createASCIIStamp(shape, country)
	leftPadForStamps := strings.Repeat(" ", sideLength/2+1)
	stampLines := strings.Split(stampASCII, "\n")
	for _, stampLine := range stampLines {
		paddedStampLine := leftPadForStamps + stampLine
		paddedStampLine += strings.Repeat(" ", sideLength-len(paddedStampLine))
		sideLines = append(sideLines, paddedStampLine)
	}

	// use height to add number of necessary space lines between end of stamp
	recipASCII := createPersonASCII("To", to, sideLength)
	recipHeight := strings.Count(recipASCII, "\n") + 1
	spaceBetweenRecipAndStamp := height - len(sideLines) - 1 - recipHeight
	for l := 0; l < spaceBetweenRecipAndStamp; l++ {
		sideLines = append(sideLines, spaceLine)
	}

	recipLines := strings.Split(recipASCII, "\n")
	sideLines = append(sideLines, recipLines...)

	sideLines = append(sideLines, spaceLine)
	return strings.Join(sideLines, "\n")
}

func createASCIIContentLines(info string, maxLineLen, maxLines int) []string {
	infoLen := len(info)
	currentLine := ""
	lines := []string{}
	for c, char := range info {
		currentLine += string(char)
		if len(currentLine) == maxLineLen-3 && len(lines) == maxLines-1 {
			currentLine += "..."
			lines = append(lines, currentLine)
			break
		}
		if len(currentLine) == maxLineLen {
			lines = append(lines, currentLine)
			currentLine = ""
		}
		if c == infoLen-1 {
			addedSpaces := strings.Repeat(" ", maxLineLen-len(currentLine))
			currentLine += addedSpaces
			lines = append(lines, currentLine)
		}
	}
	paddedLines := []string{}
	for _, line := range lines {
		paddedLines = append(paddedLines, " "+line+" ")
	}
	return paddedLines
}

func createPersonASCII(prefix string, person Person, sideLength int) string {
	topLine := fmt.Sprintf("%s:%s", prefix, strings.Repeat(" ", sideLength-(len(prefix)+1)))
	acceptableLength := sideLength - 2

	personLines := []string{topLine}
	if person.Name != "" {
		nameLines := createASCIIContentLines(person.Name, acceptableLength, 2)
		personLines = append(personLines, nameLines...)
	}
	if person.Email != "" {
		email := fmt.Sprintf("<%s>", person.Email)
		emailLines := createASCIIContentLines(email, acceptableLength, 2)
		personLines = append(personLines, emailLines...)
	}

	return strings.Join(personLines, "\n")
}

func breakupMessageForASCII(msg string) string {
	lines := []string{}
	currentLine := ""
	for i, c := range msg {
		currentLine += string(c)
		if len(currentLine) == asciiBorderWidth-1 {
			if c != ' ' && c != '-' && c != '\t' && c != '.' && c != ',' && c != ':' && c != ';' {
				currentLine += "-"
			}
			lines = append(lines, currentLine)
			currentLine = ""
		}
		if i == len(msg)-1 {
			lines = append(lines, currentLine)
		}
	}

	return strings.Join(lines, "\n")
}
