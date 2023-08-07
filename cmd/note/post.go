package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	noteofcli "github.com/NoteOf/noteof-cli"
	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type NewCmd struct {
	editor string

	config *noteofcli.Config
	api    *sdk.AuthenticatedAPI
}

func (*NewCmd) Name() string     { return "new" }
func (*NewCmd) Synopsis() string { return "post a new note" }
func (*NewCmd) Usage() string {
	return `new:
	post a new note.
`
}

func (p *NewCmd) SetFlags(f *flag.FlagSet) {}
func (p *NewCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	body, err := noteofcli.Edit(p.editor, "")
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	n := &sdk.Note{
		CurrentText: &sdk.NoteText{
			NoteTextValue: string(body),
		},
		Meta: sdk.NoteMeta{
			"created:with": "noteof-cli",
		},
	}

	n2, err := p.api.PostNewNote(n)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(n2.PublicID, getTitleLine(n2.CurrentText.NoteTextValue))

	return subcommands.ExitSuccess
}
