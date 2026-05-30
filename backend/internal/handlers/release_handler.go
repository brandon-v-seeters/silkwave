package handlers

import (
	"github.com/brandon-v-seeters/go-silk-wave/internal/access"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/brandon-v-seeters/go-silk-wave/internal/storage"
)

const presignedURLExpireSeconds int64 = 3600 // 1 hour
const maxAvatarUploadBytes int64 = 5 * 1024 * 1024

type ReleaseHandler struct {
	db            *database.ArangoDB
	accessSource  access.Source
	releases      *repository.ReleaseRepository
	storageClient *storage.Client
}

func NewReleaseHandler(db *database.ArangoDB, accessSource access.Source, releases *repository.ReleaseRepository, storageClient *storage.Client) *ReleaseHandler {
	return &ReleaseHandler{db: db, accessSource: accessSource, releases: releases, storageClient: storageClient}
}
