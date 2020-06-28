package services

import (
	"context"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

type TaskService struct {
	tr repos.TaskRepo
}

func NewTaskService(tr repos.TaskRepo) TaskService {
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
