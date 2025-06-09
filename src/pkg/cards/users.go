package cards

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/agladfield/postcart/pkg/pdb"
	"github.com/google/uuid"
)

// func GetSenderIDByEmail(email string) (string, error) {
// 	sender, senderErr := GetSenderByEmail(email)
// 	if senderErr != nil {
// 		return "", senderErr
// 	}
// 	return sender.ID, nil
// }

func GetSenderByEmail(email string) (*pdb.Sender, error) {
	db := pdb.Obtain()
	tx, txErr := db.DB.Begin()
	if txErr != nil {
		// should signal retry
		return nil, txErr
	}
	defer tx.Rollback()
	qtx := db.WithTx(tx)
	// query for user
	senderRes, senderErr := qtx.GetSenderByEmail(context.Background(), email)
	if senderErr != nil {
		if senderErr != sql.ErrNoRows {
			return nil, senderErr
		}
		if senderRes != nil && senderErr != sql.ErrNoRows {
			commitErr := tx.Commit()
			if commitErr != nil {
				return nil, commitErr
			}
			return senderRes, nil
		}
	}
	fmt.Println("sender was nil, we are making a new one")
	newSender := createSenderData(email)
	createErr := qtx.CreateSender(context.Background(), *newSender)
	if createErr != nil {
		return nil, createErr
	}

	// if none found for email, create new user
	// use transactions
	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, commitErr
	}

	return &pdb.Sender{
		ID:        newSender.ID,
		Created:   newSender.Created,
		LastSent:  newSender.LastSent,
		Email:     email,
		Sent:      newSender.Sent,
		Fails:     newSender.Fails,
		Delivered: newSender.Delivered,
		Blocked:   newSender.Blocked,
	}, nil
}

func createSenderData(email string) *pdb.CreateSenderParams {
	now := time.Now().Unix()
	newSender := pdb.CreateSenderParams{
		ID:        uuid.New().String(),
		Created:   now,
		LastSent:  now,
		Email:     email,
		Sent:      1,
		Fails:     0,
		Delivered: 0,
		Blocked:   0,
	}
	return &newSender
}

func BlockSender(userID string) error {
	//
	return nil
}

func BlockSenderEmail(email string) error {
	// add to trigger rules
	return nil
}

func BlockRecipientEmail(email string) error {
	// add to forbidden
	return nil
}
