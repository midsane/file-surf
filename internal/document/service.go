package document

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/midsane/file-surf/internal/storage"
	"github.com/midsane/file-surf/internal/tenant"
)

type Service struct {
	repo       *Repository
	tenantRepo *tenant.Repository
	storage    *storage.S3Storage
}

func NewService(repo *Repository, tenantRepo *tenant.Repository, storage *storage.S3Storage) *Service {
	return &Service{
		repo:       repo,
		tenantRepo: tenantRepo,
		storage:    storage,
	}
}

func (s *Service) UploadDocument(ctx context.Context, tenantID string, fileHeader *multipart.FileHeader) (*Document, error) {
	// Validate tenant
	t, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("tenant not found")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	docID := uuid.New().String()
	key := fmt.Sprintf("%s/%s", tenantID, docID)

	if err := s.storage.Upload(ctx, key, file, fileHeader.Size, fileHeader.Header.Get("Content-Type")); err != nil {
		return nil, err
	}

	doc := &Document{
		ID:       docID,
		TenantID: tenantID,
		FileName: fileHeader.Filename,
		S3Key:    key,
		Size:     fileHeader.Size,
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Service) GetDocuments(ctx context.Context, tenantID string) ([]Document, error) {
	return s.repo.GetByTenant(ctx, tenantID)
}