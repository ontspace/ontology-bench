package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/smartcontract/service/native/shardmgmt"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontspace/ontology-bench/sharding-bench/src/config"
)

type ShardInitParam struct {
	Path string `json:"path"`
}

func ShardInit(sdk *sdk.OntologySdk, cfg *config.Config, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config from %s: %s", configFile, err)
	}

	param := &ShardInitParam{}
	if err := json.Unmarshal(data, param); err != nil {
		return fmt.Errorf("unmarshal shard init param: %s", err)
	}

	user, ok := getAccountByPassword(sdk, param.Path, cfg.Password)
	if !ok {
		return fmt.Errorf("get account failed")
	}

	method := shardmgmt.INIT_NAME
	contractAddress := utils.ShardMgmtContractAddress
	txHash := common.Uint256{}
	txHash, err = sdk.Native.InvokeNativeContract(cfg.GasPrice, cfg.GasLimit, user, 0,
		contractAddress, method, []interface{}{})
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("shard init txHash is :%s", txHash.ToHexString())
	return nil
}
