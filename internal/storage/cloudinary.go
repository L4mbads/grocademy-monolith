package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryStorage struct {
	Cloudinary *cloudinary.Cloudinary
	Context    context.Context
}

func NewCloudinaryStorage() (*CloudinaryStorage, error) {
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()

	return &CloudinaryStorage{
		Cloudinary: cld,
		Context:    ctx,
	}, nil
}

func (c *CloudinaryStorage) UploadFile(file *multipart.FileHeader, path string, oldURL string) (string, error) {

	defer func() {
		println("removing from local")
		err := os.Remove(path)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	if oldURL != "" {
		oldID, err := extractPublicID(oldURL)
		println("Parsed publicID: " + oldID)
		if err != nil {
			fmt.Println("Error parsing publicID: ", err)

		} else {
			dresp, derr := c.Cloudinary.Upload.Destroy(c.Context, uploader.DestroyParams{
				PublicID: oldID,
			})
			if derr != nil {
				fmt.Println("Error deleting from cloud: ", derr)
			} else {
				if dresp.Result == "ok" {
					println("Cloud deletion success")
				} else {
					println("Cloud deletion failed")
				}
			}
		}
	}

	resp, err := c.Cloudinary.Upload.Upload(c.Context, path, uploader.UploadParams{
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(true),
	})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return resp.SecureURL, nil
}

func extractPublicID(cloudinaryURL string) (string, error) {
	parsedURL, err := url.Parse(cloudinaryURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %w", err)
	}

	pathSegments := strings.Split(parsedURL.Path, "/")

	// Find the "upload" segment and then look for the public ID
	for i, segment := range pathSegments {
		if segment == "upload" && i+1 < len(pathSegments) {
			// Check if there's a version number (e.g., v123456789)
			if strings.HasPrefix(pathSegments[i+1], "v") && len(pathSegments[i+1]) > 1 {
				// Public ID is after the version number
				if i+2 < len(pathSegments) {
					publicID := strings.Join(pathSegments[i+2:], "/")
					// Remove file extension if present
					if dotIndex := strings.LastIndex(publicID, "."); dotIndex != -1 {
						publicID = publicID[:dotIndex]
					}
					return publicID, nil
				}
			} else {
				// Public ID is directly after "upload" (no version)
				publicID := strings.Join(pathSegments[i+1:], "/")
				// Remove file extension if present
				if dotIndex := strings.LastIndex(publicID, "."); dotIndex != -1 {
					publicID = publicID[:dotIndex]
				}
				return publicID, nil
			}
		}
	}

	return "", fmt.Errorf("public ID not found in URL: %s", cloudinaryURL)
}
