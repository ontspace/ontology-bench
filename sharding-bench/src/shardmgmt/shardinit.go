package main

import (
	"fmt"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/smartcontract/service/native/shardmgmt"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontspace/ontology-bench/sharding-bench/src/config"
)

func ShardInit(sdk *sdk.OntologySdk, cfg *config.Config, user *sdk.Account) error {
	method := shardmgmt.INIT_NAME
	contractAddress := utils.ShardMgmtContractAddress
	txHash := common.Uint256{}
	var err error
	txHash, err = sdk.Native.InvokeNativeContract(cfg.GasPrice, cfg.GasLimit, user, 0,
		contractAddress, method, []interface{}{})
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("shard init txHash is :%s", txHash.ToHexString())
	return nil
}
