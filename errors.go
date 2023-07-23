package privnote

import "errors"

var (
	// ErrInvalidNote is returned when the note is invalid.
	// This is usually because the note has been destroyed or
	// does not exist.
	ErrInvalidNote = errors.New("note is invalid")

	// ErrEmptyNoteID is returned when the note ID is empty.
	ErrEmptyNoteID = errors.New("note ID cannot be empty")

	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = errors.New("note password is invalid")
)
