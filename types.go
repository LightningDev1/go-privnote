package privnote

import (
	"encoding/json"
)

// Note is how a private note is represented in the Privnote API.
type Note struct {
	// Data is the encrypted content of the note.
	Data string `json:"data"`

	// HasManualPass is whether the note has a manual password.
	HasManualPass bool `json:"has_manual_pass"`

	// NoteLink is the URL of the note.
	NoteLink string `json:"note_link"`

	// DontAsk is whether the user should be asked for confirmation before
	// destroying the note.
	DontAsk bool `json:"dont_ask"`

	// Unknown fields
	Policy    int    `json:"policy"`
	ExpiresJS string `json:"expires_js"`
}

func (n *Note) UnmarshalJSON(data []byte) error {
	// Bool fields (HasManualPass, DontAsk) can be integers in the response.
	// So we need to unmarshal them as integers and convert them to bools.

	type Alias Note
	aux := &struct {
		HasManualPass int `json:"has_manual_pass"`
		DontAsk       int `json:"dont_ask"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	n.HasManualPass = aux.HasManualPass != 0
	n.DontAsk = aux.DontAsk != 0

	return nil
}
