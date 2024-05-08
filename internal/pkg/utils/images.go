package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
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

func WriteResizedFile(file *multipart.FileHeader, folderName string) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	img, _, err := image.Decode(uploadedFile)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return "", err
	}

	// Создаем новое изображение с нужными размерами
	resizedImg := image.NewRGBA(image.Rect(0, 0, 215, 295))

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
	filename := RandString(8) + extension
	fullpath := dirName + "/" + filename
	destination, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	err = jpeg.Encode(destination, resizedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return "", err
	}

	// if _, err := io.Copy(destination, uploadedFile); err != nil {
	// 	return "", err
	// }

	return fullpath, nil
}

func DecodeImageWithScaling(filename string) (string, error) {
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

	return encoded, nil
}

func DecodeImage(filename string) (string, error) {
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

	return encoded, nil
}
