package fasttar

// NameString is a zero-allocation string(Name(header))
func NameString(header []byte) string {
	return unsafeString(Name(header))
}

// LinkNameString is a zero-allocation string(LinkName(header))
func LinkNameString(header []byte) string {
	return unsafeString(LinkName(header))
}

// UserNameString is a zero-allocation string(UserName(header))
func UserNameString(header []byte) string {
	return unsafeString(UserName(header))
}

// GroupNameString is a zero-allocation string(GroupName(header))
func GroupNameString(header []byte) string {
	return unsafeString(GroupName(header))
}
