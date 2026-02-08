package register

import (
	"os"
	"path"

	boxed "github.com/David/Boxed"
	"github.com/google/uuid"
)

// Creates a directory based on the uuid provided. The user must save their files in this directory.
func CreateDirectory(uuid uuid.UUID) (string, error) {
	var folderPath string = path.Join(boxed.GetInstance().FolderPath, uuid.String())
	return folderPath, os.MkdirAll(folderPath, os.ModePerm)
}

// Deletes the directory created by CreateDirectory function.
func DeleteDirectory(uuid uuid.UUID) error {
	var folderPath string = path.Join(boxed.GetInstance().FolderPath, uuid.String())
	return os.RemoveAll(folderPath)
}
