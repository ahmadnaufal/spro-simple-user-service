package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	query := `
		INSERT INTO
			users
			(id, full_name, phone_number, hashed_password)
		VALUES
			(:id, :full_name, :phone_number, :hashed_password)
		RETURNING id
	`

	bindQuery, args, err := sqlx.Named(query, input)
	if err != nil {
		return
	}

	err = r.Db.GetContext(ctx, &output, sqlx.Rebind(sqlx.DOLLAR, bindQuery), args...)

	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (output UserOutput, err error) {
	query := `
		SELECT
			id,
			full_name,
			phone_number,
			hashed_password,
			login_count
		FROM
			users
		WHERE
			phone_number = $1
		LIMIT 1
	`

	err = r.Db.GetContext(ctx, &output, query, phoneNumber)

	return
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (output UserOutput, err error) {
	query := `
		SELECT
			id,
			full_name,
			phone_number,
			hashed_password,
			login_count
		FROM
			users
		WHERE
			id = $1
		LIMIT 1
	`

	err = r.Db.GetContext(ctx, &output, query, id)

	return
}

func (r *Repository) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (err error) {
	query := `
		UPDATE
			users
		SET
			full_name = $1,
			phone_number = $2
		WHERE
			id = $3
	`

	res, err := r.Db.ExecContext(ctx, query, input.FullName, input.PhoneNumber, id)
	if err != nil {
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		err = fmt.Errorf("unexpected behavior: expected update 1 row but got %d", rowsAffected)
		return
	}

	return
}

func (r *Repository) IncrementLoginCount(ctx context.Context, id string) error {
	query := `
		UPDATE
			users
		SET
			login_count = login_count + 1
		WHERE
			id = $1
	`

	res, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		err = fmt.Errorf("unexpected behavior: expected update 1 row but got %d", rowsAffected)
		return err
	}

	return nil
}
