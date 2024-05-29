package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
)

const (
	staticDirectory = "./uploads"
	quality         = 90
	filenameLen     = 8
	startX          = 0
	startY          = 0
	endX            = 215
	endY            = 295
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
	filename := RandString(filenameLen) + extension
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

func WriteResizedFile(file *multipart.FileHeader, folderName string) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	img, _, err := image.Decode(uploadedFile)
	if err != nil {
		log.Println("Ошибка при декодировании изображения:", err)

		return "", err
	}

	// Создаем новое изображение с нужными размерами
	resizedImg := image.NewRGBA(image.Rect(startX, startY, endX, endY))

	// Рисуем исходное изображение в новом изображении
	draw.CatmullRom.Scale(resizedImg, resizedImg.Rect, img, img.Bounds(), draw.Over, nil)

	currentTime := time.Now()

	dirName := fmt.Sprintf("%s/%s/%d-%02d-%02d", staticDirectory, folderName,
		currentTime.Year(), currentTime.Month(), currentTime.Day())

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	filename := RandString(filenameLen) + extension
	fullpath := dirName + "/" + filename

	destination, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}

	defer destination.Close()

	err = jpeg.Encode(destination, resizedImg, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Println("Ошибка при сохранении изображения:", err)

		return "", err
	}

	return fullpath, nil
}

func DecodeImageWithScaling(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Ошибка при открытии изображения: %v", err)

		return "", nil
	}

	var img image.Image

	img, err = ScaleImage(file)
	if err != nil {
		log.Printf("Ошибка при масштабировании изображения: %v", err)

		return "", nil
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Printf("Ошибка при кодировании изображения: %v", err)
	}

	content := buf.Bytes()
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}

func DecodeImage(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Ошибка при открытии изображения: %v", err)

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
