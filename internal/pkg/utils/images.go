package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	staticDirectory = "./uploads"
)

func WriteFile(file *multipart.FileHeader, folderName string) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	currentTime := time.Now()

	dirName := fmt.Sprintf("%s/%s/%d-%02d-%02d", staticDirectory, folderName,
		currentTime.Year(), currentTime.Month(), currentTime.Day())

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	filename := RandString(8) + extension
	fullpath := dirName + "/" + filename
	destination, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, uploadedFile); err != nil {
		return "", err
	}

	return fullpath, nil
}

func DecodeImage(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", nil
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", nil
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}
