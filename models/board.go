package models

type Board struct {
	ID          interface{} // data persistence backend defines the type
	Name        string      `validate:"required,lte=500"`
	Description string      `validate:"lte=1000"`
	UserID      interface{} // data persistence backend defines the type
}

func (b Board) Validate() error {
	return validate.Struct(b)
}
