package service

import (
	d "gotodo/dao"
	m "gotodo/proto"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Create(sc *m.Task) (*m.Task, error)
	Delete(id string) error
	Get() ([]*m.Task, error)
}

type todoService struct {
	logger *logrus.Entry
	dao    d.Dao
}

func NewService(logger *logrus.Entry, dao d.Dao) Service {
	return &todoService{
		logger: logger.WithField("struct", "TodoService"),
		dao:    dao,
	}
}

func (s *todoService) Create(todo *m.Task) (*m.Task, error) {
	return s.dao.Create(todo)
}

func (s *todoService) Delete(id string) error {
	return s.dao.Delete(id)
}

func (s *todoService) Get() ([]*m.Task, error) {
	return s.dao.Get()
}
