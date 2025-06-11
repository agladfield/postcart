package cards

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/agladfield/postcart/pkg/shared/tools/geo"
)

const cardsParseErrFmtStr = "cards parse err: %w"

func Parse(inbound *postmark.InboundData) (Params, error) {
	details, err := parseEmailBody(inbound.TextBody)
	if err != nil {
		if err == errMissingFromField {
			details.From = Person{
				Email: inbound.FromFull.Email,
				Name:  inbound.FromFull.Name,
			}
		} else {
			return Params{}, err
		}
	}

	details.ID = inbound.MessageID
	details.Subject = inbound.Subject

	if len(inbound.Attachments) > 0 && env.AllowAttachments() {
		details.Attachment = &inbound.Attachments[0]
	}

	return details, nil
}

// Person represents a name and email pair
type Person struct {
	Name  string
	Email string
}

// Params holds all extracted email parameters
type Params struct {
	ID         string
	To         Person
	From       Person
	Artwork    enum.ArtworkEnum
	Style      enum.StyleEnum
	Font       enum.FontEnum
	Border     enum.BorderEnum
	StampShape enum.StampShapeEnum
	Textured   enum.TexturedEnum
	Country    string
	Subject    string
	Message    string
	Attachment *postmark.EmailAttachment
}

func (p *Params) toJDBJobRecord() jdb.JobRecord {
	attachmentType := ""
	if p.Attachment != nil {
		attachmentType = p.Attachment.ContentType
	}
	return jdb.JobRecord{ID: p.ID,
		ToEmail:        p.To.Email,
		ToName:         p.To.Name,
		FromEmail:      p.From.Email,
		FromName:       p.From.Name,
		Artwork:        int8(p.Artwork),
		Style:          int8(p.Style),
		Font:           int8(p.Font),
		Border:         int8(p.Border),
		StampShape:     int8(p.StampShape),
		Textured:       int8(p.Textured),
		Country:        p.Country,
		Subject:        p.Subject,
		Message:        p.Message,
		AttachmentType: attachmentType,
	}
}

// Field key variations for flexible matching
var fieldKeys = map[string][]string{
	"to":       {"to:"},
	"from":     {"from:"},
	"artwork":  {"artwork:"},
	"artstyle": {"style:"},
	"border":   {"border:"},
	"font":     {"font:"},
	"shape":    {"shape:"},
	"country":  {"country:"},
	"textured": {"textured:"},
}

// Enum string mappings
var artworkMap = map[string]enum.ArtworkEnum{
	"attach":     enum.ArtworkAttachment,
	"attached":   enum.ArtworkAttachment,
	"attachment": enum.ArtworkAttachment,
	"mountains":  enum.ArtworkMountains,
	"lake":       enum.ArtworkLakeside,
	"lakeside":   enum.ArtworkLakeside,
	"city":       enum.ArtworkCity,
	"island":     enum.ArtworkIslands,
	"islands":    enum.ArtworkIslands,
}

var artStyleMap = map[string]enum.StyleEnum{
	"painting":      enum.StylePainting,
	"photo":         enum.StylePhotograph,
	"photograph":    enum.StylePhotograph,
	"vintage":       enum.StyleVintagePhoto,
	"vintage photo": enum.StyleVintagePhoto,
	"vintage-photo": enum.StyleVintagePhoto,
	"illustrated":   enum.StyleIllustrated,
	"illustration":  enum.StyleIllustrated,
	"cartoon":       enum.StyleIllustrated,
}

var borderMap = map[string]enum.BorderEnum{
	"none":    enum.BorderStandard,
	"classic": enum.BorderStandard,
	"default": enum.BorderStandard,
	"lines":   enum.BorderLines,
	"cubes":   enum.BorderCubes,
	"stripes": enum.BorderStripes,
	"art":     enum.BorderPhoto,
	"artwork": enum.BorderPhoto,
	"photo":   enum.BorderPhoto,
}

var shapeMap = map[string]enum.StampShapeEnum{
	"classic":        enum.StampShapeRectClassic,
	"default":        enum.StampShapeRectClassic,
	"rect":           enum.StampShapeRect,
	"square":         enum.StampShapeRect,
	"circle":         enum.StampShapeCircle,
	"circle-classic": enum.StampShapeCircleClassic,
}

var fontMap = map[string]enum.FontEnum{
	"typewriter": enum.FontTypewriter,
	"polite":     enum.FontPolite,
	"marker":     enum.FontMarker,
	"midcentury": enum.FontMidCentury,
	"vintage":    enum.FontMidCentury,
}

var texturedMap = map[string]enum.TexturedEnum{
	"yes":      enum.TexturedEnabled,
	"enabled":  enum.TexturedEnabled,
	"true":     enum.TexturedEnabled,
	"no":       enum.TexturedDisabled,
	"disabled": enum.TexturedDisabled,
	"false":    enum.TexturedDisabled,
}

// nearestString calculates the character distance between two strings
func nearestString(s, t string) int {
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}
	if s[0] == t[0] {
		return nearestString(s[1:], t[1:])
	}
	a := nearestString(s[1:], t) + 1
	b := nearestString(s, t[1:]) + 1
	c := nearestString(s[1:], t[1:]) + 1
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	return a
}

// parseEnum finds the closest enum value based on Levenshtein distance
func parseEnum[T comparable](input string, enumMap map[string]T, defaultValue T, threshold int) T {
	if input == "" {
		return defaultValue
	}
	input = strings.ToLower(input)
	minDistance := int(^uint(0) >> 1) // max int
	var closestEnum T
	for str, enum := range enumMap {
		dist := nearestString(input, strings.ToLower(str))
		if dist < minDistance {
			minDistance = dist
			closestEnum = enum
		}
	}
	if minDistance <= threshold {
		return closestEnum
	}
	return defaultValue
}

// parsePerson extracts name and email from a string
func parsePerson(input string) Person {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "<") && strings.Contains(input, ">") {
		parts := strings.SplitN(input, "<", 2)
		name := strings.TrimSpace(parts[0])
		email := strings.TrimSpace(parts[1])
		email = strings.TrimSuffix(email, ">")
		return Person{Name: name, Email: email}
	} else if strings.Contains(input, "@") {
		return Person{Email: input}
	}
	return Person{Name: input}
}

var errMissingFromField = errors.New("missing required field: From")

// parseEmailBody parses the email body into Params
func parseEmailBody(body string) (Params, error) {
	// start := time.Now()
	lines := strings.Split(body, "\n")
	var messageLines []string
	inMessage := false
	params := make(map[string]string)

	// Separate parameter lines from message lines
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if inMessage {
			messageLines = append(messageLines, trimmed)
			continue
		}
		lowerLine := strings.ToLower(trimmed)
		found := false
		for field, keys := range fieldKeys {
			for _, key := range keys {
				if strings.HasPrefix(lowerLine, key) {
					value := strings.TrimSpace(line[len(key):])
					params[field] = value
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			inMessage = true
			messageLines = append(messageLines, trimmed)
		}
	}

	// Populate the struct
	var result Params

	// Required field: To
	if val, ok := params["to"]; ok {
		result.To = parsePerson(val)
	} else {
		return Params{}, fmt.Errorf(cardsParseErrFmtStr, errors.New("missing required field: To"))
	}

	// Enum fields with fuzzy matching
	result.Artwork = parseEnum(params["artwork"], artworkMap, enum.ArtworkUnknown, 2)
	result.Border = parseEnum(params["border"], borderMap, enum.BorderUnknown, 2)
	result.StampShape = parseEnum(params["shape"], shapeMap, enum.StampShapeUnknown, 2)
	result.Font = parseEnum(params["font"], fontMap, enum.FontUnknown, 2)
	result.Style = parseEnum(params["artstyle"], artStyleMap, enum.StyleUnknown, 2)
	result.Textured = parseEnum(params["textured"], texturedMap, enum.TexturedUnknown, 2)

	// Country with normalization

	result.Country = geo.GetCountry(params["country"])

	// Message
	result.Message = strings.Join(messageLines, "\n")

	// Required field: From with special handling
	if val, ok := params["from"]; ok {
		trimmedVal := strings.TrimSpace(val)
		lowerVal := strings.ToLower(trimmedVal)
		if trimmedVal == "" || strings.Contains(lowerVal, "anon") || strings.Contains(lowerVal, "anonymous") {
			result.From = Person{Name: "Anonymous"}
		} else {
			result.From = parsePerson(val)
		}
	} else {
		return result, errMissingFromField
	}

	return result, nil
}

// © Arthur Gladfield
