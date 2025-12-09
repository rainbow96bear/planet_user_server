package repository

import "github.com/rainbow96bear/planet_user_server/internal/planet_err"

var (
	ErrNotFound          = planet_err.ErrNotFound
	ErrAlreadyExists     = planet_err.ErrAlreadyExists
	ErrNicknameDuplicate = planet_err.ErrNicknameDuplicate
)
