package controller

import "errors"

// ErrHasChildren is returned when a resource has dependent child entities.
var ErrHasChildren = errors.New("resource has children")
