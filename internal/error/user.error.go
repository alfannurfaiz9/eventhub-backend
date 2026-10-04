package custom_error

import "errors"

var ChangeUserInvalidPassword = errors.New("invalid password")
var ChangeUserrInvalidLength = errors.New("password at least 6 character")
var AllFieldIsEmpty = errors.New("all field is empty")
