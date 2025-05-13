package myerrors

import "errors"

var ErrUserAlreadyExists = errors.New("the user with such data already exists: try changing the username or email")

var ErrIncorrectPassword = errors.New("incorrect password provided")

var ErrUserNotFound = errors.New("user not found")

var ErrFieldUserEmpty = errors.New("empty user field")
