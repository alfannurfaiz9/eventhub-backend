package custom_error

import "errors"

var ChangeUserInvalidPassword = errors.New("invalid password")
var ChangeUserrInvalidLength = errors.New("password at least 6 character")
