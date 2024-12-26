package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImage uploads a file to Cloudinary and returns the public URL.
func UploadImage(file *multipart.FileHeader) (string, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatalf("Cloudinary environment variables are not set")
	}

	// Create a new Cloudinary instance
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("error creating Cloudinary instance: %v", err)
		return "", fmt.Errorf("error creating Cloudinary instance: %w", err)
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file: %v", err)
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Upload the file to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder: "your-folder-name", // Specify folder in Cloudinary (optional)
	}
	uploadResult, err := cld.Upload.Upload(context.TODO(), f, uploadParams)
	if err != nil {
		log.Printf("error uploading file to Cloudinary: %v", err)
		return "", fmt.Errorf("error uploading file to Cloudinary: %w", err)
	}

	// Return the secure URL of the uploaded file
	return uploadResult.SecureURL, nil
}
