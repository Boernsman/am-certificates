# Austro Magnum Zertifikat Backend


## Usage

Command-Line Flags:

    flag.String("port", "8080", ...): Defines the --port argument, which allows the user to specify the port (default is 8080).
    flag.String("certificates", "./certificates", ...): Defines the --certificates argument, allowing the user to specify the folder where the certificates are stored (default is ./certificates).

Folder Check:

    The code checks if the provided certificates folder exists. If it doesn’t exist, it creates the folder.

Serving Static Files:

    The r.Static("/certificates", *certificatesFolder) serves the static files from the certificates folder provided via the command-line argument.

Starting the Server:

    The server starts on the port specified by the --port argument.

## Build 

Install the required libraries:

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/jung-kurt/gofpdf
go get github.com/go-gomail/gomail
```


```bash
go run main.go
```


## Docker


