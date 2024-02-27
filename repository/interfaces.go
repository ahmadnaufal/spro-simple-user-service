// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateUser(context.Context, CreateUserInput) (CreateUserOutput, error)
	GetUserByPhoneNumber(context.Context, string) (UserOutput, error)
	GetUserByID(context.Context, string) (UserOutput, error)
	UpdateUser(context.Context, string, UpdateUserInput) error
	IncrementLoginCount(context.Context, string) error
}
