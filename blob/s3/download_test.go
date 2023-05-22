package s3

import (
	"context"
	"testing"

	"github.com/imumesh18/bifrost/blob"
)

func TestDownload(t *testing.T) {
	url := "http://localhost:4566"
	region := "us-east-1"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "upload",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := New(context.Background(), WithRegion(region), WithURL(url), WithAccessKeyID("hello"), WithSecretAccessKey("world"))
			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if _, err := c.Download(context.Background(), blob.WithDownloadBucket("hello"), blob.WithDownloadObject("world")); err != nil {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
