package main

import (
	cmd "github.com/Ghamster0/os-rq-fsender/cmd/fsender/command"
	"github.com/Ghamster0/os-rq-fsender/pkg/command"
)

func main() {
	command.Execute(cmd.Root)
}
