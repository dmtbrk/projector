package models

type Task struct {
	ID          interface{}
	Name        string `validate:"required,lte=500"`
	Description string `validate:"lte=5000"`
	ColumnID    interface{}
}

func (t Task) Validate() error {
	return validate.Struct(t)
}
