package storage

import (
	"mime"
	"path/filepath"
)

func GetExtensionFromContentType(contentType string) string {
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		return ""
	}
	// Remove the leading dot
	return exts[0][1:]
}

func GetExtensionFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}
