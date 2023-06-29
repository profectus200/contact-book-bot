package callbacks

import (
	"context"
	"errors"
	"github.com/profectus200/contact-book-bot/internal/types"
)

const (
	// editContactKeyboard
	ChangeContactName        string = "ChangeContactName"
	ChangeContactPhone       string = "ChangeContactPhone"
	ChangeContactBirthday    string = "ChangeContactBirthday"
	ChangeContactDescription string = "ChangeContactDescription"
	ChangeContactDone        string = "ChangeContactDone"
	DeleteContact            string = "DeleteContact"
)

type callbackHandler interface {
	SendMessage(text string, userID int64) error
	DoneMessage(userID int64, messageID int) error
	DeleteMessage(userID int64, messageID int) error
	ShowAlert(text string, messageID string) error
}

type contactsDB interface {
	DeleteContact(ctx context.Context, userID int64, contactID int) error
}

type usersDB interface {
	SetCurrentState(ctx context.Context, userID int64, state types.CurrentState) error
	GetCurrentState(ctx context.Context, userID int64) (*types.UserStateType, bool)
	ToWaitState(ctx context.Context, userID int64) error
}

type Model struct {
	tgClient   callbackHandler
	contactsDB contactsDB
	usersDB    usersDB
}

func New(tgClient callbackHandler, contactsDB contactsDB, usersDB usersDB) *Model {
	return &Model{
		tgClient:   tgClient,
		contactsDB: contactsDB,
		usersDB:    usersDB,
	}
}

type CallbackData struct {
	FromID     int64
	MessageID  int
	Data       string
	CallbackID string
}

func (s *Model) IncomingCallback(ctx context.Context, data *CallbackData) error {
	switch data.Data {
	case ChangeContactName:
		return s.toWriteNameState(ctx, data)
	case ChangeContactPhone:
		return s.toWritePhoneState(ctx, data)
	case ChangeContactBirthday:
		return s.toWriteBirthdayState(ctx, data)
	case ChangeContactDescription:
		return s.toWriteDescriptionState(ctx, data)
	case ChangeContactDone:
		return s.saveContact(data)
	case DeleteContact:
		return s.deleteContact(ctx, data)
	}

	return errors.New("Callback handler for data '" + data.Data + "' was not found.")
}
