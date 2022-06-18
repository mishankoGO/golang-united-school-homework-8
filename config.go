package main

import (
	"errors"
	"fmt"
)

// define constants
const Permission = 0644

// define errors and necessary strings
var OperationMissing = errors.New("-operation flag has to be specified")
var FileNameMissing = errors.New("-fileName flag has to be specified")
var OperationNotAllowed = func(s string) error {
	return errors.New(fmt.Sprintf("Operation %s not allowed!", s))
}
var ItemMissing = errors.New("-item flag has to be specified")
var IdMissing = errors.New("-id flag has to be specified")
var IdExists = func(s string) string {
	return fmt.Sprintf("Item with id %s already exists", s)
}
var IdNotFound = func(s string) string {
	return fmt.Sprintf("Item with id %s not found", s)
}
