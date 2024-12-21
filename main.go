package main

//module github.com/Oleg323-creator

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strings"
)

var erc20ABI = `[{"constant":false,"inputs":[{"name":"recipient","type":"address"},
{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],
"payable":false,"stateMutability":"nonpayable","type":"function"}]`

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/8ab0a7925db44ab094fa1b6546409a6b")
	if err != nil {
		log.Fatalf("Error connecting to Sepolia node: %v", err)
	}
	fmt.Println("Connection success!")

	tokenAddress := common.HexToAddress("0xB2A727C1250EFEF3766D6928EC2C65EF50dc0557")
	fromAddress := common.HexToAddress("0x886577048713f65d6e26e61e82597A523887645B")
	toAddress := common.HexToAddress("0xAa8ff5ed1dA7832b6b361F90b9bA6D7b384Ea5E9")

	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Error getting nonce: %v", err)
	}
	fmt.Printf("Using nonce: %d\n", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting gas price: %v", err)
	}

	gasLimit := uint64(100000)

	tokenABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf("Error readind ABI contract: %v", err)
	}

	amount := new(big.Int)
	amount.SetString("1000000000000000", 10) //WEI

	data, err := tokenABI.Pack("transfer", toAddress, amount)
	if err != nil {
		log.Fatalf("Error using data for method: %v", err)
	}

	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gasPrice, data)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(11155111)), privateKey)
	if err != nil {
		log.Fatalf("Error singing tx: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error sending tx: %v", err)
	}

	log.Printf("Tx has sent! Hash: %s\n", signedTx.Hash().Hex())
}
