package privnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/LightningDev1/go-privnote/util"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

// Client is a Privnote client for creating and opening notes.
type Client struct {
	// BaseURL is the base URL for the Privnote API.
	// Defaults to "https://privnote.com".
	BaseURL string

	// HTTPClient is the client used to make HTTP requests to the Privnote API.
	HTTPClient tls_client.HttpClient
}

// CreateNoteData is the data required for creating a note.
type CreateNoteData struct {
	// Data is the content of the note.
	Data string

	// Password is the password for the note. Will be generated if empty.
	Password string

	// Duration is the time until the note expires. Defaults to 24 hours.
	// Maximum duration is 30 days or 720 hours.
	Duration time.Duration

	// NoConfirmation determines whether the user should be asked for confirmation
	// before destroying the note.
	NoConfirmation bool

	// NotifyEmail is the email address to notify when the note is read.
	NotifyEmail string

	// EmailRefName is a reference name of the note to include in the email.
	EmailRefName string
}

// CreateNote creates a new note and returns the URL.
func (c *Client) CreateNote(noteData CreateNoteData) (string, error) {
	// Encrypt the note data.
	encryptedData, password, err := util.Encrypt(noteData.Data, noteData.Password)
	if err != nil {
		return "", err
	}

	data := url.Values{
		// Data is the encrypted content of the note.
		"data": {encryptedData},

		// DataType is the type of data to be stored in the note.
		"data_type": {"T"},

		// HasManualPass is whether the note has a manual password.
		"has_manual_pass": {fmt.Sprint(noteData.Password != "")},

		// DurationHours is the time until the note expires in hours.
		"duration_hours": {fmt.Sprint(noteData.Duration.Hours())},

		// DontAsk is whether the user should be asked for confirmation before
		// destroying the note.
		"dont_ask": {fmt.Sprint(noteData.NoConfirmation)},

		// NotifyEmail is the email address to notify when the note is read.
		"notify_email": {noteData.NotifyEmail},

		// NotifyRef is a reference name of the note to include in the email.
		"notify_ref": {noteData.EmailRefName},
	}

	reader := strings.NewReader(data.Encode())

	// Create the note.
	req, err := http.NewRequest("POST", c.BaseURL+"/legacy/", reader)
	if err != nil {
		return "", err
	}

	req.Header = Headers.Clone()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	// Decode the JSON response.
	var note Note

	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		return "", err
	}

	if note.HasManualPass {
		return note.NoteLink, nil
	}

	return note.NoteLink + "#" + password, nil
}

// ReadNoteFromLink opens and decrypts a note from its link.
func (c *Client) ReadNoteFromLink(noteLink string) (string, error) {
	// Validate the note link.
	// The note link must include a # to separate the note ID and password.
	if noteLink == "" || !strings.Contains(noteLink, "#") {
		return "", ErrInvalidNote
	}

	// Fetch the note data.
	req, err := http.NewRequest("DELETE", noteLink, nil)
	if err != nil {
		return "", err
	}

	req.Header = Headers.Clone()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	// Decode the JSON response.
	var note Note

	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		// The response is not valid JSON.
		// This is usually because the note has been destroyed or does not exist,
		// so return ErrInvalidNote.
		return "", errors.Join(err, ErrInvalidNote)
	}

	// Decrypt the note data.
	password := strings.Split(noteLink, "#")[1]

	plainText, err := util.Decrypt(note.Data, password)
	if err != nil {
		// The ciphertext is invalid. This is usually caused by the note data
		// being invalid, so return ErrInvalidNote.
		if errors.Is(err, util.ErrInvalidCipherText) {
			return "", errors.Join(err, ErrInvalidNote)
		}

		// The padding is invalid. This is usually caused by decrypting the note
		// with an incorrect password, which causes the plaintext to be invalid,
		// thus having invalid padding.
		if errors.Is(err, util.ErrInvalidPadding) {
			return "", errors.Join(err, ErrInvalidPassword)
		}

		return "", err
	}

	return plainText, nil

}

// ReadNoteFromID opens and decrypts a note from its ID and password.
func (c *Client) ReadNoteFromID(noteID, password string) (string, error) {
	// Validate the note ID and password.
	if noteID == "" {
		return "", ErrEmptyNoteID
	}

	if password == "" {
		return "", ErrInvalidPassword
	}

	// Construct the note link and pass it to ReadNoteFromLink.
	return c.ReadNoteFromLink(c.BaseURL + "/" + noteID + "#" + password)
}
