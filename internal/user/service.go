package user

import (
	"context"
	"github.com/OlegKapat/Rest-api-mongo/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (u User, r error) {
	return
}
