package main

import (
	"context"
	"flag"
	"log"
	"os"

	noteofcli "github.com/NoteOf/noteof-cli"
	sdk "github.com/NoteOf/sdk-go"
	"github.com/google/subcommands"
)

func main() {
	config, err := noteofcli.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	editor := config.GetEditor()
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	if editor == "" {
		editor = "vi"
	}

	uapi := &sdk.UnauthenticatedAPI{}
	aapi := sdk.NewAuthenticatedApi(config.GetToken())

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&InitCmd{config, uapi}, "")

	list := &ListCmd{api: aapi}
	subcommands.Register(list, "")
	subcommands.Register(subcommands.Alias("ls", list), "")

	get := &GetCmd{api: aapi}
	subcommands.Register(get, "")
	subcommands.Register(subcommands.Alias("cat", get), "")

	delete := &DeleteCmd{api: aapi}
	subcommands.Register(delete, "")
	subcommands.Register(subcommands.Alias("rm", delete), "")

	archive := &ArchiveCmd{api: aapi}
	subcommands.Register(archive, "")

	subcommands.Register(&NewCmd{editor, config, aapi}, "")
	subcommands.Register(&EditCmd{editor, aapi}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
