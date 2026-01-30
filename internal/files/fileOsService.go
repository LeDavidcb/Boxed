package files

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func saveFile(fpath string, file *multipart.FileHeader) error {
	// Try to create the directory before the files is created.

	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(fpath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
func createAndSaveThumbnail(fpath string) {
	log.Println("#### TODO FUNC WAS CALLED, createAndSaveThumbnail().", fpath)
}
func deleteFile(fpath string) {
	err := os.RemoveAll(fpath)
	if err != nil {
		log.Println("Fatal error: coudln't remove the following file:", fpath)
		// maybe re-add the entry that was deleted in the database.
	}
}
