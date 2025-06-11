// Package upload wraps different object upload providers for uploading
// images to a bucket and returning consumable urls
package upload

import "github.com/agladfield/postcart/pkg/shared/env"

func UploadImage(bytes []byte) (string, error) {
	if env.GCPCredsPath() != "" && env.GCPBucket() != "" {
		return uploadImageWithGoogleCloud(bytes)
	} else {
		return uploadImageToTmpFiles(bytes, "output.jpg")
	}
}

// © Arthur Gladfield
