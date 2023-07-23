package main

import (
	"fmt"

	"github.com/LightningDev1/go-privnote"
)

func main() {
	client := privnote.NewClient()

	// Create a note.
	noteLink, err := client.CreateNote(privnote.CreateNoteData{
		Data: "Hello, World!",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(noteLink)

	// Read the note.
	note, err := client.ReadNoteFromLink(noteLink)
	if err != nil {
		panic(err)
	}

	fmt.Println(note)
}
