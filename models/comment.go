package models

type Comment struct {
	ID     interface{} // data persistence backend defines the type
	Text   string      `validate:"required,lte=5000"`
	TaskID interface{} // data persistence backend defines the type
}

func (c Comment) Validate() error {
	return validate.Struct(c)
}
