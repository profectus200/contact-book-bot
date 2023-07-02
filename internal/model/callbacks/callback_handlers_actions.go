package callbacks

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/types"
)

func (s *Model) toWriteNameState(ctx context.Context, data *CallbackData) error {
	// Change state of the user - he is now entering name for this user and this messageID.
	contactID := data.MessageID
	if state, ok := s.usersDB.GetCurrentState(ctx, data.FromID); ok {
		if state.CurrentState.ContactID != 0 {
			contactID = state.CurrentState.ContactID
		}
	}

	err := s.usersDB.SetCurrentState(ctx, data.FromID, types.CurrentState{
		ContactID: contactID,
		MessageID: data.MessageID,
		State:     types.EditingName,
	})

	if state, ok := s.usersDB.GetCurrentState(ctx, data.FromID); ok {
		fmt.Println(state.CurrentState.State, state.CurrentState.ContactID, state.CurrentState.MessageID)
	}
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	// Show notification about the action.
	return s.tgClient.ShowAlert("Enter the name of your contact:", data.CallbackID)
}

func (s *Model) toWritePhoneState(ctx context.Context, data *CallbackData) error {
	// Change state of the user - he is now entering name for this user and this messageID.

	contactID := data.MessageID
	if state, ok := s.usersDB.GetCurrentState(ctx, data.FromID); ok {
		if state.CurrentState.ContactID != 0 {
			contactID = state.CurrentState.ContactID
		}
	}

	err := s.usersDB.SetCurrentState(ctx, data.FromID, types.CurrentState{
		ContactID: contactID,
		MessageID: data.MessageID,
		State:     types.EditingPhone,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	// Show notification about the action.
	return s.tgClient.ShowAlert("Enter the phone of your contact:", data.CallbackID)
}

func (s *Model) toWriteBirthdayState(ctx context.Context, data *CallbackData) error {
	// Change state of the user - he is now entering name for this user and this messageID.
	contactID := data.MessageID
	if state, ok := s.usersDB.GetCurrentState(ctx, data.FromID); ok {
		if state.CurrentState.ContactID != 0 {
			contactID = state.CurrentState.ContactID
		}
	}

	err := s.usersDB.SetCurrentState(ctx, data.FromID, types.CurrentState{
		ContactID: contactID,
		MessageID: data.MessageID,
		State:     types.EditingBirthday,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	// Show notification about the action.
	return s.tgClient.ShowAlert("Enter the birthday of your contact in format 'dd.mm':", data.CallbackID)
}

func (s *Model) toWriteDescriptionState(ctx context.Context, data *CallbackData) error {
	// Change state of the user - he is now entering name for this user and this messageID.
	contactID := data.MessageID
	if state, ok := s.usersDB.GetCurrentState(ctx, data.FromID); ok {
		if state.CurrentState.ContactID != 0 {
			contactID = state.CurrentState.ContactID
		}
	}

	err := s.usersDB.SetCurrentState(ctx, data.FromID, types.CurrentState{
		ContactID: contactID,
		MessageID: data.MessageID,
		State:     types.EditingDescription,
	})
	if err != nil {
		return errors.Wrap(err, "cannot SetCurrentState")
	}

	// Show notification about the action.
	return s.tgClient.ShowAlert("Enter description of your contact:", data.CallbackID)
}

func (s *Model) saveContact(data *CallbackData) error {
	err := s.tgClient.ShowAlert("Saved", data.CallbackID)
	if err != nil {
		return errors.Wrap(err, "cannot ShowAlert")
	}

	return s.tgClient.DeleteMessage(data.FromID, data.MessageID)

}

func (s *Model) deleteContact(ctx context.Context, data *CallbackData) error {
	state, _ := s.usersDB.GetCurrentState(ctx, data.FromID)
	err := s.contactsDB.DeleteContact(ctx, data.FromID, state.CurrentState.ContactID)
	if err != nil {
		return errors.Wrap(err, "cannot DeleteContact")
	}

	err = s.tgClient.ShowAlert("Deleted", data.CallbackID)
	if err != nil {
		return errors.Wrap(err, "cannot ShowAlert")
	}

	return s.tgClient.DeleteMessage(data.FromID, data.MessageID)
}
