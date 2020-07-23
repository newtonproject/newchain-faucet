## NewChainFaucet

This is a server program that can get free money for the address of the NewChain blockchain.
* Call faucet to get free money
* Call balance to check your money

## QuickStart

### Download from releases

Binary archives are published at https://release.cloud.diynova.com/newton/NewChainFaucet/.

### Building the source

To get from gitlab via `go get`, this will get source and executable program.

#### Windows

install command

```bash
git clone https://github.com/newtonproject/newchain-faucet && cd newchain-faucet && make install
```

run NewChainFaucet

```bash
%GOPATH%/bin/NewChainFaucet.exe
```

#### Linux or Mac

install:

```bash
git clone https://github.com/newtonproject/newchain-faucet && cd newchain-faucet && make install
```
run NewChainFaucet

```bash
$GOPATH/bin/NewChainFaucet
```

### Usage

#### Help

Use command `NewChainFaucet help` to display the usage.

```bash
Usage:
  NewChainFaucet [flags]
  NewChainFaucet [command]

Available Commands:
  account     Manage NewChain accounts
  help        Help about any command
  init        Initialize config file
  start       start NewChainFaucet server
  version     Get version of NewChainFaucet CLI

Flags:
  -c, --config path            The path to config file (default "./config.toml")
  -h, --help                   help for NewChainFaucet
  -i, --rpcURL url             Geth json rpc or ipc url (default "https://rpc1.newchain.newtonproject.org")
  -w, --walletPath directory   Wallet storage directory (default "./wallet/")

Use "NewChainFaucet [command] --help" for more information about a command.
```

#### Use config.toml

You can use a configuration file to simplify the command line parameters.

One available configuration file `config.toml` is as follows:


```conf
rpcurl = "http://192.168.168.33"
walletpath = "./wallet/"

[faucet]
  amount = "16888"
  from = "0x83B4aB41173385A265788b835d8Ee5d3b84081D4"
  port = 8888
  unit = "NEW"
  password = "newton"
```

#### Initialize config file

```bash
# Initialize config file
newchain-faucet init
```

Just press Enter to use the default configuration, and it's best to create a new user.

```bash
$ NewChainFaucet init
Initialize config file
Enter file in which to save (./config.toml):
Enter the wallet storage directory (./wallet/):
Enter geth json rpc or ipc url (https://rpc1.newchain.newtonproject.org):
Create a new account or not: [Y/n]
Your new account is locked with a password. Please give a password. Do not forget this password.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
0xF69501bE4271eEC6AA21b975240b34f088bAb754
Your configuration has been saved in  ./config.toml
```

#### Create account

```bash
# Create an account
newchain-faucet account new

# Create 10 accounts
newchain-faucet account new -n 10
```

### List all accounts

```bash
# list all accounts of the walletPath
newchain-faucet account list
```

### Start faucet server

Make sure the default wallet address has a large balance before starting the server

```bash
# Start faucet server with rpc url
newchain-faucet start
```

### Get balance

* Open the browser and enter the url http://localhost:8888/balance?address=0xDB2C9C06E186D58EFe19f213b3d5FaF8B8c99481

* Use curl command

```bash
# Use curl command
curl http://localhost:8888/balance?address=0xDB2C9C06E186D58EFe19f213b3d5FaF8B8c99481
```

### Get faucet

* Open the browser and enter the url http://localhost:8888/faucet?address=0xDB2C9C06E186D58EFe19f213b3d5FaF8B8c99481

* Use curl command

```bash
# Use curl command
curl http://localhost:8888/faucet?address=0xDB2C9C06E186D58EFe19f213b3d5FaF8B8c99481
```
