package minio

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Dependencies struct {
	MinioService MinioService
}

func InjectDependencies() (*Dependencies, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("MinIO environment variables not set")
	}

	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to initialize MinIO client: %v", err)
		return nil, err
	}

	log.Println("MinIO client successfully initialized")

	// Create MinIO service instance
	minioService := NewMinioService(client)

	return &Dependencies{
		MinioService: minioService,
	}, nil
}
