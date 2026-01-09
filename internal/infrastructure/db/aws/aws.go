package aws

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(cfg *config.Config) *minio.Client {
	const op = "aws.NewClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	minioClient, err := minio.New(cfg.AWS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AWS.MinioRootUser, cfg.AWS.MinioRootPassword, ""),
		Secure: cfg.AWS.UseSSL,
	})
	if err != nil {
		logger.WithError(err).Error("failed to create AWS client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		logger.WithError(err).Error("cannot connect to AWS")
	}

	return minioClient
}
