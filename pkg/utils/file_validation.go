package utils

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

// IsValidImageFile - Validates if uploaded file is a valid image
// Tujuan: Memvalidasi apakah file yang diupload adalah gambar yang valid
// Parameter: file (*multipart.FileHeader) - file header dari multipart form
// Return: bool - true jika file valid, false jika tidak
// Penjelasan: Memeriksa ekstensi file dan MIME type untuk memastikan file adalah gambar
func IsValidImageFile(file *multipart.FileHeader) bool {
    if file == nil {
        return false
    }

    // Check file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
    
    validExt := false
    for _, validExtension := range validExtensions {
        if ext == validExtension {
            validExt = true
            break
        }
    }
    
    if !validExt {
        return false
    }

    // Check MIME type
    validMimeTypes := []string{
        "image/jpeg",
        "image/jpg", 
        "image/png",
        "image/gif",
        "image/webp",
    }

    contentType := file.Header.Get("Content-Type")
    for _, validMimeType := range validMimeTypes {
        if contentType == validMimeType {
            return true
        }
    }

    return false
}