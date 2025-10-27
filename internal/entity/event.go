package entity

type EventNotification struct {
	Event           string
	CheckedUserID   int64
	CheckedByUserID int64
	FamilyName      string
}
