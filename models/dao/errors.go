package dao

import "fmt"

type BaseError interface {
	error
	StatusCode() int
}

type ModelNotFound struct {
	Id    int
	Model string
	Code  int
}

func (e *ModelNotFound) Error() string {
	return fmt.Sprintf("No query results for model [%s] %d", e.Model, e.Id)
}

func (e *ModelNotFound) StatusCode() int {
	return e.Code
}
