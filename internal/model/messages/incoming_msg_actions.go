package messages

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/types"
)

func (s *Model) editContactAfterEditing(ctx context.Context, contact *types.Contact, userID int64, messageID int) error {
	message := contact.ToString()
	return s.tgClient.EditContactMessage(message, userID, messageID)
}

func (s *Model) nameEntered(ctx context.Context, msg *Message, userState types.CurrentState) error {
	name := msg.Text

	contact, err := s.contactsDB.GetContact(ctx, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot GetContact")
	}

	if contact == nil {
		contact = types.NewContact()
		contact.ContactID = userState.ContactID
		err := s.initializeContact(ctx, msg, contact)
		if err != nil {
			return errors.Wrap(err, "cannot InitializeContact")
		}
	}

	contact.Name = name
	err = s.contactsDB.WriteName(ctx, name, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot WriteName")
	}

	err = s.tgClient.DeleteMessage(msg.UserID, msg.MessageID)
	if err != nil {
		return errors.Wrap(err, "cannot DeleteMessage")
	}

	err = s.usersDB.ToWaitState(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot ToWaitState")
	}

	return s.editContactAfterEditing(ctx, contact, msg.UserID, userState.MessageID)
}

func (s *Model) phoneEntered(ctx context.Context, msg *Message, userState types.CurrentState) error {
	phone := msg.Text

	contact, err := s.contactsDB.GetContact(ctx, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot GetContact")
	}

	if contact == nil {
		contact = types.NewContact()
		contact.ContactID = userState.ContactID
		err := s.initializeContact(ctx, msg, contact)
		if err != nil {
			return errors.Wrap(err, "cannot InitializeContact")
		}
	}

	contact.Phone = phone
	err = s.contactsDB.WritePhone(ctx, phone, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot WritePhone")
	}

	err = s.tgClient.DeleteMessage(msg.UserID, msg.MessageID)
	if err != nil {
		return errors.Wrap(err, "cannot DeleteMessage")
	}

	err = s.usersDB.ToWaitState(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot ToWaitState")
	}

	return s.editContactAfterEditing(ctx, contact, msg.UserID, userState.MessageID)
}

func (s *Model) birthdayEntered(ctx context.Context, msg *Message, userState types.CurrentState) error {
	year := time.Now().Year() - 1000
	dateString := fmt.Sprintf("%s.%d", msg.Text, year)
	birthday, err := time.Parse("02.01.2006", dateString)
	if err != nil {
		return errors.Wrap(err, "Cannot time.Parse")
	}

	contact, err := s.contactsDB.GetContact(ctx, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot GetContact")
	}

	if contact == nil {
		contact = types.NewContact()
		contact.ContactID = userState.ContactID
		err := s.initializeContact(ctx, msg, contact)
		if err != nil {
			return errors.Wrap(err, "cannot InitializeContact")
		}
	}

	contact.Birthday = birthday
	err = s.contactsDB.WriteBirthday(ctx, birthday, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot WriteBirthday")
	}

	err = s.tgClient.DeleteMessage(msg.UserID, msg.MessageID)
	if err != nil {
		return errors.Wrap(err, "cannot DeleteMessage")
	}

	err = s.usersDB.ToWaitState(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot ToWaitState")
	}

	return s.editContactAfterEditing(ctx, contact, msg.UserID, userState.MessageID)
}

func (s *Model) descriptionEntered(ctx context.Context, msg *Message, userState types.CurrentState) error {
	description := msg.Text

	contact, err := s.contactsDB.GetContact(ctx, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot GetContact")
	}

	if contact == nil {
		contact = types.NewContact()
		contact.ContactID = userState.ContactID
		err := s.initializeContact(ctx, msg, contact)
		if err != nil {
			return errors.Wrap(err, "cannot InitializeContact")
		}
	}

	contact.Description = description
	err = s.contactsDB.WriteDescription(ctx, description, msg.UserID, userState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot WriteDescription")
	}

	err = s.tgClient.DeleteMessage(msg.UserID, msg.MessageID)
	if err != nil {
		return errors.Wrap(err, "cannot DeleteMessage")
	}

	err = s.usersDB.ToWaitState(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot ToWaitState")
	}

	return s.editContactAfterEditing(ctx, contact, msg.UserID, userState.MessageID)
}

func (s *Model) initializeContact(ctx context.Context, msg *Message, contact *types.Contact) error {
	err := s.contactsDB.WriteContact(ctx, msg.UserID, contact)
	if err != nil {
		return errors.Wrap(err, "cannot WriteContact")
	}

	return nil
}

func (s *Model) addContact(ctx context.Context, userID int64) error {
	err := s.usersDB.SetCurrentState(ctx, userID, types.CurrentState{
		ContactID: 0,
		MessageID: 0,
		State:     types.WaitState,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	return s.tgClient.EditContact(types.NewContact().ToString(), userID)
}

func (s *Model) getContact(ctx context.Context, userID int64) error {
	err := s.tgClient.SendMessage(getContactMsg, userID)
	if err != nil {
		return errors.Wrap(err, "cannot SendMessage")
	}

	err = s.usersDB.SetCurrentState(ctx, userID, types.CurrentState{
		ContactID: -1,
		State:     types.EditingSearchPhrase,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	return nil
}

func (s *Model) searchPhraseEntered(ctx context.Context, msg *Message) error {
	searchPhrase := msg.Text

	contacts, err := s.contactsDB.GetContactByName(ctx, msg.UserID, searchPhrase)
	if err != nil {
		return errors.Wrap(err, "cannot GetContactByName")
	}

	err = s.usersDB.ToWaitState(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot ToWaitState")
	}

	text := ""
	for _, contact := range contacts {
		text += contact.ToString() + "-----------------------------\n"
	}

	return s.tgClient.SendMessage(text, msg.UserID)
}

func (s *Model) editContact(ctx context.Context, userID int64) error {
	err := s.tgClient.SendMessage(editContactMsg, userID)
	if err != nil {
		return errors.Wrap(err, "cannot SendMessage")
	}

	err = s.usersDB.SetCurrentState(ctx, userID, types.CurrentState{
		ContactID: -1,
		State:     types.EditingEditID,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	return nil
}

func (s *Model) editIDEntered(ctx context.Context, msg *Message) error {
	ID, err := strconv.Atoi(msg.Text)
	if err != nil {
		return errors.Wrap(err, "cannot convert string to int")
	}

	contact, err := s.contactsDB.GetContact(ctx, msg.UserID, ID)
	if err != nil {
		return errors.Wrap(err, "cannot GetContact")
	}

	err = s.usersDB.SetCurrentState(ctx, msg.UserID, types.CurrentState{
		ContactID: contact.ContactID,
		State:     types.WaitState,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	return s.tgClient.EditContact(contact.ToString(), msg.UserID)
}

func (s *Model) listContacts(ctx context.Context, userID int64) error {
	contacts, err := s.contactsDB.GetAllContacts(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "cannot GetAllContacts")
	}

	text := ""
	if len(contacts) == 0 {
		text = "You don't have any contacts saved yet!"
	} else {
		for _, contact := range contacts {
			text += contact.ToString() + "-----------------------------\n"
		}
	}

	return s.tgClient.SendMessage(text, userID)
}
