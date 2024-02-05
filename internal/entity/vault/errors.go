package vault

import "errors"

var (
	ErrObjectNotExists = errors.New("object does not exist")
	ErrDuplicate       = errors.New("already exists")
	ErrEmptyMetaID     = errors.New("meta id must be non empty")
)
