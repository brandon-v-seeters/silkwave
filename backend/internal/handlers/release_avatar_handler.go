package handlers

import (
	"net/http"
	"strings"

	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *ReleaseHandler) UploadAvatar(c *gin.Context) {
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	var req models.CreateAvatarUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if req.FileSize <= 0 || req.FileSize > maxAvatarUploadBytes {
		respondError(c, http.StatusBadRequest, "invalid_request", "Avatar image must be 5MB or smaller")
		return
	}
	if !strings.HasPrefix(strings.ToLower(req.FileType), "image/") {
		respondError(c, http.StatusBadRequest, "invalid_request", "Avatar upload must be an image")
		return
	}

	ext := storage.GetExtensionFromFileName(req.FileName)
	if ext == "" {
		ext = storage.GetExtensionFromContentType(req.FileType)
	}
	ext = strings.ToLower(ext)
	if !isAllowedAvatarExtension(ext) {
		respondError(c, http.StatusBadRequest, "invalid_request", "Avatar image must be png, jpg, jpeg, or webp")
		return
	}

	storagePath := h.storageClient.Resolver().Avatar(userKey, ext)
	presignedURL, err := h.storageClient.GetUploadAvatarURL(
		c.Request.Context(),
		userKey,
		ext,
		req.FileType,
		presignedURLExpireSeconds,
	)
	if err != nil {
		logger.Error("failed to generate avatar upload URL", err, zap.String("userKey", userKey))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to generate avatar upload URL")
		return
	}

	respondOK(c, models.CreateAvatarUploadResponse{
		PresignedUrl: presignedURL,
		StoragePath:  storagePath,
	}, "")
}

func isAllowedAvatarExtension(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "webp":
		return true
	default:
		return false
	}
}
