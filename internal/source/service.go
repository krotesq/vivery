package source

import (
	"context"
)

type service struct {
	repository *repository
}

func newService(repository *repository) *service {
	return &service{repository: repository}
}

func (service *service) findAll(ctx context.Context, accountID string) ([]source, error) {
	return service.repository.findAll(ctx, accountID)
}

func (service *service) findWithRtmpByID(ctx context.Context, id string, accountID string) (*source, *rtmp, error) {
	return service.repository.findWithRtmpByID(ctx, id, accountID)
}

func (service *service) createWithRtmp(ctx context.Context, name string, description string, url string, streamKey string, accountID string) (*source, *rtmp, error) {
	return service.repository.createWithRtmp(ctx, name, description, url, streamKey, accountID)
}

func (service *service) deleteByID(ctx context.Context, id string, accountID string) error {
	return service.repository.deleteByID(ctx, id, accountID)
}
