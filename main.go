package main

import (
	// "encoding/json"
	"fmt"
	"os"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
)

func main() {
	// ./trx-usdt u addr
	if len(os.Args) < 3 {
		fmt.Println("参数错误")
		return
	}

	// for idx, args := range os.Args {
		// fmt.Println("参数 ", idx, ": ", args)
	// }

	// TSQTUybJVrjFtPAwRMFTEoBxL5yUyGvEJM
	if os.Args[1] == "u" || os.Args[1] == "U" {
		GetBalance(os.Args[2])
		// GetTrc10Balance(os.Args[2])
		GetTrc20Balance(os.Args[2])
	} else if os.Args[1] == "t" || os.Args[1] == "T" {
		// 输入参数 ./trx-usdt t 20 fromAddr toAddr seedAddr secretKey
		// go run main.go t 10.39 TUzCsPXG9NZLer8FXV9iYi7PmudrbEoseJ TATnrJCSeqEqT56b8HVMgLgfBEZgZyS8dU TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj 252d7678e1c64b149a2494249f1c04164576ebb4fe4d9989fa534dbd54a1b852
		if len(os.Args) < 7 {
			fmt.Println("参数错误")
			return
		}

		// 转账金额
		money     := os.Args[2] // 20
		fromAddr  := os.Args[3] // 转账地址
		toAddr    := os.Args[4] // 收款地址
		seedAddr  := os.Args[5] // 合约地址
		secretKey := os.Args[6] // 私钥
		TransferTrc20(money, fromAddr, toAddr, seedAddr, secretKey)
	} else if os.Args[1] == "tt" || os.Args[1] == "TT" {
		// go run main.go tt 10 TUzCsPXG9NZLer8FXV9iYi7PmudrbEoseJ TATnrJCSeqEqT56b8HVMgLgfBEZgZyS8dU 252d7678e1c64b149a2494249f1c04164576ebb4fe4d9989fa534dbd54a1b852
		if len(os.Args) < 6 {
			fmt.Println("参数错误")
			return
		}

		money     := os.Args[2] // 20
		fromAddr  := os.Args[3] // 转账地址
		toAddr    := os.Args[4] // 收款地址
		secretKey := os.Args[5] // 私钥
		TransferTrx(money, fromAddr, toAddr, secretKey)
	} else {
		fmt.Println("暂不支持的查询类型")
	}
}

func GetBalance(addr string) {
	c, err := grpcs.NewClient("3.225.171.164:50051") // 正式网
	// c, err := grpcs.NewClient("grpc.nile.trongrid.io:50051") // 测试网
	if err != nil {
		fmt.Println(err)
		return
	}

	acc, err := c.GetTrxBalance(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// d, _ := json.Marshal(acc)

	// fmt.Println(string(d))

	fmt.Println(" ================================================")
	fmt.Println(" TRX Balance：", acc.GetBalance()) // 需要除以 1000000000000000000 后续测试发现是6个0
	fmt.Println(" TRX Balance Num：", fmt.Sprintf("%.6f", float64(acc.GetBalance())/1000000))
	fmt.Println(" ================================================")
}

// func GetTrc10Balance(addr string) {
// 	// c, err := grpcs.NewClient("grpc.trongrid.io:50051") // 正式网
// 	c, err := grpcs.NewClient("grpc.nile.trongrid.io:50051") // 测试网
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	amount, err := c.GetTrc10Balance(addr, "1002000")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(amount)
// }

func GetTrc20Balance(addr string) {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051") // 正式网
	// c, err := grpcs.NewClient("grpc.nile.trongrid.io:50051") // 测试网
	if err != nil {
		fmt.Println(err)
		return
	}
	// amount, err := c.GetTrc20Balance("TMgHF7GKoDEM1AvpgLVVizsPW7JVccDVuK", "TLdfZSUTwAJXxbav6od8iYCBSaW3EveYxm") // 代码测试的不要的
	amount, err := c.GetTrc20Balance(addr, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t") // 正式网
	// amount, err := c.GetTrc20Balance(addr, "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj") // 测试网
	if err != nil {
		fmt.Println(err)
		return
	}

	amountF, _ := strconv.ParseFloat(amount.String(), 64)

	fmt.Println(" ================================================")
	fmt.Println(" USDT TRC20 Balance：", amount.String()) // 需要除以 1000000
	fmt.Println(" USDT TRC20 Balance Num：", fmt.Sprintf("%.6f", amountF/1000000))
	fmt.Println(" ================================================")

	num2, _ := TokenWeiBigIntToEthStr(amount, 6)
	fmt.Println(" USDT TRC20 Balance Num2：", num2)
	fmt.Println(" ================================================")
}

func TransferTrx(money, fromAddr, toAddr, secretKey string) {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051") // 正式网
	// c, err := grpcs.NewClient("grpc.nile.trongrid.io:50051") // 测试网
	if err != nil {
		fmt.Println(err)
		return
	}

	amountF, _ := strconv.ParseFloat(money, 64)
	amountF = amountF * 1000000 // 精度
	amountF2 := strconv.FormatFloat(amountF, 'f', -1, 64)

	moneyInt, _ := strconv.ParseInt(amountF2, 10, 64)

	tx, err := c.Transfer(fromAddr, toAddr, moneyInt)
	if err != nil {
		fmt.Println(err)
		return
	}

	signTx, err := sign.SignTransaction(tx.Transaction, secretKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 交易哈希：", common.BytesToHexString(tx.GetTxid()))
}

func TransferTrc20(money, fromAddr, toAddr, seedAddr, secretKey string) {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051") // 正式网
	// c, err := grpcs.NewClient("grpc.nile.trongrid.io:50051") // 测试网
	if err != nil {
		fmt.Println(err)
		return
	}

	amountF, _ := strconv.ParseFloat(money, 64)
	amountF = amountF * 1000000 // 精度
	amountF2 := strconv.FormatFloat(amountF, 'f', -1, 64)

	moneyInt, _ := strconv.ParseInt(amountF2, 10, 64)

	amount := big.NewInt(moneyInt) // 转账金额 U
	// amount = amount.Mul(amount, big.NewInt(1000000)) // 精度
	tx, err := c.TransferTrc20(fromAddr, toAddr, seedAddr, amount, 100000000) // gaslimit
	if err != nil {
		fmt.Println(err)
		return
	}

	signTx, err := sign.SignTransaction(tx.Transaction, secretKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 交易哈希：", common.BytesToHexString(tx.GetTxid()))
}

func TokenWeiBigIntToEthStr(wei *big.Int, tokenDecimals int64) (string, error) {
	balance, err := decimal.NewFromString(wei.String())
	if err != nil {
		return "0", err
	}

	balanceStr := balance.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(tokenDecimals))).StringFixed(int32(tokenDecimals))

	return balanceStr, nil
}
