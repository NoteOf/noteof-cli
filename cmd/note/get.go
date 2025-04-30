package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/charmbracelet/glamour"
	"github.com/google/subcommands"
	"github.com/mattn/go-isatty"
)

type outputType string

const (
	outputJson   outputType = "json"
	outputPlain  outputType = "plain"
	outputPretty outputType = "pretty"
)

type GetCmd struct {
	api *sdk.AuthenticatedAPI

	output outputType
}

func (*GetCmd) Name() string     { return "get" }
func (*GetCmd) Synopsis() string { return "get a note" }
func (*GetCmd) Usage() string {
	return `get <noteID>:
	get a note.
`
}

func (p *GetCmd) SetFlags(f *flag.FlagSet) {
	p.output = outputPlain
	if isatty.IsTerminal(os.Stdout.Fd()) {
		p.output = outputPretty
	}

	f.Func("output", fmt.Sprintf("set the output format, options are: %s %s %s", outputPretty, outputPlain, outputPretty), func(s string) error {
		switch outputType(s) {
		case outputJson:
			p.output = outputJson
		case outputPlain:
			p.output = outputPlain
		case outputPretty:
			p.output = outputPretty
		default:
			return fmt.Errorf("invalid output type: %s", s)
		}

		return nil
	})
}

func (p *GetCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Fatal("Expects exactly one noteID argument")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	switch p.output {
	case outputPretty:
		out, err := glamour.Render(n.CurrentText.NoteTextValue, "notty")
		if err != nil {
			fmt.Println(n.CurrentText.NoteTextValue)
		}

		fmt.Println(out)
	case outputPlain:
		fmt.Println(n.CurrentText.NoteTextValue)
	case outputJson:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.Encode(n)
	}

	return subcommands.ExitSuccess
}
