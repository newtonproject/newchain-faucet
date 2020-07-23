package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sirupsen/logrus"
)

func showSuccess(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

func showError(fields logrus.Fields, msg string, args ...interface{}) {
	logrus.WithFields(fields).Errorf(msg, args...)
}

// getPassPhrase retrieves the password associated with an account,
// requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool) (string, error) {
	// prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Enter passphrase (empty for no passphrase): ")
	if err != nil {
		return "", err
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Enter same passphrase again: ")
		if err != nil {
			return "", err
		}
		if password != confirm {
			return "", fmt.Errorf("Passphrases do not match")
		}
	}
	return password, nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// DenominationString is for denomination string
//const DenominationString = "Available unit: Wei, Ada, Babbage, Shannon, Szabo, Finney, Ether, Einstein, Douglas, Gwei"
const DenominationString = "Available unit: NEW, WEI"

// DenominationList is array for denomination string
// var DenominationList = []string{"Wei", "Ada", "Babbage", "Shannon", "Szabo", "Finney", "Ether", "Einstein", "Douglas", "Gwei"}
var DenominationList = []string{"NEW", "WEI"}

func getDenominationByUnit0(unit string) *big.Float {
	bf := new(big.Float)
	switch unit {
	case "Ether":
		bf.SetFloat64(params.Ether)
	case "Wei":
		bf.SetFloat64(params.Wei)
	default:
		bf.SetFloat64(params.Wei)
	}

	return bf
}

func getDenominationByUnit(unit string) *big.Float {
	bf := new(big.Float)
	switch unit {
	case "NEW":
		bf.SetFloat64(params.Ether)
	case "WEI":
		bf.SetFloat64(params.Wei)
	default:
		bf.SetFloat64(params.Wei)
	}

	return bf
}

func getAmountWei(amountStr, unit string) (*big.Int, bool) {

	amountFloat, ok := new(big.Float).SetString(amountStr)
	if !ok {
		return nil, ok
	}
	amountFloat.Mul(amountFloat, getDenominationByUnit(unit))
	amount := new(big.Int)
	amountFloat.Int(amount)

	return amount, ok
}

func getAmountTextByUnit(amount *big.Float, unit string) string {
	amountF := amount.Quo(amount, getDenominationByUnit(unit))

	precF := getDenominationByUnit(unit).Quo(getDenominationByUnit(unit), getDenominationByUnit("Wei"))
	precI := new(big.Int)
	precF.Int(precI)
	prec := len(precI.String()) - 1

	text := amountF.Text('f', prec)

	if unit == "NEW" && strings.Contains(text, ".") {
		for i := len(text) - 1; i > 0; i-- {
			if text[i] != '0' {
				if text[i] == '.' {
					return text[:i]
				}
				return text[:i+1]
			}
		}
	}

	return text
}

func getAmountTextUnitByUnit(amount *big.Float, unit string) string {
	amountF := amount.Quo(amount, getDenominationByUnit(unit))

	precF := getDenominationByUnit(unit).Quo(getDenominationByUnit(unit), getDenominationByUnit("Wei"))
	precI := new(big.Int)
	precF.Int(precI)
	prec := len(precI.String()) - 1

	text := amountF.Text('f', prec)

	if unit == "NEW" && strings.Contains(text, ".") {
		for i := len(text) - 1; i > 0; i-- {
			if text[i] != '0' {
				if text[i] == '.' {
					return text[:i]
				}
				return text[:i+1]
			}
		}
	}

	return text
}

func getWeiAmountTextByUnit(amount *big.Int, unit string) string {
	amountF := new(big.Float)
	amountF.SetInt(amount)

	return getAmountTextByUnit(amountF, unit)
}

func getWeiAmountTextUnitByUnit1(amount *big.Int, unit string) string {
	if unit == "" {
		NEW1 := "1000000000000000000"
		new1Big := new(big.Int)
		new1Big.SetString(NEW1, 10)
		if amount.Cmp(new1Big) < 0 {
			// show in WEI
			return amount.String() + " WEI"
		}

		unit = "NEW"
	}

	amountF := new(big.Float)
	amountF.SetInt(amount)

	return getAmountTextByUnit(amountF, unit) + " " + unit

}

func getWeiAmountTextUnitByUnit(amount *big.Int, unit string) string {
	if unit == "" {
		NEW1 := "1000000000000000000"
		new1Big := new(big.Int)
		new1Big.SetString(NEW1, 10)
		if amount.Cmp(new1Big) < 0 {
			// show in WEI
			return amount.String() + " WEI"
		}

		unit = "NEW"
	}

	if unit == "WEI" {
		return amount.String() + " WEI"
	} else if unit == "NEW" {
		const NEW_WEI = 18
		amountStr := amount.String()
		amountLen := len(amountStr)

		if amountLen <= NEW_WEI {
			for i := 0; i <= NEW_WEI+1-amountLen; i++ {
				amountStr = "0" + amountStr
			}
			if len(amountStr) != NEW_WEI+1 {
				return amount.String() + " WEI"
			}
			amountLen = NEW_WEI + 1
		}
		return amountStr[:amountLen-NEW_WEI] + "." + amountStr[amountLen-NEW_WEI:] + " NEW"
	}

	return amount.String() + " WEI"

}

func createNewAccount(walletPath string, numOfNew int) error {

	wallet := keystore.NewKeyStore(walletPath,
		keystore.LightScryptN, keystore.LightScryptP)

	walletPassword, err := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	for i := 0; i < numOfNew; i++ {
		account, err := wallet.NewAccount(walletPassword)
		if err != nil {
			fmt.Println("Account error:", err)
			return err
		}
		fmt.Println(account.Address.Hex())
	}

	return nil
}

// TODO: add transaction receipt later by command `receipt`
func showTransactionReceipt(url, txStr string) {
	var jsonStr = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["%s"],"id":1}`, txStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	clientHttp := &http.Client{}

	resp, err := clientHttp.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var body json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			fmt.Println(err)
			return
		}

		bodyStr, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			fmt.Println("JSON marshaling failed: ", err)
			return
		}
		fmt.Printf("%s\n", bodyStr)

		return
	}

}
