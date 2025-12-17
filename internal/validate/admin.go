package validate

func AdminPermission(uid int64, creatorID int64) bool {
	return uid == creatorID
}
