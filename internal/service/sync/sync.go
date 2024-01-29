package sync

import "context"

type Service struct {
	client  client
	storage storage
}

func New() *Service {
	return &Service{}
}

func Run(ctx context.Context) error {
	return nil
}
