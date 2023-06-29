package types

// State tells what we are modifying un the expense.
type State int

const (
	EditingName State = iota + 1
	EditingPhone
	EditingBirthday
	EditingDescription
	EditingSearchPhrase
	EditingEditID
	WaitState
)

// CurrentState contains id on the expense we are modifying now, and what we are modifying.
type CurrentState struct {
	ContactID int
	MessageID int
	State     State
}
