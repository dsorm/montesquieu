package users

// User is a data structure that should represent its structure in Store.
// Anyone that can login is a user.
type User struct {
	// The unique identifier for each user, used for relation with other
	// structures
	ID uint64

	// A real name of the user, for example "John Doe"
	DisplayName string

	// example: "john_doe"
	Login string

	// User's password for login. Usually left empty, unless changing password
	Password string
}

// Author is a kind of User which can publish Articles
type Author struct {
	User

	// Authors have their own ID's, separate from User ID's
	AuthorID uint64
}

// Admin is a kind of User that has access to the Admin Panel
type Admin User

// HashPassword hashes password from user's input, so it can be safely put
// into a database
func (u *User) HashPassword(input string) {
	// TODO
}

// VerifyPassword verifies the password that has user provided ('input') against
// the hash from a database ('hash')
func (u *User) VerifyPassword(hash string, input string) {
	// TODO
}
