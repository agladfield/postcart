package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/tools/random"
	"github.com/davidbyttow/govips/v2/vips"
)

var postcardEntries = []cards.Params{
	{
		ID: "cookie-recipe",
		To: cards.Person{
			Name:  "Grandma",
			Email: "grandma@email.madeup",
		}, From: cards.Person{
			Name:  "Arthur",
			Email: "arthur@email.madeup",
		},
		StampShape: enum.StampShapeRectClassic,
		Border:     enum.BorderStandard,
		Artwork:    enum.ArtworkCity,
		Textured:   enum.TexturedEnabled,
		Font:       enum.FontMarker,
		Country:    "PC",
		Message:    "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
	},
	{
		ID: "network-state",
		To: cards.Person{
			Name:  "Surangel Samuel Whipps Jr.",
			Email: "pres@gov.pw",
		}, From: cards.Person{
			Name:  "Balajis",
			Email: "balajis@ns.edu",
		},
		StampShape: enum.StampShapeCircleClassic,
		Border:     enum.BorderLines,
		Artwork:    enum.ArtworkIslands,
		Font:       enum.FontPolite,
		Textured:   enum.TexturedDisabled,
		Country:    "PW",
		Message:    "Thanks for doing the pod and taking Palau's future seriously.\n You've got some exciting programs in the\n works and I can't wait to see what else you come up with. \nIt's only fitting as one of the world's most unique currency innovators. \nCongrats on the islands!",
	},
	{
		ID: "vineyard-woes",
		To: cards.Person{
			Name:  "Data",
			Email: "data@androids.star.fleet",
		}, From: cards.Person{
			Name:  "Jean-Luc Picard",
			Email: "jeanlucpicard@star.fleet",
		},
		StampShape: enum.StampShapeRect,
		Border:     enum.BorderStripes,
		Artwork:    enum.ArtworkLakeside,
		Font:       enum.FontMidCentury,
		Textured:   enum.TexturedEnabled,
		Country:    "FR",
		Message: `My vineyards are suffering tremendously from this heat Data. I would be most appreciative if you could come up with a solution for my particular species of grapes that does not end up watering down the wine.

Thanks, Jean-Luc`,
	},
	{
		ID: "west-virginia",
		To: cards.Person{
			Name:  "Mountains",
			Email: "mountains@nature.gov",
		}, From: cards.Person{
			Name:  "John D.",
			Email: "johnd@music.guy",
		},
		StampShape: enum.StampShapeRectClassic,
		Border:     enum.BorderCubes,
		Artwork:    enum.ArtworkMountains,
		Font:       enum.FontTypewriter,
		Country:    "US",
		Message: `"Country roads, take me home
To the place I belong
West Virginia, mountain mama
Take me home, country roads"`,
	},
	{
		ID: "attachment",
		To: cards.Person{
			Name:  "President Coffee",
			Email: "pres@bigcoffee.coffee",
		}, From: cards.Person{
			Name:  "Prime",
			Email: "primagen@terminal.coffee",
		},
		Border:     enum.BorderPhoto,
		Artwork:    enum.ArtworkAttachment,
		Country:    "BR",
		StampShape: enum.StampShapeCircle,
		Font:       enum.FontMarker,
		Message: `Dear President Coffee,
I know your margins are hurting since you didn't have the bright idea to sell coffee over ssh. My team and I are willing to offer you a deal to save your company. Have your people reach out to my people and maybe you don't have to go bankrupt.
Best,
The Coffeeagen`,
		Attachment: &postmark.EmailAttachment{
			Content:     "./pkg/cards/res/artwork/attachment.jpeg",
			ContentType: "image/jpeg",
		},
	},
}

func demo() error {
	random.SetSeed(1234)
	ctx := context.Background()
	dirErr := os.MkdirAll("./demo", 0700)
	if dirErr != nil {
		return dirErr
	}
	var wg sync.WaitGroup
	prepErr := cards.Prepare(ctx, &wg)
	if prepErr != nil {
		return prepErr
	}
	defer cards.Close()
	for _, postcard := range postcardEntries {
		if postcard.Attachment != nil {
			attachmentBytes, loadErr := os.ReadFile(postcard.Attachment.Content)
			if loadErr != nil {
				return loadErr
			}
			postcard.Attachment.Content = base64.StdEncoding.EncodeToString(attachmentBytes)
		}

		unified, err := cards.Create(ctx, &postcard)
		if err != nil {
			return err
		}

		_, byteErr := unified.UnifiedImage.ToBytes()
		if byteErr != nil {
			return byteErr
		}

		bytes, _, exportErr := unified.UnifiedImage.ExportJpeg(&vips.JpegExportParams{
			Quality: 90,
		})
		if exportErr != nil {
			return exportErr
		}

		writeImageErr := os.WriteFile(fmt.Sprintf("./demo/%s-image.jpg", postcard.ID), bytes, 0600)
		if writeImageErr != nil {
			return writeImageErr
		}

		writeASCIIErr := os.WriteFile(fmt.Sprintf("./demo/%s-ascii.txt", postcard.ID), []byte(unified.UnifiedText), 0600)
		if writeASCIIErr != nil {
			return writeASCIIErr
		}

		fmt.Println()
		fmt.Println(unified.UnifiedText)
	}
	return nil
}

func main() {
	demoErr := demo()
	if demoErr != nil {
		log.Fatalln(demoErr)
	}
	os.Exit(0)
}

// © Arthur Gladfield
