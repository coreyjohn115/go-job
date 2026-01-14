package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

func Test1(t *testing.T) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	ctx := context.Background()

	resp, err := c.GetVersion(ctx)
	if err != nil {
		log.Fatalf("GetVersion: %v", err)
	}
	log.Println("GetVersion", resp.SolanaCore)

	address := os.Getenv("address")
	balance, err := c.GetBalance(ctx, address)
	fmt.Println(balance)

	info, err := c.GetAccountInfo(ctx, address)
	if err != nil {
		log.Fatalf("GetAccountInfo: %v", err)
	}
	log.Printf("GetAccountInfo%v", info)

	balance1, err1 := c.GetBalanceWithConfig(context.TODO(), address, client.GetBalanceConfig{
		Commitment: rpc.CommitmentFinalized,
	})
	if err1 != nil {
		log.Fatalf("GetBalanceWithConfig: %v", err1)
	}
	log.Printf("GetBalanceWithConfig %v", balance1)

	// 获取最新的区块高度
	slot, err := c.GetSlot(ctx)
	if err != nil {
		log.Fatal("获取最新slot失败:", err)
	}
	fmt.Printf("最新slot: %d\n", slot)
	// 获取最新区块
	recentBlock, err := c.GetBlock(ctx, slot)
	if err != nil {
		panic("查询失败: " + err.Error())
	}

	fmt.Printf("区块高度: %d\n", recentBlock.BlockHeight)
	fmt.Printf("交易数量: %d\n", len(recentBlock.Transactions))
}

// 实时交易监控器
type TransactionMonitor struct {
	client        *client.Client
	lastSignature string
}

func TestMonitorStart(t *testing.T) {
	ctx := context.Background()
	transactionMonitor := &TransactionMonitor{
		client: client.NewClient(rpc.DevnetRPCEndpoint),
	}

	fmt.Println("🔍 开始监控交易...")
	ticker := time.NewTicker(time.Duration(6000))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 停止监控")
			return
		case <-ticker.C:
			transactionMonitor.pollRecentTransactions(ctx)
		}
	}
}

// 轮询最近交易
func (tm *TransactionMonitor) pollRecentTransactions(ctx context.Context) {
	// 获取最近的区块
	slot, err := tm.client.GetSlot(ctx)
	if err != nil {
		log.Printf("获取区块失败: %v", err)
		return
	}

	// 获取区块的交易签名
	blockSignatures, err := tm.client.GetBlock(ctx, slot)
	if err != nil {
		log.Printf("获取区块交易失败: %v", err)
		return
	}

	if blockSignatures == nil || len(blockSignatures.Transactions) == 0 {
		return
	}

	// 处理新交易
	for _, txSig := range blockSignatures.Transactions {
		tm.processTransaction(ctx, string(txSig.Transaction.Signatures[0]))
	}
}

// 处理单个交易
func (tm *TransactionMonitor) processTransaction(ctx context.Context, signature string) {
	// 如果是已经处理过的交易，跳过
	if signature == tm.lastSignature {
		return
	}

	// 获取交易详情
	tx, err := tm.client.GetTransaction(ctx, signature)
	if err != nil {
		log.Printf("获取交易详情失败: %v", err)
		return
	}

	if tx == nil {
		return
	}

	fmt.Printf("\n发现新交易: %s\n", signature)

	// 分析交易
	tm.analyzeTransaction(tx)
	tm.lastSignature = signature
}

// 分析交易
func (tm *TransactionMonitor) analyzeTransaction(tx *client.Transaction) {
	fmt.Printf("📊 交易分析:\n")
	fmt.Printf("  区块: %d\n", tx.Slot)

	if tx.BlockTime != nil {
		timestamp := time.Unix(int64(*tx.BlockTime), 0)
		fmt.Printf("  时间: %s\n", timestamp.Format("2006-01-02 15:04:05"))
	}

	// 检查交易状态
	if tx.Meta != nil {
		if tx.Meta.Err != nil {
			fmt.Printf("  状态: ❌ 失败\n")
			fmt.Printf("  错误: %v\n", tx.Meta.Err)
		} else {
			fmt.Printf("  状态: ✅ 成功\n")
		}

		// 计算费用
		if tx.Meta.Fee != 0 {
			fmt.Printf("  费用: %d lamports (%.6f SOL)\n",
				tx.Meta.Fee, float64(tx.Meta.Fee)/1e9)
		}

		// 计算单元
		if tx.Meta.ComputeUnitsConsumed != nil {
			fmt.Printf("  计算单元: %d\n", *tx.Meta.ComputeUnitsConsumed)
		}
	}
}
