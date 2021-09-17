package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	noteofcli "github.com/NoteOf/noteof-cli"
	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type EditCmd struct {
	editor string

	api *sdk.AuthenticatedAPI
}

func (*EditCmd) Name() string     { return "edit" }
func (*EditCmd) Synopsis() string { return "edit a note" }
func (*EditCmd) Usage() string {
	return `edit <noteID>:
	edit a note.
`
}

func (p *EditCmd) SetFlags(f *flag.FlagSet) {}
func (p *EditCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Fatal("Expects exactly one noteID argument")
	}

	editor := p.editor
	if fs.NArg() > 1 {
		editor = strings.Join(fs.Args()[1:], " ")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	body, err := noteofcli.Edit(editor, n.CurrentText.NoteTextValue)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	n.CurrentText.NoteTextValue = string(body)

	n2, err := p.api.PutUpdateNote(n)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(n.PublicID, getTitleLine(n2.CurrentText.NoteTextValue))

	return subcommands.ExitSuccess
}
