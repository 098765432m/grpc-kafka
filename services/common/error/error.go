package common_error

import "errors"

var ErrDuplicateRecord = errors.New("duplicate record")
var ErrNoRows = errors.New("no rows found")
