package config

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var CLD *cloudinary.Cloudinary

func Cloudinary() {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	CLD = cld
}
func UploaderCloudinary(file multipart.File, filename string) (string, error) {
	ctx := context.Background()

	res, err := CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: "rag_docs/" + filename,
	})
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}
