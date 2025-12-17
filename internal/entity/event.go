package entity

type EventType string

const (
	EventBalanceChecked    EventType = "balance_checked"
	EventDeletedFromFamily EventType = "deleted_from_family"
	EventJoinedFamily      EventType = "joined_family"
	EventLeavedFromFamily  EventType = "leaved_from_family"
)

type EventNotification struct {
	Type        EventType
	RecipientID int64
	FamilyName  string
	Data        map[string]any
}
