package types

type Student struct {
	ID string 
	Name string `validate:"required"`
	Email string `validate:"required"`
	Age int `validate:"required"`

}