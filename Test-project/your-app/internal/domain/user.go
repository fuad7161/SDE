// # Business models & rules (NO external imports)

package domain

import "errors"

type User struct {
	ID   string
	Name string
	Age  int
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.Age < 0 {
		return errors.New("age cannot be negative")
	}

	return nil
}
