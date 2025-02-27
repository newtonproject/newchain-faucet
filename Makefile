
BUILD_DATE=`date +%Y%m%d%H%M%S`
BUILD_COMMIT=`git rev-parse --short HEAD`

default:
	go build -ldflags "-X github.com/newtonproject/newchain-faucet/cli.buildCommit=${BUILD_COMMIT}\
	    -X github.com/newtonproject/newchain-faucet/cli.buildDate=${BUILD_DATE}" -o newchain-faucet

newton:
	@echo "Modifying go.mod for newton version..."
	cp go.mod go.mod.bak
	echo "replace github.com/ethereum/go-ethereum => github.com/newtonproject/newchain v1.9.18-newton-1.2" >> go.mod
	go mod tidy
	go build -ldflags "-X github.com/newtonproject/newchain-faucet/cli.buildCommit=${BUILD_COMMIT}\
		-X github.com/newtonproject/newchain-faucet/cli.buildDate=${BUILD_DATE}" -o newchain-faucet-newton
	mv go.mod.bak go.mod
	go mod tidy

all: default newton
