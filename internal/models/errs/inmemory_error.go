package errs

import "errors"

var (
	ErrQuotesIsEmpty = errors.New("quotes list is empty")
	ErrQuoteNotFound = errors.New("quote not found")
)
