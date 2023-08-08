package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

type MetaCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*MetaCmd) Name() string     { return "meta" }
func (*MetaCmd) Synopsis() string { return "subcommands for working with note metadata" }
func (*MetaCmd) Usage() string {
	return `meta:
	subcommands for working with note metadata.
`
}

func (p *MetaCmd) SetFlags(f *flag.FlagSet) {}

func (p *MetaCmd) Execute(ctx context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	cmdr := subcommands.NewCommander(fs, subcommands.DefaultCommander.Name()+" "+p.Name())
	cmdr.Register(&MetaSetCmd{api: p.api}, "")
	cmdr.Register(&MetaGetCmd{api: p.api}, "")

	list := &MetaListCmd{api: p.api}
	cmdr.Register(list, "")
	cmdr.Register(subcommands.Alias("ls", list), "")

	delete := &MetaDeleteCmd{api: p.api}
	cmdr.Register(delete, "")
	cmdr.Register(subcommands.Alias("rm", delete), "")
	cmdr.Register(cmdr.HelpCommand(), "")

	return cmdr.Execute(ctx)
}

type MetaSetCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*MetaSetCmd) Name() string     { return "set" }
func (*MetaSetCmd) Synopsis() string { return "set a note metadata item" }
func (*MetaSetCmd) Usage() string {
	return `set <noteID> <key> <value>:
	set note metadata item by key.
`
}

func (p *MetaSetCmd) SetFlags(f *flag.FlagSet) {}

func (p *MetaSetCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 3 {
		log.Fatal("Expects exactly three arguments: <noteID> <key> <value>")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	n.Meta[fs.Arg(1)] = fs.Arg(2)
	_, err = p.api.PutUpdateNote(n)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

type MetaGetCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*MetaGetCmd) Name() string     { return "get" }
func (*MetaGetCmd) Synopsis() string { return "get a note metadata item" }
func (*MetaGetCmd) Usage() string {
	return `get <noteID> <key>:
	get note metadata item by key.
`
}

func (p *MetaGetCmd) SetFlags(f *flag.FlagSet) {}

func (p *MetaGetCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Fatal("Expects exactly two arguments: <noteID> <key>")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := n.Meta[fs.Arg(1)]
	if !ok {
		log.Println("key not found")
		return subcommands.ExitFailure
	}

	fmt.Println(v)

	return subcommands.ExitSuccess
}

type MetaListCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*MetaListCmd) Name() string     { return "list" }
func (*MetaListCmd) Synopsis() string { return "list note metadata items" }
func (*MetaListCmd) Usage() string {
	return `list <noteID>:
	list note metadata items.
`
}

func (p *MetaListCmd) SetFlags(f *flag.FlagSet) {}

func (p *MetaListCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Fatal("Expects exactly one argument: <noteID>")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range n.Meta {
		fmt.Println(k, v)
	}

	return subcommands.ExitSuccess
}

type MetaDeleteCmd struct {
	api *sdk.AuthenticatedAPI
}

func (*MetaDeleteCmd) Name() string     { return "delete" }
func (*MetaDeleteCmd) Synopsis() string { return "delete a note metadata item" }
func (*MetaDeleteCmd) Usage() string {
	return `delete <noteID> <key>:
	delete note metadata item by key.
`
}

func (p *MetaDeleteCmd) SetFlags(f *flag.FlagSet) {}

func (p *MetaDeleteCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Fatal("Expects exactly two arguments: <noteID> <key>")
	}

	i := fs.Arg(0)
	n, err := p.api.GetNote(i)
	if err != nil {
		log.Fatal(err)
	}

	key := fs.Arg(1)
	delete(n.Meta, key)
	_, err = p.api.PutUpdateNote(n)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(key, "DELETED")

	return subcommands.ExitSuccess
}
