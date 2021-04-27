package main

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	bech322 "github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"runtime"
	"strconv"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	privateKey := ""
	host := ""
	times := 20000
	batchNum := 2
	interval := time.Second * 5
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)
	provider := provider2.NewProvider(host)
	chainIdString, err := provider.GetNetworkId()
	chainId, _ := strconv.ParseInt(chainIdString, 10, 64)
	if err != nil {
		panic(err)
	}
	address := keytools.GetAddressFromPrivateKey(util.DecodeHex(privateKey))
	bech32Addr, _ := bech322.ToBech32Address(address)
	pubKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)

	gasPrice, err := provider.GetMinimumGasPrice()
	if err != nil {
		panic(err)
	}

	balAndNonce, _ := provider.GetBalance(address)
	startNonce := balAndNonce.Nonce

	for i := 0; i < times; i++ {
		var transactions []*transaction.Transaction
		for i := 0; i < batchNum; i++ {
			txn := &transaction.Transaction{
				Version:      strconv.FormatInt(int64(util.Pack(int(chainId), 1)), 10),
				SenderPubKey: util.EncodeHex(pubKey),
				ToAddr:       bech32Addr,
				Amount:       "1000",
				GasPrice:     gasPrice,
				GasLimit:     "50",
				Code:         "",
				Data:         "",
				Priority:     true,
			}

			transactions = append(transactions, txn)
		}

		err2 := wallet.SignBatchWithNonce(transactions, *provider, startNonce+1)
		if err2 != nil {
			panic(err2)
		}
		batchSendingResult, err := wallet.SendBatchOneGo(transactions, *provider)
		if err != nil {
			panic(err)
		} else {
			for _, tx := range batchSendingResult {
				rsp, _ := json.Marshal(tx)
				fmt.Println(string(rsp))
			}
		}

		time.Sleep(interval)
		startNonce += int64(batchNum)

	}
}
