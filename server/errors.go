package server

// All the different errors
const (
	ErrPasswordMissing     string = "Password Not Present"
	ErrNotRegistered       string = "No records found"
	ErrInternalError       string = "An Error Occured"
	ErrPasswordMatchFailed string = "Passwords do not match"
	ErrURLMissing          string = "URL Missing"
	ErrPasswordTooShort    string = "Password too short"
	ErrURLNotPresent       string = "URL not present"
)
