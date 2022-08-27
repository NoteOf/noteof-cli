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
	append bool

	api *sdk.AuthenticatedAPI
}

func (*EditCmd) Name() string     { return "edit" }
func (*EditCmd) Synopsis() string { return "edit a note" }
func (*EditCmd) Usage() string {
	return `edit <noteID>:
	edit a note.
`
}

func (p *EditCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.append, "append", false, "append to the note")
}

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

	startingText := ""
	if !p.append {
		startingText = n.CurrentText.NoteTextValue
	}

	body, err := noteofcli.Edit(editor, startingText)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	if p.append {
		out := strings.TrimSuffix(n.CurrentText.NoteTextValue, "\n")
		if out != "" {
			out += "\n"
		}

		n.CurrentText.NoteTextValue = out + strings.TrimSuffix(string(body), "\n") + "\n"
	} else {
		n.CurrentText.NoteTextValue = string(body)
	}

	n2, err := p.api.PutUpdateNote(n)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(n.PublicID, getTitleLine(n2.CurrentText.NoteTextValue))

	return subcommands.ExitSuccess
}
