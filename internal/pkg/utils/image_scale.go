package utils

import (
	"image"
	"os"

	"golang.org/x/image/draw"
)

func ScaleImage(file *os.File) (image.Image, error) {

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// // Вычисляем координаты центра изображения
	// centerX := img.Bounds().Dx() / 2
	// centerY := img.Bounds().Dy() / 2

	// minSide := min(img.Bounds().Dx(), img.Bounds().Dy())
	// halfMinSide := minSide / 2

	// // Создаем новый изображение для центрального квадрата
	// centerSquare := image.NewRGBA(image.Rect(0, 0, minSide, minSide))

	// // Копируем центр изображения в новый квадрат
	// for i := 0; i < minSide; i++ {
	// 	for j := 0; j < minSide; j++ {
	// 		x := centerX - halfMinSide + i
	// 		y := centerY - halfMinSide + j
	// 		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
	// 			centerSquare.Set(i, j, img.At(x, y))
	// 		} else {
	// 			centerSquare.Set(i, j, color.Transparent)
	// 		}
	// 	}
	// }

	// Создаем новое изображение с нужными размерами
	resizedImg := image.NewRGBA(image.Rect(0, 0, 215, 295))

	// Рисуем исходное изображение в новом изображении
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Rect, img, img.Bounds(), draw.Over, nil)

	return resizedImg, nil
}
