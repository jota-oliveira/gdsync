// Package google - Client for Google Drive

package google

import (
	"context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveClient struct {
	service *drive.Service
}

func NewDriveClient(credentialsPath string) (*DriveClient, error) {
	ctx := context.Background()

	srv, err := drive.NewService(
		ctx,
		option.WithCredentials(credentialsPath),
		option.WithScopes(drive.DriveScope),
	)
	if err != nil {
		return nil, err
	}

	return &DriveClient{
		service: srv,
	}, nil
}
