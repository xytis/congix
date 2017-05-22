package helper

// BoolToPtr returns the pointer to a boolean
func BoolToPtr(b bool) *bool {
	return &b
}

// StringToPtr returns the pointer to a string
func StringToPtr(str string) *string {
	return &str
}
