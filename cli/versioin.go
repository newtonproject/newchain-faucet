package cli

import (
	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Get version of " + cli.name + " CLI",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			version := cli.version
			if isNewton() {
				version = version + "-newton"
			} else {
				version = version + "-ethereum"
			}
			showSuccess(version)
		},
	}

	return cmd
}

func isNewton() bool {
	p1 := crypto.S256().Params()
	p2 := elliptic.P256().Params()

	return p1.Gx.Cmp(p2.Gx) == 0 &&
		p1.Gy.Cmp(p2.Gy) == 0 &&
		p1.N.Cmp(p2.N) == 0 &&
		p1.B.Cmp(p2.B) == 0
}
