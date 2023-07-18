package database

import (
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/types"
	"golang.org/x/net/context"
	"time"
)

type contactsDB struct {
	db *sql.DB
}

func NewContactsDB(db *sql.DB) *contactsDB {
	return &contactsDB{
		db: db,
	}
}

func (db *contactsDB) WriteContact(ctx context.Context, fromID int64, contact *types.Contact) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"WriteContact",
	)
	defer span.Finish()

	const query = `
		INSERT INTO contacts(
			tg_user_id,
			contact_id,
			name,
			phone,
		    birthday,
		    description
		) values (
			$1, $2, $3, $4, $5, $6
		);
	`

	_, err := db.db.ExecContext(ctx, query,
		fromID,
		contact.ContactID,
		contact.Name,
		contact.Phone,
		contact.Birthday,
		contact.Description,
	)
	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *contactsDB) GetContact(ctx context.Context, userID int64, contactID int) (*types.Contact, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"GetContact",
	)
	defer span.Finish()

	const query = `
		SELECT 
			name,
			phone,
			birthday,
			description
		FROM contacts
		WHERE 
			tg_user_id = $1 AND contact_id = $2
	`

	contact := types.NewContact()

	err := db.db.QueryRowContext(ctx, query,
		userID,
		contactID,
	).Scan(&contact.Name, &contact.Phone, &contact.Birthday, &contact.Description)
	contact.ContactID = contactID

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrap(err, "cannot Scan")
		}

		// In another case we haven't just found the needed row.
		return nil, nil
	}

	return contact, nil
}

func (db *contactsDB) GetAllContacts(ctx context.Context, userID int64) ([]*types.Contact, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"GetAllContacts",
	)
	defer span.Finish()

	const query = `
		SELECT 
			contact_id,
			name,
			phone,
			birthday,
			description
		FROM contacts
		WHERE 
			tg_user_id = $1
	`

	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot QueryContext")
	}
	defer rows.Close()

	contacts := []*types.Contact{}
	for rows.Next() {
		contact := types.NewContact()
		err := rows.Scan(&contact.ContactID, &contact.Name, &contact.Phone, &contact.Birthday, &contact.Description)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan contact row")
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "cannot Scan")
	}

	return contacts, nil
}

func (db *contactsDB) GetContactByName(ctx context.Context, userID int64, name string) ([]*types.Contact, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"GetContactByName",
	)
	defer span.Finish()

	const query = `
		SELECT 
			contact_id,
			phone,
			birthday,
			description
		FROM contacts
		WHERE 
			tg_user_id = $1 AND name = $2
	`

	rows, err := db.db.QueryContext(ctx, query, userID, name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot QueryContext")
	}
	defer rows.Close()

	contacts := []*types.Contact{}
	for rows.Next() {
		contact := types.NewContact()
		err := rows.Scan(&contact.ContactID, &contact.Phone, &contact.Birthday, &contact.Description)
		if err != nil {
			return nil, errors.Wrap(err, "cannot Scan")
		}

		contact.Name = name
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (db *contactsDB) DeleteContact(ctx context.Context, userID int64, contactID int) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"DeleteContact",
	)
	defer span.Finish()

	const query = `
		DELETE FROM
			contacts
		WHERE
			tg_user_id = $1 AND
			contact_id = $2
	`
	_, err := db.db.ExecContext(ctx, query,
		userID,
		contactID,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *contactsDB) WriteName(ctx context.Context, name string, userID int64, contactID int) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"WriteName",
	)
	defer span.Finish()

	const query = `
		UPDATE 
			contacts
		SET
			name = $1
		WHERE
			tg_user_id = $2 AND
			contact_id = $3
	`

	_, err := db.db.ExecContext(ctx, query,
		name,
		userID,
		contactID,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *contactsDB) WritePhone(ctx context.Context, phone string, userID int64, contactID int) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"WritePhone",
	)
	defer span.Finish()

	const query = `
		UPDATE 
			contacts
		SET
			phone = $1
		WHERE
			tg_user_id = $2 AND
			contact_id = $3
	`

	_, err := db.db.ExecContext(ctx, query,
		phone,
		userID,
		contactID,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *contactsDB) WriteBirthday(ctx context.Context, birthday time.Time, userID int64, contactID int) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"WriteBirthday",
	)
	defer span.Finish()

	const query = `
		UPDATE 
			contacts
		SET
			birthday = $1
		WHERE
			tg_user_id = $2 AND
			contact_id = $3
	`

	_, err := db.db.ExecContext(ctx, query,
		birthday,
		userID,
		contactID,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}

func (db *contactsDB) WriteDescription(ctx context.Context, description string, userID int64, contactID int) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"WriteDescription",
	)
	defer span.Finish()

	const query = `
		UPDATE 
			contacts
		SET
			description = $1
		WHERE
			tg_user_id = $2 AND
			contact_id = $3
	`

	_, err := db.db.ExecContext(ctx, query,
		description,
		userID,
		contactID,
	)

	if err != nil {
		return errors.Wrap(err, "cannot ExecContent")
	}

	return nil
}
