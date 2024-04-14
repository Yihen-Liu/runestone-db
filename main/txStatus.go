package main

import (
	"context"
	"encoding/json"
	"fmt"
	"galactic-monitor/config"
	"galactic-monitor/log"
	"galactic-monitor/models"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const QueryTxStatusError = -1
const TxIsSuccess = 1
const TxIsFailed = 0
const TxIsPending = 2

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func verifyTriasTransfer(client *ethclient.Client, height uint64) {

}

func verifyUsdcTransfer(client *ethclient.Client, logs []*types.Log, from, to string, amount float64) bool {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	contractAbi, err := abi.JSON(strings.NewReader(USDCTokenAbi))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		if vLog.Address.Hex() != USDCTokenAddress {
			continue
		}

		if vLog.Topics[0].Hex() == logTransferSigHash.Hex() {
			var transferEvent LogTransfer

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			if transferEvent.From.Hex() != from {
				continue
			}

			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			if transferEvent.To.Hex() != to {
				continue
			}

			//transferEvent.Tokens.Div(big.NewInt())
		}
	}

	return false
}

type BscResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Status string `json:"status"`
	} `json:"result"`
}

func queryBscTxStatus(txhash string) int {
	baseUrl := config.AppConf.ChainRpc["BNBChain"]
	requestUrl := fmt.Sprintf("%s&module=transaction&action=gettxreceiptstatus&txhash=%s", baseUrl, txhash)
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("Error requesting data:", err)
		return QueryTxStatusError
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return QueryTxStatusError
	}

	var bscResp BscResponse
	err = json.Unmarshal(body, &bscResp)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return QueryTxStatusError
	}

	num, err := strconv.Atoi(bscResp.Result.Status)
	if err != nil {
		return QueryTxStatusError
	}
	return num
}

func queryTxStatus(txhash string, token string, receiver string, chainId int32) int {
	chain, err := models.QueryChainByReceiver(receiver, chainId)
	if err != nil {
		log.Errorf("query chain by receiver err:%s, receiver:%s, chainid:%d", err.Error(), receiver, chainId)
		return QueryTxStatusError
	}

	if strings.TrimSpace(chain.Label) == "BNB Chain" {
		return queryBscTxStatus(txhash)
	}

	client, err := ethclient.Dial(config.AppConf.ChainRpc[chain.Label])
	if err != nil {
		log.Errorf("txhash:%s, dial rpc err:%s, url:%s, chain:%s,", txhash, err.Error(), config.AppConf.ChainRpc[chain.Label], chain.Label)
		return QueryTxStatusError
	}
	defer client.Close()

	txHash := common.HexToHash(txhash)
	_, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatalf("query tx by hash err:%s, tx hash:%s, token:%s", err.Error(), txhash, token)
		return QueryTxStatusError
	}
	if isPending == true {
		return TxIsPending
	}

	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal("query tx receipt err:%s", err.Error())
		return QueryTxStatusError
	}
	// 1. 如果token是ETH，验证有效性

	//2. 如果是TRIAS，验证有效性

	//3. 如果是USDC，验证有效性

	return int(receipt.Status)
}

func StartMonitorNewOrder(ctx context.Context) {
	for {
		select {
		case <-time.After(3 * time.Second):
			log.Info("monitor new order")
			//1. 查出所有的pending状态的订单
			orders, err := models.QueryPendingOrders()
			if err != nil {
				log.Errorf("query pending orders err:", err.Error())
				continue
			}
			if len(orders) == 0 {
				log.Debug("no order need to be handled, wait some minutes...")
				continue
			}

			// 2. 去链上查询交易状态, 如果是1，则success； 如果是2， 则failed；
			for _, order := range orders {
				status := queryTxStatus(order.TxHash, order.Type, order.ReceiverAdddress, order.ChainId)
				if status == TxIsPending || status == QueryTxStatusError {
					continue
				}
				// 更新订单状态
				if status == TxIsFailed {
					models.UpdateOrderStatus(order.TxHash, "failed")
				}
				if status == TxIsSuccess {
					models.UpdateOrderStatus(order.TxHash, "success")
					//activeUserStatus(order.InitiatorAddress)
				}
				//bsc的apikey调用有限制，所以每次调用间隔1秒
				time.Sleep(time.Second * 1)
			}
		case <-ctx.Done():
			return
		}
	}
}
