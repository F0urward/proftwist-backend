package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/auth"
)

type AuthAWSRepository struct {
	client *minio.Client
}

func NewAuthAWSRepository(awsClient *minio.Client) auth.AWSRepository {
	return &AuthAWSRepository{client: awsClient}
}

func (aws *AuthAWSRepository) PutObject(ctx context.Context, input entities.UploadInput) (*minio.UploadInfo, error) {
	const op = "AuthAWSRepository.PutObject"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"bucket_name": input.BucketName,
		"file_name":   input.Name,
		"file_size":   input.Size,
	})

	fileName := aws.generateFileName(input.Name)

	options := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := aws.client.PutObject(ctx, input.BucketName, fileName, input.File, input.Size, options)
	if err != nil {
		logger.WithError(err).Error("failed to put object")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("generated_file_name", fileName).Info("successfully uploaded object")
	return &uploadInfo, nil
}

func (aws *AuthAWSRepository) GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error) {
	const op = "AuthAWSRepository.GetObject"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"bucket":    bucket,
		"file_name": fileName,
	})

	object, err := aws.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		logger.WithError(err).Error("failed to get object")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved object")
	return object, nil
}

func (aws *AuthAWSRepository) RemoveObject(ctx context.Context, bucket string, fileName string) error {
	const op = "AuthAWSRepository.RemoveObject"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"bucket":    bucket,
		"file_name": fileName,
	})

	if err := aws.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{}); err != nil {
		logger.WithError(err).Error("failed to remove object")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully removed object")
	return nil
}

func (aws *AuthAWSRepository) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
