package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/client"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/shardping"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontspace/ontology-bench/sharding-bench/src/config"
)

func main() {
	log.InitLog(log.InfoLog, log.PATH, log.Stdout)
	if len(os.Args) < 3 {
		log.Errorf("missed config file and wallet path")
		return
	}
	configPath := os.Args[1]
	walletPath := os.Args[2]

	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		log.Error(err)
		return
	}

	wallet, err := sdk.OpenWallet(walletPath)
	if err != nil {
		log.Errorf("parse wallet err: %s", err)
		return
	}
	account, err := wallet.GetDefaultAccount([]byte(cfg.Password))
	if err != nil {
		log.Errorf("get account err: %s", err)
		return
	}

	shardTxTest(cfg, account)
}

func shardTxTest(cfg *config.Config, account *sdk.Account) {
	txNum := cfg.TxNum * cfg.ShardPerNode * len(cfg.Rpc)
	if txNum > math.MaxUint32 {
		txNum = math.MaxUint32
	}
	exitChan := make(chan int)
	routineNum := len(cfg.Rpc) * cfg.ShardPerNode
	txNumPerRoutine := txNum / routineNum
	startTestTime := time.Now().UnixNano() / 1e6
	for i, server := range cfg.Rpc {
		for j := 0; j < cfg.ShardPerNode; j++ {
			go func(ipaddr string, shardId uint64) {
				sendTxSdk := sdk.NewOntologySdk()
				rpcClient := client.NewRpcClient()
				rpcAddress := fmt.Sprintf("http://%s:%d", ipaddr, shardId*10+20336)
				rpcClient.SetAddress(rpcAddress)
				sendTxSdk.SetDefaultClient(rpcClient)

				for k := 0; j < txNumPerRoutine; k++ {
					txPayload := fmt.Sprintf("%d", k)
					if err := sendShardPing(sendTxSdk, cfg, account, shardId, 0, txPayload); err != nil {
						log.Errorf("send ping to %s, shard %d failed: %s", ipaddr, shardId, err)
						return
					}
					time.Sleep(time.Microsecond * 10)
				}
				exitChan <- 1
			}(server, uint64(i*cfg.ShardPerNode+j+1))
		}
	}
	for i := 0; i < routineNum; i++ {
		<-exitChan
	}
	endTestTime := time.Now().UnixNano() / 1e6
	log.Infof("send tps is %f", float64(txNum*1000)/float64(endTestTime-startTestTime))
}

func sendShardPing(sdk *sdk.OntologySdk, cfg *config.Config, user *sdk.Account, fromShardID, toShardID uint64, txt string) error {
	tFromShardId, _ := types.NewShardID(fromShardID)
	tToShardId, _ := types.NewShardID(toShardID)
	param := shardping.ShardPingParam{
		FromShard: tFromShardId,
		ToShard:   tToShardId,
		Param:     txt,
	}
	buf := new(bytes.Buffer)
	if err := param.Serialize(buf); err != nil {
		return fmt.Errorf("failed to ser shard deposit gas param: %s", err)
	}

	method := shardping.SEND_SHARD_PING_NAME
	contractAddress := utils.ShardPingAddress
	txParam := []interface{}{buf.Bytes()}

	txHash, err := sdk.Native.InvokeShardNativeContract(fromShardID, cfg.GasPrice, cfg.GasLimit, user, 0, contractAddress, method, txParam)
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("shard send shard %d ping txHash is :%s", fromShardID, txHash.ToHexString())
	return nil
}
