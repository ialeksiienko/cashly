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

type CustomError[T any] struct {
	Data T
	Msg  string
	Code ErrorCode
}

func (e *CustomError[T]) Error() string {
	return e.Msg
}

func New[T any](msg string, code ErrorCode, data T) *CustomError[T] {
	return &CustomError[T]{Msg: msg, Code: code, Data: data}
}
