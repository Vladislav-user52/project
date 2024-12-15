package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func main() {

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
