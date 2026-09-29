package custom_error

import "errors"

var RegisterInvalidLength = errors.New("email and password at least 6 character")
var RegisterAlreadyExist = errors.New("user already exist")
var LoginInvalidEmailOrPassword = errors.New("invalid email or password")
var EmptyLoginField = errors.New("username or password cannot be empty")
var NotificationNotFound = errors.New("you don't have notification yet")
