package utils

import (
	"image"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

func ScaleImage(file *os.File) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Создаем новое изображение с нужными размерами
	resizedImg := image.NewRGBA(image.Rect(startX, startY, endX, endY))

	// Рисуем исходное изображение в новом изображении
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Rect, img, img.Bounds(), draw.Over, nil)

	return resizedImg, nil
}
