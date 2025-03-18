package google

import (
	"context"
	"io"
	"mime/multipart"
	"net/url"

	"cloud.google.com/go/storage"
)

func UploadFileToGCS(ctx context.Context, bucket string, path string, f multipart.File) (string, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	sw := storageClient.Bucket(bucket).Object(path).NewWriter(ctx)

	if _, err := io.Copy(sw, f); err != nil {
		return "", err
	}

	if err := sw.Close(); err != nil {
		return "", err
	}

	gcsPath := "/" + bucket + "/" + sw.Attrs().Name

	finalPath, err := url.Parse(gcsPath)
	return finalPath.EscapedPath(), err
}
