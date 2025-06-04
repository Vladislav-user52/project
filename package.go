package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/spreadsheet"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Использование: converter <файл> <исходное_расширение> <целевое_расширение>")
		fmt.Println("Пример: converter image.jpg jpg pdf")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	sourceExt := strings.ToLower(os.Args[2])
	targetExt := strings.ToLower(os.Args[3])

	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "." + targetExt

	switch sourceExt + "->" + targetExt {
	case "jpg->pdf", "jpeg->pdf", "png->pdf":
		err := ImageToPDF(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "txt->pdf", "docx->pdf":
		err := TXTtoPDF(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "pptx->pdf":
		err := PPTXtoPDF(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "xlsx->pdf":
		err := XLSXtoPDF(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "pdf->jpg", "pdf->jpeg", "pdf->png":
		err := PDFtoImage(inputFile, outputFile, targetExt)
		if err != nil {
			log.Fatal(err)
		}
	case "pdf->txt":
		err := PDFtoTXT(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "pdf->pptx":
		err := PDFtoPPTX(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	case "pdf->xlsx":
		err := PDFtoXLSX(inputFile, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Неподдерживаемый формат конвертации")
	}

	fmt.Printf("Конвертация успешно завершена: %s -> %s\n", inputFile, outputFile)
}

// ImageToPDF конвертирует JPG/PNG в PDF
func ImageToPDF(inputFile, outputFile string) error {
	img, err := loadImage(inputFile)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pageWidth, pageHeight := 210.0, 297.0 // A4 размер в мм
	width, height := float64(img.Bounds().Size().X)*0.264583, float64(img.Bounds().Size().Y)*0.264583

	aspectRatio := width / height
	if aspectRatio > pageWidth/pageHeight {
		height = pageWidth / aspectRatio
		width = pageWidth
	} else {
		width = pageHeight * aspectRatio
		height = pageHeight
	}

	pdf.Image(inputFile, (pageWidth-width)/2, (pageHeight-height)/2, width, height, false, "", 0, "")

	return pdf.OutputFileAndClose(outputFile)
}

// TXTtoPDF конвертирует TXT или DOCX в PDF
func TXTtoPDF(inputFile, outputFile string) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	width, _ := pdf.GetPageSize()
	pdf.MultiCell(width, 10, string(content), "", "L", false)

	return pdf.OutputFileAndClose(outputFile)
}

// PPTXtoPDF конвертирует PowerPoint в PDF
func PPTXtoPDF(inputFile, outputFile string) error {
	ppt, err := presentation.Open(inputFile)
	if err != nil {
		return err
	}

	return ppt.SaveToFile(outputFile)
}

// XLSXtoPDF конвертирует Excel в PDF
func XLSXtoPDF(inputFile, outputFile string) error {
	ss, err := spreadsheet.Open(inputFile)
	if err != nil {
		return err
	}

	return ss.SaveToFile(outputFile)
}

// PDFtoImage конвертирует PDF в JPG/PNG
func PDFtoImage(inputFile, outputFile, format string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.SetResolution(300, 300)
	if err != nil {
		return err
	}

	err = mw.ReadImage(inputFile)
	if err != nil {
		return err
	}
	err = mw.SetImageFormat(format)
	if err != nil {
		return err
	}

	return mw.WriteImage(outputFile)
}

// PDFtoTXT конвертирует PDF в текстовый файл
func PDFtoTXT(inputFile, outputFile string) error {
	doc, err := document.Open(inputFile)
	if err != nil {
		return err
	}

	text := doc.ExtractText().Text()
	return os.WriteFile(outputFile, []byte(text), 0644)
}

// PDFtoPPTX конвертирует PDF в PowerPoint
func PDFtoPPTX(inputFile, outputFile string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.SetResolution(300, 300)
	if err != nil {
		return err
	}

	err = mw.ReadImage(inputFile)
	if err != nil {
		return err
	}

	ppt := presentation.New()
	defer ppt.Close()

	for i := 0; i < int(mw.GetNumberImages()); i++ {
		mw.SetIteratorIndex(i)
		tempFile := fmt.Sprintf("temp_page_%d.jpg", i)
		err := mw.WriteImage(tempFile)
		if err != nil {
			return err
		}
		defer os.Remove(tempFile)

		slide := ppt.AddSlide()
		_, err = slide.AddImageFromFile(tempFile)
		if err != nil {
			return err
		}
	}

	return ppt.SaveToFile(outputFile)
}

// PDFtoXLSX конвертирует PDF в Excel
func PDFtoXLSX(inputFile, outputFile string) error {
	doc, err := document.Open(inputFile)
	if err != nil {
		return err
	}

	text := doc.ExtractText().Text()
	ss := spreadsheet.New()
	defer ss.Close()

	sheet := ss.AddSheet()
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.SetString(line)
		if i > 1000 { // Ограничение на количество строк
			break
		}
	}

	return ss.SaveToFile(outputFile)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return jpeg.Decode(file)
}
