package files

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

// saveFile saves to the file system.
// Parameters:
//   - fpath: The path where the file's going to be saved.
//   - file: a pointer to the FileHeader, you could find this in a http request.
//
// Returns:
//   - An error if any os-related operation went wrong.
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

// createAndSaveThumbnail generates a thumbnail for the given file and saves it to disk.
// This function is currently a placeholder for future implementation.
//
// Parameters:
//   - fpath: The file path where the thumbnail is expected to be saved.
func createAndSaveThumbnail(fpath string) {
	log.Println("#### TODO FUNC WAS CALLED, createAndSaveThumbnail().", fpath)
}

// deleteFile removes a file or directory from the file system.
//
// Parameters:
//   - fpath: The path to the file or directory to be deleted.
//
// Logs:
//   - Logs an error message if the file cannot be deleted successfully.
func deleteFile(fpath string) {
	err := os.RemoveAll(fpath)
	if err != nil {
		log.Println("Fatal error: coudln't remove the following file:", fpath)
		// maybe re-add the entry that was deleted in the database.
	}
}
