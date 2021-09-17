package main

import (
	"context"
	"flag"
	"log"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type ArchiveCmd struct {
	api *sdk.AuthenticatedAPI

	unarchive bool
}

func (*ArchiveCmd) Name() string     { return "archive" }
func (*ArchiveCmd) Synopsis() string { return "archive (or unarchive) a note" }
func (*ArchiveCmd) Usage() string {
	return `archive <noteID>:
	archive a note.
`
}

func (p *ArchiveCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.unarchive, "u", false, "unarchive")
}

func (p *ArchiveCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Fatal("Expects exactly one noteID argument")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	n.Archived = !p.unarchive

	if _, err := p.api.PutUpdateNote(n); err != nil {
		log.Fatal(err)
	}

	return subcommands.ExitSuccess
}
