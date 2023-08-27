package todo

import (
	"github.com/google/uuid"
)

type Todo struct {
	Id          uuid.UUID
	Description string
}

func New(description string) *Todo {
	return &Todo{
		Id:          uuid.New(),
		Description: description,
	}
}
