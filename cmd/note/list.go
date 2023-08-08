package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type ListCmd struct {
	api      *sdk.AuthenticatedAPI
	archived bool
}

func (*ListCmd) Name() string     { return "list" }
func (*ListCmd) Synopsis() string { return "list notes" }
func (*ListCmd) Usage() string {
	return `list:
	list notes.
`
}

func (p *ListCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.archived, "a", false, "all -- include archived notes")
}

func (p *ListCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	notes, err := p.api.GetNotes()
	if err != nil {
		log.Fatal(err.Error())
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, n := range notes {
		if n.Archived && !p.archived {
			continue
		}

		fmt.Fprintf(w, "%s\t%s", n.PublicID, getTitleLine(n.CurrentText.NoteTextValue))

		fmt.Fprint(w, "\t[")
		if p.archived {
			if n.Archived {
				fmt.Fprint(w, "a")
			} else {
				fmt.Fprint(w, " ")
			}
		}

		if n.Starred {
			fmt.Fprint(w, "*")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, "]")

		fmt.Fprint(w, "\t", n.CurrentText.Created)

		fmt.Fprintln(w)
	}

	w.Flush()

	return subcommands.ExitSuccess
}

func getTitleLine(s string) string {
	out := strings.Split(strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(s), "#")), "\n")[0]
	trimmed, didtrim := substr(out, 50)
	if didtrim {
		return trimmed + "..."
	}

	return trimmed
}

func substr(input string, length int) (string, bool) {
	data := []rune(input)
	if len(data) < length {
		return input, false
	}

	return string(data[:length]), true
}
