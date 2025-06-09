package cards

import (
	"errors"
	"strings"

	"github.com/agladfield/postcart/pkg/pdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/geo"
)

// Person represents a name and email pair
type Person struct {
	Name  string
	Email string
}

// EmailParams holds all extracted email parameters
type EmailParams struct {
	ID         string
	To         Person
	From       Person
	Artwork    enum.ArtworkEnum
	ArtStyle   enum.StyleEnum
	Font       enum.FontEnum
	Border     enum.BorderEnum
	StampShape enum.StampShapeEnum
	Country    string
	Subject    string
	Message    string
	Attachment *postmark.EmailAttachment
}

func (ep *EmailParams) ToQueueRequest(userID string) *pdb.SetQueuedRequestParams {
	var attachment []byte
	if ep.Attachment != nil {
		attachment = []byte(ep.Attachment.Content)
	}

	queued := pdb.SetQueuedRequestParams{
		ID:          ep.ID,
		User:        userID,
		ToEmail:     ep.To.Email,
		ToName:      ep.To.Name,
		FromName:    ep.From.Name,
		FromEmail:   ep.From.Email,
		ArtworkEnum: int64(ep.Artwork),
		BorderEnum:  int64(ep.Border),
		FontEnum:    int64(ep.Font),
		ShapeEnum:   int64(ep.StampShape),
		StyleEnum:   int64(ep.ArtStyle),
		Country:     ep.Country,
		Message:     ep.Message,
		Attachment:  attachment,
	}

	return &queued
}

// func (ep *EmailParams) ToCompletedRequest() *pdb.QueuedRequest {
// 	var attachment []byte
// 	if ep.Attachment != nil {
// 		attachment = []byte(ep.Attachment.Content)
// 	}

// 	queued := pdb.QueuedRequest{
// 		ID:          ep.ID,
// 		User:        "user",
// 		ToEmail:     ep.To.Email,
// 		ToName:      ep.To.Name,
// 		FromName:    ep.From.Name,
// 		FromEmail:   ep.From.Email,
// 		ArtworkEnum: int64(ep.Artwork),
// 		BorderEnum:  int64(ep.Border),
// 		FontEnum:    int64(ep.Font),
// 		ShapeEnum:   int64(ep.StampShape),
// 		StyleEnum:   int64(ep.ArtStyle),
// 		Country:     ep.Country,
// 		Message:     ep.Message,
// 		HadAttachment:  false,
// 	}

// 	return &queued
// }

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
	"rect":    enum.BorderStandard,
	"line":    enum.BorderLines,
	"lines":   enum.BorderLines,
	"cube":    enum.BorderCubes,
	"cubes":   enum.BorderCubes,
	"stripe":  enum.BorderStripes,
	"stripes": enum.BorderStripes,
	"striped": enum.BorderStripes,
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

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s, t string) int {
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}
	if s[0] == t[0] {
		return levenshteinDistance(s[1:], t[1:])
	}
	a := levenshteinDistance(s[1:], t) + 1
	b := levenshteinDistance(s, t[1:]) + 1
	c := levenshteinDistance(s[1:], t[1:]) + 1
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
		dist := levenshteinDistance(input, strings.ToLower(str))
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

// parseEmailBody parses the email body into EmailParams
func parseEmailBody(body string) (EmailParams, error) {
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
	var result EmailParams

	// Required field: To
	if val, ok := params["to"]; ok {
		result.To = parsePerson(val)
	} else {
		return EmailParams{}, errors.New("missing required field: To")
	}

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
		return EmailParams{}, errors.New("missing required field: From")
	}

	// Enum fields with fuzzy matching
	result.Artwork = parseEnum(params["artwork"], artworkMap, enum.ArtworkUnknown, 2)
	result.Border = parseEnum(params["border"], borderMap, enum.BorderUnknown, 2)
	result.StampShape = parseEnum(params["shape"], shapeMap, enum.StampShapeUnknown, 2)
	result.Font = parseEnum(params["font"], fontMap, enum.FontUnknown, 2)
	// fmt.Println("took:", time.Since(start))
	// artStyleStart := time.Now()
	result.ArtStyle = parseEnum(params["artstyle"], artStyleMap, enum.StyleUnknown, 2)
	// fmt.Println("art style took:", time.Since(artStyleStart))

	// Country with normalization

	result.Country = geo.GetCountry(params["country"])

	// Message
	result.Message = strings.Join(messageLines, "\n")

	return result, nil
}
