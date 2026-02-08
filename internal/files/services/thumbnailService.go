package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
)

// createAndSaveThumbnail generates a thumbnail for the given file and saves it to disk.
//
// Parameters:
//   - fpath: The file path where the thumbnail is expected to be saved.
func CreateAndSaveThumbnail(inPath, outPath, mime, originalName string, thumbnailEntry uuid.UUID, repository *repositories.ThumbnailRepository) error {
	if inPath == "" {
		return fmt.Errorf("Input path is empty")
	}

	// Check if input file exists
	if _, err := os.Stat(inPath); os.IsNotExist(err) {
		return fmt.Errorf("Input file does not exist: %s", inPath)
	}

	if outPath == "" {
		return fmt.Errorf("Output path is empty")
	}

	// Ensure parent directories for the output file exist
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("Failed to create parent directories for output file: %w", err)
	}

	mimeParts := strings.Split(mime, "/")
	if len(mimeParts) != 2 {
		return fmt.Errorf("Invalid MIME type: %s", mime)
	}
	mimeType := mimeParts[0]
	switch mimeType {
	case "video":
		log.Println("LLEGOOwnflka")
		err := GenerateVideoThumbnail(context.Background(), inPath, outPath)
		if err != nil {
			return err
		}
		// Fill thumbnail row
		thumbnail, err := repository.GetByID(thumbnailEntry)
		if err != nil {
			return err
		}
		thumbnail.StoragePath = outPath
		thumbnail.OriginalName = originalName
		return repository.UpdateByID(thumbnail)

	case "image":
		err := GenerateImageThumbnail(context.Background(), inPath, outPath)
		if err != nil {
			return err
		}
		// Fill thumbnail row
		thumbnail, err := repository.GetByID(thumbnailEntry)
		if err != nil {
			return err
		}
		thumbnail.StoragePath = outPath
		thumbnail.OriginalName = originalName
		return repository.UpdateByID(thumbnail)
	default:
		log.Println("Unsupported MIME type")
		return fmt.Errorf("Unsupported MIME type")
	}
}
func GetVideoDuration(c context.Context, iPath string) (float64, error) {
	cmd := exec.CommandContext(
		c,
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		iPath,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// GenerateVideoThumbnail will generate a thumbnail at output path, for a given file (input)
// Depends on ffmpeg to work.
func GenerateVideoThumbnail(c context.Context, input, output string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	duration, err := GetVideoDuration(ctx, input)
	if err != nil {
		return err
	}

	half := duration / 2

	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-ss", fmt.Sprintf("%.2f", half),
		"-i", input,
		"-vframes", "1",
		"-vf", "scale=320:-1",
		output,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, out)
	}

	return nil
}

func GenerateImageThumbnail(c context.Context, input, output string) error {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		c,
		"ffmpeg",
		"-y",
		"-i", input,
		"-vf", "scale=320:-1", // same logic as video thumbnails
		output,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, out)
	}

	return nil
}
