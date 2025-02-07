package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func ImageToPDF() {
	imgPath := "image.jpg" // Путь к изображению
	pdfPath := "image.pdf" // Путь для сохранения PDF

	img, err := loadImage(imgPath)
	if err != nil {
		panic(err)
	}
	pdf := gofpdf.New("P", "mm", "A4", "") // новый файл
	pdf.AddPage()

	pageWidth, pageHeight := 210.0, 297.0 // размеры A4 в мм
	width, height := float64(img.Bounds().Size().X)*0.264583, float64(img.Bounds().Size().Y)*0.264583

	// Масштабирование изображения
	aspectRatio := width / height
	if aspectRatio > pageWidth/pageHeight {
		height = pageWidth / aspectRatio
		width = pageWidth
	} else {
		width = pageHeight * aspectRatio
		height = pageHeight
	}

	// Добавляем изображение в PDF
	pdf.Image(imgPath, (pageWidth-width)/2, (pageHeight-height)/2, width, height, false, "", 0, "")

	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		panic(err)
	}
}

func loadImage(path string) (image.Image, error) {
	// Открытие файла с изображением
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Декодирование изображения из файла
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func TXTtoPDF() {
	text, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.MoveTo(0, 0)
	pdf.SetFont("Arial", "", 25)
	width, _ := pdf.GetPageSize()
	pdf.MultiCell(width, 10, string(text), "0", "C", false)
	err = pdf.OutputFileAndClose("hello.pdf")
	if err == nil {
		fmt.Println("PDF generated successfully")
	}
}

func main() {
	var input string
	fmt.Println("если вы хотите конвертировать txt/docx в pdf введите 1")
	fmt.Println("если вы хотите конвертировать jpg в pdf введите 2")
	_, err := fmt.Scanln(&input) // Считываем ввод пользователя
	if err != nil {
		fmt.Println("Ошибка при вводе:", err)
		return
	}
	if input == "1" {
		TXTtoPDF()
	}
	if input == "2" {
		ImageToPDF()
	}

}
