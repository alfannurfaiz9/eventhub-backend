package custom_error

import "errors"

var EventNotFound = errors.New("event not found")
var EventErrorPage = errors.New("page must be a positive number")
