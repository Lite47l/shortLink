package model

import "errors"

var (
	ErrNotFound          = errors.New("Not Found")
	ErrInvalidUrl        = errors.New("Invalid URL")
	ErrShortUrlCollision = errors.New("Short URL Collision error")
)
