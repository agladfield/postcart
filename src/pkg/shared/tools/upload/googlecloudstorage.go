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

var gcpStorageClient *storage.Client

func uploadImageWithGoogleCloud(bytes []byte) (string, error) {
	if gcpStorageClient == nil {
		var clientErr error
		gcpStorageClient, clientErr = storage.NewClient(context.Background(), option.WithCredentialsFile("/Users/aglad/Downloads/postcart-ce3696517869.json"))
		if clientErr != nil {
			return "", clientErr
		}
	}

	objectName := fmt.Sprintf("%s.jpg", uuid.New().String())
	bucket := gcpStorageClient.Bucket(env.GCPBucket())
	object := bucket.Object(objectName)
	background := object.NewWriter(context.Background())
	_, writeErr := background.Write(bytes)
	if writeErr != nil {
		return "", writeErr
	}

	writeCloseErr := background.Close()
	if writeCloseErr != nil {
		return "", writeCloseErr
	}

	signedURL, signingErr := bucket.SignedURL(objectName, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})
	if signingErr != nil {
		return "", signingErr
	}

	//

	return signedURL, nil
}
