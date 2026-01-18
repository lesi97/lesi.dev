package utils

/*
Utility function to truncate a string if the length of the string is greater than the specified character limit

It modifies the original string by adding "..." at the end if truncation occurs
The final length will not exceed the given limit

Example:
func init() {
	message := "This is a test message"
	fmt.Println(message) 		// Output: This is a test message
	TruncateString(&message, 17)	// The address of the message is passed through, not the message directly
	fmt.Println(message) 		// Output: This is a test...
}
*/
func TruncateString(message *string, charLimit int) {
	clone := *message
	if len(*message) > (charLimit - 3) {
		*message = clone[:(charLimit-3)] + "..."
	}
}
