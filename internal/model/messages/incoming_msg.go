package messages

import (
	"context"
	"time"

	"github.com/profectus200/contact-book-bot/internal/types"
)

type messageSender interface {
	SendMessage(text string, userID int64) error
	EditContact(text string, userID int64) error
	EditContactMessage(text string, userID int64, messageID int) error
	DeleteMessage(userID int64, messageID int) error
}

type contactsDB interface {
	WriteContact(ctx context.Context, userID int64, contact *types.Contact) error
	GetContact(ctx context.Context, userID int64, contactID int) (*types.Contact, error)
	GetContactByName(ctx context.Context, userID int64, name string) ([]*types.Contact, error)
	GetAllContacts(ctx context.Context, userID int64) ([]*types.Contact, error)
	WriteName(ctx context.Context, name string, userID int64, contactID int) error
	WritePhone(ctx context.Context, phone string, userID int64, contactID int) error
	WriteBirthday(ctx context.Context, birthday time.Time, userID int64, contactID int) error
	WriteDescription(ctx context.Context, description string, userID int64, contactID int) error
}

type usersDB interface {
	ToWaitState(ctx context.Context, userID int64) error
	SetCurrentState(ctx context.Context, userID int64, state types.CurrentState) error
	GetCurrentState(ctx context.Context, userID int64) (*types.UserStateType, bool)
}

type Model struct {
	tgClient   messageSender
	contactsDB contactsDB
	usersDB    usersDB
}

func New(tgClient messageSender, contactsDB contactsDB, usersDB usersDB) *Model {
	return &Model{
		tgClient:   tgClient,
		contactsDB: contactsDB,
		usersDB:    usersDB,
	}
}

type Message struct {
	Text      string
	UserID    int64
	MessageID int
}

const (
	getContactMsg  = "Write the name of your contact:"
	editContactMsg = "Write ID of the contact you want to edit:"
)

func (s *Model) IncomingMessage(ctx context.Context, msg *Message) error {
	// Trying to recognize the command.
	switch msg.Text {
	case "/start":
		return s.tgClient.SendMessage("Hello! You can save people contacts here!:)", msg.UserID)
	case "/add_contact":
		return s.addContact(ctx, msg.UserID)
	case "/get_contact":
		return s.getContact(ctx, msg.UserID)
	case "/edit_contact":
		return s.editContact(ctx, msg.UserID)
	case "/list_contacts":
		return s.listContacts(ctx, msg.UserID)
	}

	// It is not a known command - maybe it is message to change the state.
	if userState, ok := s.usersDB.GetCurrentState(ctx, msg.UserID); ok {
		switch userState.CurrentState.State {
		case types.EditingName:
			return s.nameEntered(ctx, msg, userState.CurrentState)
		case types.EditingPhone:
			return s.phoneEntered(ctx, msg, userState.CurrentState)
		case types.EditingBirthday:
			return s.birthdayEntered(ctx, msg, userState.CurrentState)
		case types.EditingDescription:
			return s.descriptionEntered(ctx, msg, userState.CurrentState)
		case types.EditingSearchPhrase:
			return s.searchPhraseEntered(ctx, msg)
		case types.EditingEditID:
			return s.editIDEntered(ctx, msg)
		case types.WaitState:
		}
	}

	return s.tgClient.SendMessage("I do not know such a command", msg.UserID)
}
