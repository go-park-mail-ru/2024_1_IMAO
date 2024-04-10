package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func WriteFile(file *multipart.FileHeader, folderName string) error {
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	currentTime := time.Now()

	dirName := fmt.Sprintf("./uploads/%s/%d-%02d-%02d", folderName,
		currentTime.Year(), currentTime.Month(), currentTime.Day())

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	extension := filepath.Ext(file.Filename)
	filename := RandString(8) + extension
	destination, err := os.Create(dirName + "/" + filename)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, uploadedFile); err != nil {
		return err
	}

	return nil
}
