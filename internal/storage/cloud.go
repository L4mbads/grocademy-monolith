package storage

import "mime/multipart"

type CloudStorage interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
}
