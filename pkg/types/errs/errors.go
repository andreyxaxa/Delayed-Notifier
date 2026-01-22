package errs

import "errors"

var (
	ErrRecordAlreadyExists               = errors.New("record already exists")
	ErrRecordNotFound                    = errors.New("record not found")
	ErrRecordNotFoundOrAlreadyProcessing = errors.New("record not found or already processing")
	ErrFutureNotification                = errors.New("notification scheduled for future")
	ErrAlreadyCancelled                  = errors.New("notification already cancelled")
	ErrAlreadySentOrFailed               = errors.New("notification already send or failed")
	ErrInvalidChatID                     = errors.New("invalid chat ID")
)
