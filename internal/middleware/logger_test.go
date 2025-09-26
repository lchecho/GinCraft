package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldLogRequestBody(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{
			name:        "JSON content type should be logged",
			contentType: "application/json",
			expected:    true,
		},
		{
			name:        "Form data should be logged",
			contentType: "application/x-www-form-urlencoded",
			expected:    true,
		},
		{
			name:        "Multipart form data (file upload) should NOT be logged",
			contentType: "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW",
			expected:    false,
		},
		{
			name:        "Empty content type should not be logged",
			contentType: "",
			expected:    false,
		},
		{
			name:        "Plain text should not be logged",
			contentType: "text/plain",
			expected:    false,
		},
		{
			name:        "Image content should not be logged",
			contentType: "image/jpeg",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldLogRequestBody(tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsFileUpload(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{
			name:        "Multipart form data should be detected as file upload",
			contentType: "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW",
			expected:    true,
		},
		{
			name:        "Simple multipart form data should be detected as file upload",
			contentType: "multipart/form-data",
			expected:    true,
		},
		{
			name:        "JSON should not be detected as file upload",
			contentType: "application/json",
			expected:    false,
		},
		{
			name:        "Form data should not be detected as file upload",
			contentType: "application/x-www-form-urlencoded",
			expected:    false,
		},
		{
			name:        "Empty content type should not be detected as file upload",
			contentType: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFileUpload(tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShouldLogResponseBody(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{
			name:        "JSON response should be logged",
			contentType: "application/json",
			expected:    true,
		},
		{
			name:        "HTML response should be logged",
			contentType: "text/html",
			expected:    true,
		},
		{
			name:        "Empty content type should be logged by default",
			contentType: "",
			expected:    true,
		},
		{
			name:        "Binary file should not be logged",
			contentType: "application/octet-stream",
			expected:    false,
		},
		{
			name:        "Image should not be logged",
			contentType: "image/jpeg",
			expected:    false,
		},
		{
			name:        "Video should not be logged",
			contentType: "video/mp4",
			expected:    false,
		},
		{
			name:        "Audio should not be logged",
			contentType: "audio/mpeg",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldLogResponseBody(tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
