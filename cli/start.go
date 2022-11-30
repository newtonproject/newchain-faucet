package cli

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [-port 8888] [--from address] [--amount 16888] [--unit NEW]",
		Short: "start " + cli.name + " server",
		Args:  cobra.MinimumNArgs(0),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			port := viper.GetInt("faucet.port")
			cli.port = port

			amountStr := viper.GetString("faucet.amount")
			unit := viper.GetString("faucet.unit")
			d := stringInSlice(unit, DenominationList)
			if !d {
				fmt.Printf("Unit(%s) for amount error. %s.\n", unit, DenominationString)
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}

			amountWei, ok := getAmountWei(amountStr, unit)
			if !ok {
				fmt.Println("Get amount error:", amountStr)
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}
			cli.amountWei = amountWei

			fromAddress := viper.GetString("faucet.from")
			if fromAddress == "" {
				fmt.Println("Error: required flag(s) \"from\" not set")
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}
			cli.coinbase = fromAddress

			walletPath := cli.walletPath
			rpcURL := cli.rpcURL

			ks := keystore.NewKeyStore(walletPath,
				keystore.LightScryptN, keystore.LightScryptP)
			if len(ks.Accounts()) == 0 {
				fmt.Println("Empty wallet, create account first.")
				return
			}

			account, err := ks.Find(accounts.Account{Address: common.HexToAddress(fromAddress)})
			if err != nil {
				fmt.Printf("Error: keystore find account(%s) error(%v)\n", fromAddress, err)
				return
			}

			walletPassword := viper.GetString("faucet.password")
			isSetPwd := viper.IsSet("faucet.password")
			var trials int
			for trials = 0; trials < 3; trials++ {
				prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", fromAddress, trials+1, 3)
				if isSetPwd == false {
					walletPassword, _ = getPassPhrase(prompt, false)
				} else {
					fmt.Println(prompt, "\nUse the `faucet.password` in the config file")
				}
				err = ks.Unlock(account, walletPassword)
				if err == nil {
					break
				}
				walletPassword = ""
			}
			cli.password = walletPassword

			if trials == 3 {
				fmt.Printf("Error: Failed to unlock account %s (%v)\n", fromAddress, err)
				return
			}

			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				fmt.Println("Dial error:", err)
				return
			}
			ctx := context.Background()
			nonce, err := client.NonceAt(ctx, account.Address, nil)
			if err != nil {
				fmt.Println("NonceAt error:", err)
				return
			}
			cli.nonce = nonce

			// get ChainID
			networkID, err := client.NetworkID(ctx)
			if err != nil {
				fmt.Println("Get NetworkID Error: ", err)
				networkID = big.NewInt(16888)
			}
			cli.networkID = networkID

			cli.startFaucet()

		},
	}

	cmd.Flags().String("from", "", "source account seed or name")
	unitUsageString := fmt.Sprintf("unit for faucet amount. %s.", DenominationString)
	cmd.Flags().StringP("unit", "u", "NEW", unitUsageString)
	cmd.Flags().StringP("amount", "a", "16888", "Default faucet amount")
	cmd.Flags().IntP("port", "p", 8888, "Default faucet server port `url`")

	viper.BindPFlag("faucet.from", cmd.Flags().Lookup("from"))
	viper.BindPFlag("faucet.unit", cmd.Flags().Lookup("unit"))
	viper.BindPFlag("faucet.amount", cmd.Flags().Lookup("amount"))

	viper.BindPFlag("faucet.port", cmd.Flags().Lookup("port"))

	return cmd
}
