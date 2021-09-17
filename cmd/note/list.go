package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type ListCmd struct {
	api      *sdk.AuthenticatedAPI
	archived bool
}

func (*ListCmd) Name() string     { return "list" }
func (*ListCmd) Synopsis() string { return "list your notes" }
func (*ListCmd) Usage() string {
	return `list:
	list your notes.
`
}

func (p *ListCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.archived, "a", false, "include archived notes")
}

func (p *ListCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	notes, err := p.api.GetNotes()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, n := range notes {
		if n.Archived && !p.archived {
			continue
		}

		fmt.Fprintf(w, "%s\t%s", n.PublicID, getTitleLine(n.CurrentText.NoteTextValue))
	}

	return subcommands.ExitSuccess
}

func getTitleLine(s string) string {
	return strings.Split(strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(s), "#")), "\n")[0]
}
