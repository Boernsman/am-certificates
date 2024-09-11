package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"am-certificates/database"
	"am-certificates/models"
	"am-certificates/utils"
)

// ValidateCode validates a certificate code
func ValidateCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	certificate, err := database.GetEntry(code, false)
	if err != nil {
		http.Error(w, "Invalid or used code", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(certificate)
}

// GenerateCertificate generates a certificate
func GenerateCertificate(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	certificate, err := database.GetEntry(request.Code, false)
	if err != nil {
		http.Error(w, "Invalid or used code", http.StatusNotFound)
		return
	}
	// Update certificate with user info
	certificate.Name = request.Name
	certificate.Email = request.Email
	database.AssignEntry(certificate)

	// Generate PDF and send certificate via email
	pdfName := utils.GeneratePDF(certificate.Code)
	pngName := utils.ConvertPDFToPNG(certificate.Code)
	// Send email
	//utils.SendEmail(request.Email, pdfPath)

	pdfUrl := "https://zertifikat.austromagnum.at/img/" + pdfName
	pngUrl := "https://zertifikat.austromagnum.at/img/" + pngName
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
		"code":    certificate.Code,
		"png_url": pngUrl,
		"pdf_url": pdfUrl,
	})
}

func DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	email := r.URL.Query().Get("email")

	certificate, err := database.GetEntry(code, true)
	if err != nil {
		http.Error(w, "Invalid code", http.StatusNotFound)
		return
	}

	if certificate.Email != email {
		http.Error(w, "Invalid email", http.StatusNotFound)
		return
	}
	if err := database.DeleteEntry(code); err != nil {
		log.Println(err)
		http.Error(w, "Invalid code", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Certificate deleted"})
}

// CreateCertificateCodes generates a specified number of unique certificate codes
func CreateCertificateCodes(w http.ResponseWriter, r *http.Request) {
	log.Println("Create certificate codes")
	// Parse the count from query params
	countParam := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countParam)
	if err != nil || count <= 0 {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}
	log.Println("Count:", count)

	// Parse additional parameters (name, type, tags)
	certificateType := r.URL.Query().Get("type")
	if certificateType == "" {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}
	log.Println("Type:", certificateType)
	certificateTags := r.URL.Query().Get("tags")
	log.Println("Tags:", certificateTags)

	// Base URL for the certificate
	baseURL := "https://zertifikat.austromagnum.at/img/?code="

	// List to hold the generated codes and URLs
	var certificateCodes []map[string]string

	for i := 0; i < count; i++ {
		// Generate new ULID code
		newCode := utils.GenerateULID()

		// Create a new certificate entry in the database
		newCertificate := models.Certificate{
			Code:      newCode,
			Type:      certificateType,
			Tags:      certificateTags,
			Generated: false,
		}

		// Save the certificate in the database
		if err := database.DB.Create(&newCertificate).Error; err != nil {
			http.Error(w, "Failed to generate certificate code", http.StatusInternalServerError)
			return
		}

		// Add the generated code and URL to the response
		certificateCodes = append(certificateCodes, map[string]string{
			"code": newCertificate.Code,
			"url":  fmt.Sprintf("%s%s", baseURL, newCertificate.Code),
		})
	}

	// Return the generated codes in JSON
	json.NewEncoder(w).Encode(map[string]interface{}{
		"codes": certificateCodes,
	})
}

func DeleteCertificateCodes(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	log.Println("Delete certificate code:", code)
	if !utils.IsUlid(code) {
		log.Println("Received code is not a valid ULID")
		http.Error(w, "Invalid code", http.StatusNotFound)
		return
	}

	if err := database.DeleteEntry(code); err != nil {
		log.Println(err)
		http.Error(w, "Invalid code", http.StatusNotFound)
		return
	}
}
