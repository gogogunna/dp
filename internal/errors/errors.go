package errors

import "errors"

var NotFoundAccount = errors.New("accounts for this token don't exist")
