package storage

import (
	"context"
	"log"
	"mime/multipart"
	"os"

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

func (c *CloudinaryStorage) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	defer func() {
		os.Remove(path)
	}()

	resp, err := c.Cloudinary.Upload.Upload(c.Context, path, uploader.UploadParams{
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(false),
	})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return resp.SecureURL, nil
}
