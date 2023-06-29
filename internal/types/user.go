package types

type CurrentContact struct {
	ContactID int
	Contact   Contact
}

type UserStateType struct {
	CurrentState CurrentState
}
