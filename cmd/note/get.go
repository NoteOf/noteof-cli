package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/charmbracelet/glamour"
	"github.com/google/subcommands"
	"github.com/mattn/go-isatty"
)

type GetCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*GetCmd) Name() string     { return "get" }
func (*GetCmd) Synopsis() string { return "get a note" }
func (*GetCmd) Usage() string {
	return `get <noteID>:
	get a note.
`
}

func (p *GetCmd) SetFlags(f *flag.FlagSet) {}
func (p *GetCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Fatal("Expects exactly one noteID argument")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	out := n.CurrentText.NoteTextValue

	if isatty.IsTerminal(os.Stdout.Fd()) {
		out, err = glamour.Render(out, "notty")
		if err != nil {
			fmt.Println(n.CurrentText.NoteTextValue)
		}
	}

	fmt.Println(out)

	return subcommands.ExitSuccess
}
