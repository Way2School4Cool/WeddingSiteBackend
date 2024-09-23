package main

import (
	"fmt"
    "io"
	"net/http"
	"os"
	"path/filepath"
    "time"
)

const uploadPath = "./uploads"

// uploadHandler handles the file upload
func uploadHandler(response http.ResponseWriter, request *http.Request) {
    start := time.Now()
    var checkpoint time.Time

    // Set CORS headers
    response.Header().Set("Access-Control-Allow-Origin", "*")  // Allow all origins; for production, specify the allowed domain
    response.Header().Set("Access-Control-Allow-Methods", "POST")
    response.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if request.Method == http.MethodOptions {
        response.WriteHeader(http.StatusOK)  // Handle preflight requests
        return
    }

    // Parse the multipart form
    err := request.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(response, "Could not parse form", http.StatusBadRequest)
        return
    }

    if request.Method == http.MethodPost {
        file, header, err := request.FormFile("image")
            if err != nil {
            http.Error(response, "Error retrieving the file", http.StatusBadRequest)
            return
        }
		defer file.Close()

        // Limit file size to 100MB
        if header.Size > 100*1024*1024 {
            http.Error(response, "file is too large", http.StatusBadRequest)
            return
        }

        // Restrict file types to images only
        allowedTypes := map[string]bool{
            "image/gif": true,
            "image/heif": true,
            "image/jpeg": true,
            "image/raw": true,
            "image/png":  true,
            "image/webp": true,
        }

        fileType := header.Header.Get("Content-Type")
        if !allowedTypes[fileType] {
            http.Error(response, "invalid file type", http.StatusBadRequest)
            return
        }
            
		// Create the uploads directory if it doesn't exist
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			err := os.Mkdir(uploadPath, os.ModePerm)
			if err != nil {
				http.Error(response, "Unable to create upload directory\n", http.StatusInternalServerError)
				return
			}
		}

        // TODO: Hash images to prevent repeats
        
        
        // TODO: Compress files


		
        // Create a file in the uploads directory
		destFile, err := os.Create(filepath.Join(uploadPath, time.Now().String()))
		if err != nil {
			http.Error(response, "Unable to create file\n", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

        // TODO: Reformat images to webp for size

		// Copy the uploaded file to the destination file
		_, err = io.Copy(destFile, file)
		if err != nil {
			http.Error(response, "Unable to save file\n", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(response, "File successfully uploaded\n")
	} else {
		http.Error(response, "Invalid request method\n", http.StatusMethodNotAllowed)
	}

    checkpoint = time.Now()
    fmt.Printf("Saved to file @ %s\n\tSaved in: %v\n", time.Now().String(), checkpoint.Sub(start));
}

func main() {
	http.HandleFunc("/uploadimage", uploadHandler)
	fmt.Println("Server started at http://localhost:8085")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}

