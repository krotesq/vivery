package target

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/krotesq/vivery/internal/config"
)

type service struct {
	repository *repository
	config     *config.Config
}

func newService(repository *repository, c *config.Config) *service {
	return &service{
		repository: repository,
		config:     c,
	}
}

func (s *service) findAll(ctx context.Context, accountID string) ([]target, error) {
	return s.repository.findAll(ctx, accountID)
}

func (s *service) findWithRtmpByID(ctx context.Context, id string, accountID string) (*target, *rtmp, error) {
	return s.repository.findWithRtmpByID(ctx, id, accountID)
}

func (s *service) createWithRtmp(ctx context.Context, name string, description string, url string, streamKey string, accountID string) (*target, *rtmp, error) {
	if err := s.validateURL(url, []string{"rtmp", "rtmps"}); err != nil {
		return nil, nil, err
	}

	return s.repository.createWithRtmp(ctx, name, description, url, streamKey, accountID)
}

func (s *service) deleteByID(ctx context.Context, id string, accountID string) error {
	return s.repository.deleteByID(ctx, id, accountID)
}

func (s *service) validateURL(rawURL string, allowedSchemes []string) error {
	defaultErr := errors.New("invalid url")

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if !slices.Contains(allowedSchemes, parsed.Scheme) {
		return defaultErr
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))

	publicIP := strings.ToLower(strings.TrimSpace(s.config.PublicIP))
	if host == publicIP {
		return defaultErr
	}

	if slices.Contains(strings.Split(s.config.Domains, ","), host) {
		return defaultErr
	}

	return nil
}
