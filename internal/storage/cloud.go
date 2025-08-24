package storage

import "mime/multipart"

type CloudStorage interface {
	UploadFile(file *multipart.FileHeader, path string, oldID string) (string, error)
}
