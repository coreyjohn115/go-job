package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	rawUrl      string
	privateKey1 string
	privateKey2 string
)

// 加载环境变量
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rawUrl = os.Getenv("SEPOLIA_RPC_URL")
	privateKey1 = os.Getenv("PRIVATE_KEY1")
	privateKey2 = os.Getenv("PRIVATE_KEY2")
}

// 查询最新区块
func TestQueryBlock(t *testing.T) {
	loadEnv()

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block hash: ", block.Hash().Hex())
	fmt.Println("block number: ", block.Number().Uint64())
	fmt.Println("block time: ", block.Time())
	fmt.Println("block nonce: ", block.Nonce())
	fmt.Println("block transactions: ", len(block.Transactions()))
}

// 转账
func TestTransfer(t *testing.T) {
	loadEnv()

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	keyFrom, err := crypto.HexToECDSA(privateKey1)
	if err != nil {
		log.Fatal(err)
	}
	keyTo, err := crypto.HexToECDSA(privateKey2)
	if err != nil {
		log.Fatal(err)
	}

	gassPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(keyFrom.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	toAddress := crypto.PubkeyToAddress(keyTo.PublicKey)
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(50000000000), 21000, gassPrice, nil)
	tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), keyFrom)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}

	// 0x29d3a2101ca64ce4addee060d6e2ffc39847edde4d66153b0444cea6e7645ad8
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
}

// 部署合约
func TestDeployContract(t *testing.T) {
	loadEnv()

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	keyFrom, err := crypto.HexToECDSA(privateKey1)
	if err != nil {
		log.Fatal(err)
	}

	gassPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(keyFrom.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(keyFrom, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.GasPrice = gassPrice
	auth.GasLimit = 210000
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	contractAddress, tx, _, err := DeployMain(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	//0xf2842d100d3ed350ea78ba75d7e4d3abb02e572ef004c028e580b2cf9efe68ea
	fmt.Printf("合约部署交易哈希: %s\n", tx.Hash().Hex())
	//0xA7de9761c90B6A993d5039732278fdef0F85C819
	fmt.Printf("合约地址: %s\n", contractAddress.Hex())
}

// 调用合约
func TestCallContract(t *testing.T) {
	loadEnv()

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	keyFrom, err := crypto.HexToECDSA(privateKey1)
	if err != nil {
		log.Fatal(err)
	}

	gassPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(keyFrom.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(keyFrom, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.GasPrice = gassPrice
	auth.GasLimit = 210000
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	contractAddress := common.HexToAddress("A7de9761c90B6A993d5039732278fdef0F85C819")
	instance, err := NewMain(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	count, err := instance.Count(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("合约计数: %s\n", count.String())

	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Increment 交易哈希: ", tx.Hash().Hex())
}
