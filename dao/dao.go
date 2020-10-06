package dao

import (
	"errors"
	"fmt"
	model "gotodo/proto"

	"github.com/sirupsen/logrus"
)

// To generate MockDao run the command
// mockery -name Dao -inpkg -case snake .

// Consider adding the following function to bash_profile
// function mock () {
//    mockery -name $1 -inpkg -case snake .
// }
// Then you can run the command
// mock Dao
// for future struct mocks

type Dao interface {
	Create(todo *model.Task) (*model.Task, error)
	Delete(id string) error
	Get() ([]*model.Task, error)
}

type todoDao struct {
	logger *logrus.Entry
	todos  []*model.Task
	nextId int
}

func NewDao(logger *logrus.Entry) Dao {
	return &todoDao{
		logger: logger.WithField("struct", "TodoDao"),
		nextId: 0,
	}
}

func (d *todoDao) Create(todo *model.Task) (*model.Task, error) {
	logger := d.logger.WithField("todo", todo)
	logger.Debug("Create")

	if todo == nil {
		msg := "todo cannot be nil"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	todo.Id = fmt.Sprintf("%d", d.nextId)
	d.nextId++

	d.todos = append(d.todos, todo)

	logger.WithField("result", todo).Trace("Stored todo")
	return todo, nil
}

func (d *todoDao) Delete(id string) error {
	logger := d.logger.WithField("todoID", id)
	logger.Debug("Delete")

	for i, t := range d.todos {
		if t.Id == id {
			d.todos = append(d.todos[:i], d.todos[i+1:]...)
		}
	}
	return nil
}

func (d *todoDao) Get() ([]*model.Task, error) {
	d.logger.Trace("Get")
	return d.todos, nil
}
