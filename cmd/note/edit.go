package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	noteofitcli "github.com/Noteof/noteof-cli"
	sdk "github.com/Noteof/sdk-go"
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
	if fs.NArg() != 1 {
		log.Fatal("Expects exactly one noteID argument")
	}

	i, err := strconv.ParseInt(fs.Args()[0], 10, 64)
	if err != nil {
		log.Fatal("invalid id")
	}

	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	body, err := noteofitcli.Edit(p.editor, n.CurrentText.NoteTextValue)
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

	fmt.Println(n.NoteID, getTitleLine(n2.CurrentText.NoteTextValue))

	return subcommands.ExitSuccess
}
