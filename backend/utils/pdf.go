package utils

import (
	"am-certificates/database"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var (
	CertificateFolder string
	TemplateFolder    string
)

func pdfFileName(code string) string {
	return fmt.Sprintf("%s.pdf", code)
}

func pngFileName(code string) string {
	return fmt.Sprintf("%s.png", code)
}

// GeneratePDF creates a certificate with the user's name
func GeneratePDF(code string) string {

	cert, err := database.GetEntry(code, true)
	if err != nil {
		return ""
	}
	pdf := gofpdf.New("L", "mm", "A4", TemplateFolder)
	pdf.AddUTF8Font("CutiveMono", "", "CutiveMono-Regular.ttf")
	pdf.AddPage()
	pdf.SetAutoPageBreak(false, 0)

	pdf.SetFont("CutiveMono", "", 40)
	background := filepath.Join(TemplateFolder, "background.png")
	pdf.ImageOptions(background, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pageWidth, pageHeight := pdf.GetPageSize()
	// Vertical position offset to start from
	yOffset := 110.0
	textWidth := pdf.GetStringWidth(cert.Name)
	// Set the X position to center the text
	pdf.SetXY((pageWidth-textWidth)/2, yOffset)
	// Add the NAME
	pdf.CellFormat(textWidth, 10, cert.Name, "", 1, "C", false, 0, "")
	yOffset += 26
	// Add the TYPE
	textWidth = pdf.GetStringWidth(cert.Type)
	pdf.SetXY((pageWidth-textWidth)/2, yOffset)
	pdf.CellFormat(textWidth, 10, cert.Type, "", 1, "C", false, 0, "")
	pdf.SetFont("CutiveMono", "", 16)
	// Add the DATE
	yOffset = pageHeight - 25
	pdf.SetXY(5, yOffset)
	textWidth = pdf.GetStringWidth(cert.Date)
	pdf.CellFormat(textWidth, 10, cert.Date, "", 0, "L", false, 0, "")
	// Add the CODE
	textWidth = pdf.GetStringWidth(cert.Code)
	pdf.SetXY(pageWidth-(textWidth+7), yOffset)
	pdf.CellFormat(textWidth, 10, cert.Code, "", 0, "R", false, 0, "")

	pdf.SetProducer("Austro Magnum", true)
	pdf.SetSubject("Zertifikat", true)
	pdf.SetAuthor("Austro Magnum", true)
	pdf.SetCreator("Austro Magnum", true)
	pdf.SetCreationDate(time.Now())

	fileName := pdfFileName(cert.Code)
	p := filepath.Join(CertificateFolder, fileName)
	err = pdf.OutputFileAndClose(p)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
	}
	return fileName
}

// ConvertPDFToPNG converts a PDF file to PNG using pdftoppm.
func ConvertPDFToPNG(code string) string {
	// Make sure the output directory exists
	if _, err := os.Stat(CertificateFolder); err != nil {
		log.Fatalf("Certificate folder does not exist: %v", err)
	}
	pdf := CertificateFolder + "/" + pdfFileName(code)
	if _, err := os.Stat(pdf); os.IsNotExist(err) {
		log.Println("PDF file does not exist:", pdf)
		return ""
	}

	// Use pdftoppm to convert PDF to PNG
	cmd := exec.Command("pdftoppm", "-singlefile", "-png", pdfFileName(code), code)
	cmd.Dir = CertificateFolder
	if err := cmd.Run(); err != nil {
		log.Printf("failed to convert PDF to PNG: %v", err)
		return ""
	}

	log.Println("PDF successfully converted to PNG")
	return pngFileName(code)
}
