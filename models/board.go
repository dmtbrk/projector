package models

type Board struct {
	ID          interface{} // data persistence backend defines the type
	Name        string      `validate:"required,lte=500"`
	Description string      `validate:"lte=1000"`
	UserID      interface{} // data persistence backend defines the type
}

func NewBoard(name, desc string) (b Board, err error) {
	b.Name = name
	b.Description = desc

	err = b.Validate()

	return b, err
}

func (b Board) Validate() error {
	return validate.Struct(b)
}
