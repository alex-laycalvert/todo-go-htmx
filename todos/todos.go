package todos

import (
	"github.com/google/uuid"
)

type Todo struct {
	Id          uuid.UUID
	Description string
}

var todos []*Todo

func Todos() []*Todo {
	if todos == nil {
		todos = []*Todo{}
	}
	return todos
}

func Add(description string) {
	if todos == nil {
		todos = []*Todo{}
	}
	if len(description) == 0 {
		return
	}
	todos = append(todos, &Todo{
		Id:          uuid.New(),
		Description: description,
	})
}

func Remove(id string) {
	if len(id) == 0 {
		return
	}
	if todos == nil {
		todos = []*Todo{}
	}
	index := -1
	for i, t := range todos {
		if t.Id.String() == id {
			index = i
			break
		}
	}
	if index < 0 {
		return
	}
	copy(todos[index:], todos[index+1:])
	todos[len(todos)-1] = nil
	todos = todos[:len(todos)-1]
}
