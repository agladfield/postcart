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
	"github.com/davidbyttow/govips/v2/vips"
)

var postcardEntries = []cards.EmailParams{
	{
		ID: "grandma",
		To: cards.Person{
			Name:  "Grandma",
			Email: "grandma@email.madeup",
		}, From: cards.Person{
			Name:  "Arthur",
			Email: "arthur@email.madeup",
		},
		StampShape: enum.StampShapeRect,
		Border:     enum.BorderStandard,
		Artwork:    enum.ArtworkCity,
		Country:    "US",
		Subject:    "Cookie Recipe",
		Message:    "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
	},
	{
		ID: "balajis",
		To: cards.Person{
			Name:  "Surangel Samuel Whipps Jr.",
			Email: "pres@gov.pw",
		}, From: cards.Person{
			Name:  "Balajis",
			Email: "balajis@ns.edu",
		},
		StampShape: enum.StampShapeCircleClassic,
		Border:     enum.BorderCubes,
		Artwork:    enum.ArtworkIslands,
		Country:    "PW",
		Subject:    "Network States",
		Message:    "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
	},
	{
		ID: "grandma",
		To: cards.Person{
			Name:  "Grandma",
			Email: "grandma@email.madeup",
		}, From: cards.Person{
			Name:  "Arthur",
			Email: "arthur@email.madeup",
		},
		StampShape: enum.StampShapeRect,
		Border:     enum.BorderStandard,
		Artwork:    enum.ArtworkCity,
		Country:    "US",
		Subject:    "Cookie Recipe",
		Message:    "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
	},
	// {
	// 	ID: "grandma",
	// 	To: cards.Person{
	// 		Name:  "Grandma",
	// 		Email: "grandma@email.madeup",
	// 	}, From: cards.Person{
	// 		Name:  "Arthur",
	// 		Email: "arthur@email.madeup",
	// 	},
	// 	StampShape: enum.StampShapeRect,
	// 	Border:     enum.BorderStandard,
	// 	Artwork:    enum.ArtworkCity,
	// 	Country:    "US",
	// 	Subject:    "Cookie Recipe",
	// 	Message:    "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
	// },
	{
		ID: "attachment",
		To: cards.Person{
			Name:  "President Coffee",
			Email: "pres@bigcoffee.coffee",
		}, From: cards.Person{
			Name:  "Prime",
			Email: "primagen@terminal.coffee",
		},
		Border:  enum.BorderPhoto,
		Artwork: enum.ArtworkAttachment,
		Country: "BR",
		Subject: "Coffee Deal",
		Message: "It was good visiting you in the city the other day Grandma. If you could get me the recipe for the chocolate chip cookies that would be spectacular. I know I loved them and they're a real fan favorite.",
		Attachment: &postmark.EmailAttachment{
			Content:     "./pkg/cards/res/artwork/attachment.jpeg",
			ContentType: "image/jpeg",
		},
	},
}

func demo() error {
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
	for _, postcard := range postcardEntries {
		// should do the generation postjob
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
