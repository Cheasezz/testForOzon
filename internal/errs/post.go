package errs

import "errors"

var (
	ErrCreatePost    = errors.New("create post error")
	ErrPostDsntExist = errors.New("post does not exist ")
)
