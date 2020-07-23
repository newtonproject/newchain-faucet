package cli

import "testing"

func TestStart(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("start")
}
