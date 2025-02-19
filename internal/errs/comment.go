package errs

import "errors"

var (
	ErrOffsetToBig        = errors.New("offset for commentaries pagination is to big")
	ErrStartTransaction   = errors.New("failed to start transaction")
	ErrCheckCmtAllowed    = errors.New("failed to check comments_allowed")
	ErrCmntDisabled       = errors.New("comments are disabled for this post")
	ErrFailedCommit       = errors.New("failed to commit transaction")
	ErrFailedParentCmnt   = errors.New("failed to check parent comment")
	ErrParentCmntNotExist = errors.New("parent comment does not exist")

	ErrToLongtext  = errors.New("comment is too long (max 2000 characters)")
	ErrCmntAreProh = errors.New("comments are prohibited")
)
