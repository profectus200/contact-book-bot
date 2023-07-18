package database

import (
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/types"
)

type usersDB struct {
	db *sql.DB
}

func NewUsersDB(db *sql.DB) *usersDB {
	return &usersDB{
		db: db,
	}
}

// Actions with users table.

func (db *usersDB) SetCurrentState(ctx context.Context, userID int64, state types.CurrentState) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"SetCurrentState",
	)
	defer span.Finish()

	const query = `
		INSERT INTO users(
			tg_user_id,
			contact_id,
		    message_id,
			current_state
		) VALUES (
			$1, $2, $3, $4
		)
		ON CONFLICT (tg_user_id) DO UPDATE
		SET 
			contact_id = $2,
			message_id = $3,
			current_state = $4
	`

	_, err := db.db.ExecContext(ctx, query,
		userID,
		state.ContactID,
		state.MessageID,
		state.State,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *usersDB) GetCurrentState(ctx context.Context, userID int64) (*types.UserStateType, bool) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"GetCurrentState",
	)
	defer span.Finish()

	const query = `
		SELECT
			contact_id,
			message_id,
			current_state
		FROM
			users
		WHERE
			tg_user_id = $1
	`

	var userState types.UserStateType

	err := db.db.QueryRowContext(ctx, query,
		userID,
	).Scan(&userState.CurrentState.ContactID, &userState.CurrentState.MessageID, &userState.CurrentState.State)

	if err != nil {
		return nil, false
	}

	return &userState, true
}

func (db *usersDB) ToWaitState(ctx context.Context, userID int64) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"ToWaitState",
	)
	defer span.Finish()

	const query = `
		INSERT INTO users(
			tg_user_id,
			current_state
		) VALUES (
			$1, $2
		)
		ON CONFLICT(tg_user_id)
		DO UPDATE 
			SET
			current_state = $2
	`

	_, err := db.db.ExecContext(ctx, query,
		userID,
		types.WaitState,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}
