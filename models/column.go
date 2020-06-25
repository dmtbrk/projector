package models

type Column struct {
	ID         interface{} // data persistence backend defines the type
	BoardID    interface{} // data persistence backend defines the type
	Name       string      `validate:"required,lte=255"`
	OrderIndex int         // user defined order, on column delete tasks go to the column with lower index
}

func NewColumn(name string, orderIndex int) (c Column, err error) {
	c.Name = name

	err = c.Validate()

	return c, err
}

func (c Column) Validate() error {
	return validate.Struct(c)
}
