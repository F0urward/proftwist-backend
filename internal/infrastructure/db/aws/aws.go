package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.AWS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, ""),
		Secure: cfg.AWS.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to AWS: %w", err)
	}

	log.Printf("Successfully connected to AWS S3 at %s", cfg.AWS.Endpoint)
	return minioClient, nil
}
