package s3

import (
	"bytes"
	"context"
	"testing"

	"github.com/imumesh18/bifrost/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		expectedError  error
		expectedObject *blob.Object
		name           string
	}{
		{
			name:           "successfully able to upload file to s3",
			expectedError:  nil,
			expectedObject: &blob.Object{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, err := New(
				ctx,
				WithAccessKeyID("test"),
				WithSecretAccessKey("test"),
				WithRegion("us-east-1"),
				WithEndpoint("http://localhost:4566"),
			)
			require.NoError(t, err)

			actualObject, err := storage.Upload(
				ctx,
				blob.WithUploadBody(bytes.NewReader([]byte("test"))),
				blob.WithUploadBucket("test"),
				blob.WithUploadObject("test"),
			)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedObject, actualObject)
			}

		})
	}
}
