package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Command struct {
	Name        string
	ShortName   string
	Usage       string
	Description string
	Commands    []Command
	Flags       []Flag
	Action      func(context *Context)
}

func (c Command) Run(ctx *Context) {
	// append help to flags
	c.Flags = append(
		c.Flags,
		helpFlag{"show help"},
	)

	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	firstFlagIndex := -1
	for index, arg := range ctx.Args() {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}

	var err error
	if firstFlagIndex > -1 {
		args := ctx.Args()[1:firstFlagIndex]
		flags := ctx.Args()[firstFlagIndex:]
		err = set.Parse(append(flags, args...))
	} else {
		err = set.Parse(ctx.Args()[1:])
	}

	if err != nil {
		fmt.Println("Incorrect Usage.\n")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		os.Exit(1)
	}

	context := NewContext(ctx.App, set, ctx.globalSet)
	checkCommandHelp(context, c.Name)

  args := context.Args()
	if len(args) > 0 {
		name := args[0]
		cmd := c.Command(name)
		if cmd != nil {
			cmd.Run(context)
			return
		}
	}

	c.Action(context)
}

func (c Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}

func (c Command) Command(name string) *Command {
	for _, cmd := range c.Commands {
		if cmd.HasName(name) {
			return &cmd
		}
	}

	return nil
}
