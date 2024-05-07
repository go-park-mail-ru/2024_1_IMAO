package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func WriteFile(file *multipart.FileHeader, folderName string) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	currentTime := time.Now()

	dirName := fmt.Sprintf("./uploads/%s/%d-%02d-%02d", folderName,
		currentTime.Year(), currentTime.Month(), currentTime.Day())

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	//extension := ".jpg"

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
		fmt.Printf("Ошибка при открытии изображения: %v", err)
		return "", nil
	}

	var image image.Image

	image, err = ScaleImage(file)

	if err != nil {
		fmt.Printf("Ошибка при масштабировании изображения: %v", err)
		return "", nil
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, image, &jpeg.Options{Quality: 90})
	if err != nil {
		fmt.Printf("Ошибка при кодировании изображения: %v", err)
	}

	content := buf.Bytes()

	encoded := base64.StdEncoding.EncodeToString(content)

	fmt.Println(encoded)

	return encoded, nil
}

func DecodeImageWithoutScaling(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Ошибка при открытии изображения: %v", err)
		return "", nil
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", nil
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	fmt.Println(encoded)

	return encoded, nil
}
