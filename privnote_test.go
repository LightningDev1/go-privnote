package privnote

import (
	"math/rand"
	"testing"
)

// https://coder-ipsum.tech/
var randomNotes = []string{
	`Use the native DOM algorithm, then you can compile the idiosyncratic devops!`,
	`So bubble-sorting the document object model won't do anything, we need to circle back the test-driven LIFO stream!`,
	`You can't pair the resource without pair programming the containerized SOAP RSS feed!`,
	`I'll promise the functional XML language, that should Internet Explorer the CLI controller!`,
	`Try to ship the OOP emoji, maybe it will promise the dynamic graph!`,
	`The CLI child is down, graph the scalable devops so we can uglify the FIFO Netscape!`,
	`Try to rebase the CLI Imagemagick, maybe it will pivot the ecommerce code!`,
	`The SRE presenter is down, pair the mobile hashtable so we can pivot the FP callback!`,
	`We need to grep the atomic SQL data store!`,
	`I'll circle back the senior XML RSS feed, that should Internet Explorer the AWS language!`,
}

func TestMemoizer(t *testing.T) {
	t.Log("Testing go-privnote...")

	noteData := randomNotes[rand.Intn(len(randomNotes))]

	t.Logf("Using note data: \"%s\"", noteData)

	client := NewClient()

	// Creating a note
	noteLink, err := client.CreateNote(CreateNoteData{
		Data: noteData,
	})
	if err != nil {
		t.Fatal("Failed to create note:", err)
	}

	t.Log("Created note:", noteLink)

	// Reading the readNote
	readNote, err := client.ReadNoteFromLink(noteLink)
	if err != nil {
		t.Fatal("Failed to read note:", err)
	}

	t.Logf("Read note: \"%s\"", readNote)

	if readNote != noteData {
		t.Fatal("Note data does not match!")
	}

	t.Log("Done testing go-privnote.")
}
