package custom_error

import "errors"

var ChangeUserInvalidPassword = errors.New("Invalid password")
var ChangeUserrInvalidLength = errors.New("Password at least 6 character")
var AllFieldIsEmpty = errors.New("All field is empty")
