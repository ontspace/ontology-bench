package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/shardmgmt"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontspace/ontology-bench/sharding-bench/src/config"
)

type ShardActivateParam struct {
	Path    string `json:"path"`
	ShardID uint64 `json:"shard_id"`
}

func ShardActivate(sdk *sdk.OntologySdk, cfg *config.Config, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config from %s: %s", configFile, err)
	}

	param := &ShardActivateParam{}
	if err := json.Unmarshal(data, param); err != nil {
		return fmt.Errorf("unmarshal shard activate param: %s", err)
	}

	user, ok := getAccountByPassword(sdk, param.Path, cfg.Password)
	if !ok {
		return fmt.Errorf("get account failed")
	}

	tShardId, _ := types.NewShardID(param.ShardID)
	activateParam := &shardmgmt.ActivateShardParam{
		ShardID: tShardId,
	}

	buf := new(bytes.Buffer)
	if err := activateParam.Serialize(buf); err != nil {
		return fmt.Errorf("failed to ser activate shard param: %s", err)
	}

	method := shardmgmt.ACTIVATE_SHARD_NAME
	contractAddress := utils.ShardMgmtContractAddress
	txHash, err := sdk.Native.InvokeNativeContract(cfg.GasPrice, cfg.GasLimit, user, 0,
		contractAddress, method, []interface{}{buf.Bytes()})
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("activate shard txHash is :%s", txHash.ToHexString())
	return nil
}
