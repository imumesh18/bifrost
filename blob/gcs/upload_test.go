package gcs

import (
	"bytes"
	"context"
	"testing"

	"github.com/imumesh18/bifrost/blob"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	client, err := New(context.Background(), WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "bankopen-uat",
		"private_key_id": "2164d78facb058ac7e6ac546329e3e982422615c",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCpewmjPAWG1LVe\ntXEgB+hfiI5uClmzH4Ef6E1FruZ+hdBfSa1hrchKnneJCez5F/YXyAZ0HXYkB0tg\nHQg0eEjHwvO30PhHsoA9aZpyI+Emgj5OmgalRBhJIsktKVUsrpVTh4j3u4PbBsof\nENufHWvpo2m31roKB353zY6hkP3j6Ieq6M9h5qhc6+qWnwE7NZgnTjGSbP56/0kd\nTCfMBEF8aiaIXG52fs9QfVW9HUL9wvCm7Jcyrm4wvhGqFX8NZTNBXCAeUJzvpWTz\n64TkiyBlCi3yMCj9TJH9NVEYhhOu5vqQS2pHDmSI4vsyYuYyddTZdZC1Ezu3h7fr\nz9ICy0m5AgMBAAECggEABCjIOeS0VISmHljI2dLD28tZxn7TyHr4k+hDNDd4geim\nDl6cJhJqINWuhJMVKSBSym7e5Kzf7D6vFYDk5YfBroBc85kJ7NWHQuqTbyNtCEFI\nRuzOBrKK/4l+YzC5Xcc7lDEWdhD4qZCLN5O19hvWXQWlB61aA7q/GnDVpCavKm+N\nUS+pC52eblVL1uw2VACHlksyLRGjJ3X0CDHiyBZd+moQiAL/D4imtZOvC5mSTIBR\nzYNA1X69guFSe5246zcknW8gmvtYj4R2kOe4wxkucVQhh5zl51vU4TdlGwqS5oHz\ncCJslhD3eEi1m7HXdCE541Dl4jy0Bc91MoE6JGVEoQKBgQDpqSOBg86VM1+INjjr\nDxv496KE/rPtN+bjeYbWIC73C8rqpW407uATxvTGh/efX0UtUqGB1AzruXCEPH3m\nhFnuOjkQshRfrHFIL4NkElqhfGzGQy6c0oE3fmNh+qSZinepdLPe/OHcz0LsEoF6\ngYPydYedQ83kRAn7IqNipb1nGQKBgQC5rxZNmtHlDSQF3aKMzAuuNz7jbXP+SSV9\nWhbvgmLmKOkYlDzw7Mvqyxt4inY2rX5pNTUA6r7f+i87qnsKK6SdNsTTwfRUNEw7\nlkMilKQBKkycTyOvqQkxEIiL8i5Izrosd0Ysueup9zvrpUU27pLRsOaeDD7p5wJ2\nJD+7AgZroQKBgBgCRyGxt3JhOvm2CJcukEM+vrZHrZk8Wz8YZ6Bs4iaEUa9WnEJY\nITInCVO0+N6pXWRQz0OV1FYMUeFkjdM32j2+QcrTYYCLKYCvUSLhN+rL7ClbEdkP\nUDOxiuiwZmVYcv84fJr3BQY5TbkQFbnOwQ4SwYKJSwifbR8e6gbi3NlRAoGAHvKi\nqf6S2zVMernM/OCJVdkZXzh/67LvT6wzRGob57aL2y/h1FnzRsfhZT7WoxhZiFl4\n4xU9CQGe27f3V+OcRSO6vHyIJ3yr9AaAXAQgLZ2KNUcvcHig8o+J4qFTu4jRGNYs\nWQoH0EVHtGfQWG59A/wTA+aQmdWJ4Hz8LkQRI+ECgYEAlGewysqDfbg8MBLUeloS\nBomph0Fg5dWli1vae7i0wBePnjB70EB1hAMjar+4Lchy1Er0G5RxkOJZq0GIF6Jp\nZc+mw7CKYzZ5qdQ779jcFuZaGMF5FQWwj+EzDEmm+2cXS8kvePalmc+1BvIrn6gR\nGBtnj6l1f8V1HLXmFza9bU0=\n-----END PRIVATE KEY-----\n",
		"client_email": "uat-zwitch-reports@bankopen-uat.iam.gserviceaccount.com",
		"client_id": "111064628825600123658",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/uat-zwitch-reports%40bankopen-uat.iam.gserviceaccount.com"
	  }`)), WithProjectID("bankopen-uat"))
	require.NoError(t, err)
	_, err = client.Upload(context.Background(), blob.WithUploadBody(bytes.NewBufferString("hello world")), blob.WithUploadBucket("sample-bucket"), blob.WithUploadObject("hello"))
	require.NoError(t, err)
}
