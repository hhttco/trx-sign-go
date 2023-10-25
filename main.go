package main

import (
	// "encoding/json"
	"fmt"
	"os"
	"github.com/JFJun/trx-sign-go/grpcs"
	// "math/big"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("参数错误")
		return
	}

	// for idx, args := range os.Args {
		// fmt.Println("参数 ", idx, ": ", args)
	// }

	// TSQTUybJVrjFtPAwRMFTEoBxL5yUyGvEJM
	if os.Args[1] == "u" || os.Args[1] == "U" {
		GetTrc20Balance(os.Args[2])
		GetBalance(os.Args[2])
	} else {
		fmt.Println("暂不支持的查询类型")
	}
}


func GetTrc20Balance(addr string) {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051")
	if err != nil {
		fmt.Println(err)
		return
	}
	// amount, err := c.GetTrc20Balance("TMgHF7GKoDEM1AvpgLVVizsPW7JVccDVuK", "TLdfZSUTwAJXxbav6od8iYCBSaW3EveYxm")
	amount, err := c.GetTrc20Balance(addr, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if err != nil {
		fmt.Println(err)
		return
	}

	amountF, _ := strconv.ParseFloat(amount.String(), 64)

	fmt.Println(" ================================================")
	fmt.Println(" USDT Balance：", amount.String()) // 需要除以 1000000
	fmt.Println(" USDT Balance Num：", fmt.Sprintf("%.6f", amountF/1000000))
	fmt.Println(" ================================================")
}

func GetBalance(addr string) {
	c, err := grpcs.NewClient("3.225.171.164:50051")
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
	fmt.Println(" TRX Balance：", acc.GetBalance()) // 需要除以 1000000000000000000
	fmt.Println(" TRX Balance Num：", fmt.Sprintf("%.18f", float64(acc.GetBalance())/1000000000000000000))
	fmt.Println(" ================================================")
}
