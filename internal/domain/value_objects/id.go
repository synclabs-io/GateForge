package value_objects

import "github.com/google/uuid"

type ID struct {
	value string
}

func NewID() *ID {
	return &ID{value: uuid.NewString()}
}

func (i *ID) String() string {
	return i.value
}
