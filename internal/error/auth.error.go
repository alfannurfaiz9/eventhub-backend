package custom_error

import "errors"

var RegisterInvalidLength = errors.New("Email and password at least 6 character")
var ForgotPasswordInvalidLength = errors.New("New password at least 6 character")
var RegisterAlreadyExist = errors.New("User already exist")
var UserNotFound = errors.New("User not found")
var LoginInvalidEmailOrPassword = errors.New("Invalid email or password")
var EmptyLoginField = errors.New("Username or password cannot be empty")
var NotificationNotFound = errors.New("you don't have notification yet")
