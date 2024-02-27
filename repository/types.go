// This file contains types that are used in the repository layer.
package repository

type CreateUserInput struct {
	ID             string `db:"id"`
	FullName       string `db:"full_name"`
	PhoneNumber    string `db:"phone_number"`
	HashedPassword string `db:"hashed_password"`
}

type UserOutput struct {
	ID             string `db:"id"`
	FullName       string `db:"full_name"`
	PhoneNumber    string `db:"phone_number"`
	LoginCount     int    `db:"login_count"`
	HashedPassword string `db:"hashed_password"`
}

type UpdateUserInput struct {
	FullName    string `db:"full_name"`
	PhoneNumber string `db:"phone_number"`
}

type CreateUserOutput struct {
	ID string `db:"id"`
}
