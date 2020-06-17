package services

import (
	"context"

	"github.com/ortymid/projector/models"
)

type TaskRepo interface {
	// All(context.Context) []Task
	AllByColumn(context.Context, models.Column) ([]models.Task, error)
	AllByBoard(context.Context, models.Board) ([]models.Task, error)
	Create(context.Context, models.Column, models.Task) (models.Task, error)
}

type TaskService struct {
	tr TaskRepo
}

func NewTaskService(tr TaskRepo) TaskService {
	return TaskService{tr: tr}
}

func (srv TaskService) CreateTask(c models.Column, name, desc string) (t models.Task, err error) {
	t.Name = name
	t.Description = desc

	if err = t.Validate(); err != nil {
		return t, err
	}

	return srv.tr.Create(context.TODO(), c, t)
}
