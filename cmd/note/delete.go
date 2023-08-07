package main

import (
	"context"
	"flag"
	"fmt"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type DeleteCmd struct {
	api *sdk.AuthenticatedAPI
	yes bool
}

func (*DeleteCmd) Name() string     { return "delete" }
func (*DeleteCmd) Synopsis() string { return "delete a note" }
func (*DeleteCmd) Usage() string {
	return `delete [<noteID>...]:
	delete one or more notes.
`
}

func (p *DeleteCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.yes, "y", false, "force")
}
func (p *DeleteCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	hasError := false

	if fs.NArg() == 0 {
		fmt.Println("no noteIDs given")
		return subcommands.ExitUsageError
	}

	for _, v := range fs.Args() {
		deleted, err := p.api.DeleteNote(v)
		if deleted && err == nil {
			fmt.Println(v, "DELETED")
		} else if err != nil {
			fmt.Println(v, err)
			hasError = true
		} else {
			fmt.Println(v, "NOT FOUND")
			hasError = true
		}
	}

	if hasError {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
