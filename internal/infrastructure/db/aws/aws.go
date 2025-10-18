package aws

import (
	"context"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(cfg *config.Config) *minio.Client {
	log.Println(cfg.AWS.MinioRootUser, cfg.AWS.MinioRootPassword)
	minioClient, err := minio.New(cfg.AWS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AWS.MinioRootUser, cfg.AWS.MinioRootPassword, ""),
		Secure: cfg.AWS.UseSSL,
	})
	if err != nil {
		log.Fatalf("failed to create AWS client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("cannot connect to AWS: %v", err)
	}

	return minioClient
}
