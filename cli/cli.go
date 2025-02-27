package cli

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	buildCommit string
	buildDate   string
)

// CLI represents a command-line interface. This class is
// not threadsafe.
type CLI struct {
	rootCmd    *cobra.Command
	version    string
	walletPath string
	rpcURL     string
	config     string
	testing    bool
	name       string

	port      int
	networkID *big.Int
	amountWei *big.Int
	nonce     uint64
	password  string
	coinbase  string
}

// NewCLI returns an initialized CLI
func NewCLI() *CLI {
	version := "v0.5.0"
	if buildCommit != "" {
		version = fmt.Sprintf("%s-%s", version, buildCommit)
	}
	if buildDate != "" {
		version = fmt.Sprintf("%s-%s", version, buildDate)
	}

	cli := &CLI{
		rootCmd:    nil,
		version:    version,
		walletPath: "",
		rpcURL:     "",
		testing:    false,
		config:     "",
		name:       "NewChainFaucet",
	}

	cli.buildRootCmd()
	return cli
}

// Execute parses the command line and processes it.
func (cli *CLI) Execute() {
	cli.rootCmd.Execute()
}

// setup turns up the CLI environment, and gets called by Cobra before
// a command is executed.
func (cli *CLI) setup(cmd *cobra.Command, args []string) {
	err := setupConfig(cli)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}
}

func (cli *CLI) help(cmd *cobra.Command, args []string) {
	fmt.Fprint(os.Stderr, cmd.UsageString())

	os.Exit(-1)

}

// TestCommand test command
func (cli *CLI) TestCommand(command string) string {
	cli.testing = true
	result := cli.Run(strings.Fields(command)...)
	cli.testing = false
	return result
}

// Run executes CLI with the given arguments. Used for testing. Not thread safe.
func (cli *CLI) Run(args ...string) string {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	cli.rootCmd.SetArgs(args)
	cli.rootCmd.Execute()
	cli.buildRootCmd()

	w.Close()

	os.Stdout = oldStdout

	var stdOut bytes.Buffer
	io.Copy(&stdOut, r)
	return stdOut.String()
}

// Embeddable returns a CLI that you can embed into your own Go programs. This
// is not thread-safe.
func (cli *CLI) Embeddable() *CLI {
	cli.testing = true
	return cli
}
