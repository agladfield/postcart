package upload

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const uploadGCSErrFmtStr = "upload to gcs err: %w"

var gcpStorageClient *storage.Client

func uploadImageWithGoogleCloud(bytes []byte) (string, error) {
	if gcpStorageClient == nil {
		var clientErr error
		gcpStorageClient, clientErr = storage.NewClient(context.Background(), option.WithCredentialsFile(env.GCPCredsPath()))
		if clientErr != nil {
			return "", fmt.Errorf(uploadGCSErrFmtStr, clientErr)
		}
	}

	objectName := fmt.Sprintf("%s.jpg", uuid.New().String())
	bucket := gcpStorageClient.Bucket(env.GCPBucket())
	object := bucket.Object(objectName)
	background := object.NewWriter(context.Background())
	_, writeErr := background.Write(bytes)
	if writeErr != nil {
		return "", fmt.Errorf(uploadGCSErrFmtStr, writeErr)
	}

	writeCloseErr := background.Close()
	if writeCloseErr != nil {
		return "", fmt.Errorf(uploadGCSErrFmtStr, writeCloseErr)
	}

	signedURL, signingErr := bucket.SignedURL(objectName, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})
	if signingErr != nil {
		return "", fmt.Errorf(uploadGCSErrFmtStr, signingErr)
	}

	return signedURL, nil
}

// © Arthur Gladfield
