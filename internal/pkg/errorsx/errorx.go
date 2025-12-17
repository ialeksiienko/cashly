package errorsx

type ErrorCode string

const (
	ErrCodeFamilyCodeExpired ErrorCode = "family_code_expired"

	ErrCodeFamilyNotFound ErrorCode = "family_not_found"
	ErrCodeTokenNotFound  ErrorCode = "token_not_found"

	ErrCodeUserHasNoFamily    ErrorCode = "user_has_no_family"
	ErrCodeFamilyHasNoMembers ErrorCode = "family_has_no_members"
	ErrCodeUserNotInFamily    ErrorCode = "user_not_in_family"

	ErrCodeNoPermission     ErrorCode = "no_permission"
	ErrCodeCannotRemoveSelf ErrorCode = "cannot_remove_self"

	ErrCodeFailedToGenerateInviteCode ErrorCode = "failed_to_generate_invite_code"

	ErrRequestCooldown ErrorCode = "api_request_cooldown"
)

type ErrorWithCode interface {
	error
	GetCode() ErrorCode
	GetData() any
}

func (e *Error[T]) GetCode() ErrorCode {
	return e.Code
}

func (e *Error[T]) GetData() any {
	return e.Data
}

type Error[T any] struct {
	Data T
	Msg  string
	Code ErrorCode
}

func (e *Error[T]) Error() string {
	return e.Msg
}

func New[T any](msg string, code ErrorCode, data T) *Error[T] {
	return &Error[T]{Msg: msg, Code: code, Data: data}
}
