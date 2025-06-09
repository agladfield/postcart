package upload

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	imgbb "github.com/JohnNON/ImgBB"
)

// 46a75c1b4191cc379e8018a9bd0b6e4a

const imgBBApiKey = "46a75c1b4191cc379e8018a9bd0b6e4a"

var (
	imgBBHTTPClient *http.Client
	imgBBClient     *imgbb.Client
)

// expire in 180 days
const expiration = 15_552_000

func uploadImageWithImageBB(imgBytes []byte) (string, error) {
	if imgBBHTTPClient == nil {
		imgBBHTTPClient = &http.Client{
			Timeout: 20 * time.Second,
		}
	}
	if imgBBClient == nil {
		imgBBClient = imgbb.NewClient(imgBBHTTPClient, imgBBApiKey)
	}

	newImage, newImgErr := imgbb.NewImageFromFile(hashSumForImageBB(imgBytes), expiration, imgBytes)
	if newImgErr != nil {
		return "", newImgErr
	}

	resp, uploadErr := imgBBClient.Upload(context.Background(), newImage)
	if uploadErr != nil {
		return "", uploadErr
	}

	return resp.Data.URL, nil
}

func hashSumForImageBB(b []byte) string {
	sum := md5.Sum(b)

	return hex.EncodeToString(sum[:])
}
