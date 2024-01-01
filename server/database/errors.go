package database

import (
	"errors"
)

var (
	InvalidInputParamError = errors.New("Invalid input param")
	NoRecordFoundError     = errors.New("No record found")
	CreateTableError       = errors.New("Create table error")
)
