package models

type User struct {
	ID   interface{} // data persistence backend defines the type
	Name string      `validate:"required"`
}

func (u User) Validate() error {
	return validate.Struct(u)
}
