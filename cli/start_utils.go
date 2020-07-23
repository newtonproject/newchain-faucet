package cli

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (cli *CLI) startFaucet() {
	port := cli.port
	portStr := fmt.Sprintf("%v", port)
	http.HandleFunc("/faucet", cli.faucetHandler)
	http.HandleFunc("/balance", cli.getBalanceHandler)
	addr := ":" + portStr
	fmt.Printf("Faucet serve started(%v)\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (cli *CLI) sendMoney(toAddressStr string) error {
	// to address
	if !common.IsHexAddress(toAddressStr) {
		return fmt.Errorf("Not valid hex-encoded address")
	}
	toAddress := common.HexToAddress(toAddressStr)

	// from address
	wallet := keystore.NewKeyStore(cli.walletPath,
		keystore.LightScryptN, keystore.LightScryptP)
	if len(wallet.Accounts()) == 0 {
		return fmt.Errorf("Empty wallet, create account first")
	}
	var account accounts.Account
	for _, a := range wallet.Accounts() {
		if a.Address == common.HexToAddress(cli.coinbase) {
			account = a
			break
		}
	}
	if account == (accounts.Account{}) {
		return fmt.Errorf("Error: Can NOT get the keystore file of address %v", cli.coinbase)
	}
	wallet.Unlock(account, cli.password)
	fmt.Println("account address:", account.Address.Hex())

	// amount
	amountWei := cli.amountWei
	// networ id
	networkID := cli.networkID

	// get gasLimit and gasPrice
	client, err := ethclient.Dial(cli.rpcURL)
	if err != nil {
		log.Printf("client dial error: %v", err)
		return err
	}
	defer client.Close()
	ctx := context.Background()
	var gasPrice *big.Int
	gasPrice, err = client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("SuggestGasPrice err:", err)
		gasPrice = big.NewInt(1)
	}
	msg := ethereum.CallMsg{
		From:     account.Address,
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    amountWei,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		fmt.Println("EstimateGas Error: ", err)
		gasLimit = 21000
	}

	// nonce
	nonce, _ := client.NonceAt(ctx, account.Address, nil)
	if cli.nonce < nonce {
		cli.nonce = nonce
	}
	fmt.Println("nonce: ", cli.nonce)

	tx := types.NewTransaction(cli.nonce, toAddress, amountWei, gasLimit, gasPrice, nil)
	signTx, err := wallet.SignTx(account, tx, networkID)

	err = client.SendTransaction(ctx, signTx)
	if err != nil {
		return fmt.Errorf("SendTransaction err (%v)", err)
	}

	cli.nonce++
	return nil
}

func (cli *CLI) getBalance(address string) *big.Int {
	client, err := ethclient.Dial(cli.rpcURL)
	if err != nil {
		log.Printf("client dial error: %v", err)
		return new(big.Int)
	}
	defer client.Close()
	ctx := context.Background()
	fmt.Println(address)
	balance, err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		fmt.Println("Balance error:", err)
		os.Exit(1)
	}
	fmt.Println("Balance: ", balance.String())
	return balance
}

func (cli *CLI) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	val, ok := r.Form["address"]
	if !ok {
		fmt.Fprintf(w, "Just give me a address!")
		return
	}
	if len(val) != 1 {
		fmt.Fprintf(w, "Just give me ONE address!")
		return
	}
	address := val[0]
	// TODO: check address is valid.
	log.Printf("faucet got address: %v", address)
	amount := cli.getBalance(address)

	fmt.Fprintf(w, "balance: %v", getWeiAmountTextUnitByUnit(amount, ""))
}

func (cli *CLI) faucetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	val, ok := r.Form["address"]
	if !ok {
		fmt.Fprintf(w, "Just give me a address!")
		return
	}
	if len(val) != 1 {
		fmt.Fprintf(w, "Just give me ONE address!")
		return
	}
	address := val[0]
	// TODO: check address is valid.
	log.Printf("faucet got address: %v", address)
	err := cli.sendMoney(address)
	if err != nil {
		fmt.Fprintf(w, "something is wrong: %v", err)
		return
	} else {
		fmt.Fprintf(w, "Done! go check your money.") // send data to client side
	}
}
