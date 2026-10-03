package custom_error

import "errors"

var NoRowsAffected = errors.New("failed to update")
