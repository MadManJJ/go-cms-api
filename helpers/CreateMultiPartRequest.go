package helpers

import (
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateMultipartRequest(t *testing.T) (*bytes.Buffer, string) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", "test.png")
	require.NoError(t, err)

	_, err = fileWriter.Write([]byte("dummy image data"))
	require.NoError(t, err)

	// Add optional fields
	_ = writer.WriteField("path", "images/uploads")
	_ = writer.WriteField("replace", "true")

	writer.Close()
	return &body, writer.FormDataContentType()
}
