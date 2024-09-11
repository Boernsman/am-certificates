package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"am-certificates/database"
	"am-certificates/handlers"
	"am-certificates/middleware"
	"am-certificates/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Define command-line flags
	certificatesFolder := flag.String("certificates", "/data/austromagnum/cert", "Folder to serve certificates from")
	appFolder := flag.String("app", "../frontend/", "Folder to serve the app from")
	credentialFile := flag.String("cred", "../tests/config.ini", "File containing API key and basic auth credentials")
	certFile := flag.String("cert", "/etc/ssl/certs/public.crt", "Path to the SSL certificate file (required for HTTPS)")
	keyFile := flag.String("key", "/etc/ssl/private/private.key", "Path to the SSL private key file (required for HTTPS)")
	databaseFile := flag.String("database", "/data/austromagnum/certificates.db", "Database file")
	httpsPort := flag.String("https-port", "443", "Port for the server to listen on")
	httpPort := flag.String("http-port", "80", "Port for the HTTP server to listen on for redirection")

	// Parse the command-line flags
	flag.Parse()

	// Check if the certificates folder exists, if not create it
	if _, err := os.Stat(*certificatesFolder); os.IsNotExist(err) {
		err = os.Mkdir(*certificatesFolder, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create certificates folder: %v", err)
		}
	}
	utils.CertificateFolder = *certificatesFolder
	utils.TemplateFolder = "/data/austromagnum/template"

	// Load credentials for basic authentication
	if err := middleware.LoadCredentials(*credentialFile); err != nil {
		log.Fatalf("Failed to load authentication credentials %v", err)
	}

	// Initialize the database
	database.ConnectDatabase(*databaseFile)
	log.Println("Database stats:")
	log.Println("    - Total entries:", database.GetTotalEntries())
	log.Println("    - Unused entries:", database.GetUnusedEntries())
	for k, v := range database.GetEntriesByType() {
		log.Printf("    - Type \"%s\" unused %d\n", k, v)
	}

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Set up routes
	r.HandleFunc("/valide", handlers.ValidateCode).Methods("GET")
	r.HandleFunc("/generiere", handlers.GenerateCertificate).Methods("POST")

	// Protected route for /erstellen using Basic Auth
	createRoute := r.Path("admin").Subrouter()
	createRoute.Use(middleware.BasicAuthMiddleware)
	r.HandleFunc("/erstelle", handlers.CreateCertificateCodes).Methods("GET")
	r.HandleFunc("/loesche", handlers.DeleteCertificateCodes).Methods("GET", "DELETE")

	// Serve static files (certificates folder and the web app)
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir(*certificatesFolder))))
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(*appFolder))))

	// Validate the certificate and key paths
	if *certFile == "" || *keyFile == "" {
		log.Print("SSL certificate and key files are required for HTTPS. Use --cert and --key flags to specify the paths.")
	} else {
		// Start the HTTPS server in a goroutine
		go func() {
			address := fmt.Sprintf(":%s", *httpsPort)
			log.Printf("Starting HTTPS server on port %s", *httpsPort)
			if err := http.ListenAndServeTLS(address, *certFile, *keyFile, r); err != nil {
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		}()
	}

	// Start the HTTP server for redirection to HTTPS
	httpAddress := fmt.Sprintf(":%s", *httpPort)
	log.Printf("Starting HTTP server on port %s for redirecting to HTTPS", *httpPort)
	if err := http.ListenAndServe(httpAddress, http.HandlerFunc(redirectToHTTPS)); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// redirectToHTTPS redirects HTTP requests to HTTPS
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	httpsURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
}
