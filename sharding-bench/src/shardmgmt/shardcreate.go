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

type ShardCreateParam struct {
	Path          string `json:"path"`
	ParentShardID uint64 `json:"parent_shard_id"`
}

func ShardCreate(sdk *sdk.OntologySdk, cfg *config.Config, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config from %s: %s", configFile, err)
	}

	param := &ShardCreateParam{}
	if err := json.Unmarshal(data, param); err != nil {
		return fmt.Errorf("unmarshal shard create param: %s", err)
	}

	user, ok := getAccountByPassword(sdk, param.Path, cfg.Password)
	if !ok {
		return fmt.Errorf("get account failed")
	}

	tShardId, _ := types.NewShardID(param.ParentShardID)
	createParam := &shardmgmt.CreateShardParam{
		ParentShardID: tShardId,
		Creator:       user.Address,
	}

	buf := new(bytes.Buffer)
	if err := createParam.Serialize(buf); err != nil {
		return fmt.Errorf("failed to ser createshard param: %s", err)
	}
	method := shardmgmt.CREATE_SHARD_NAME
	contractAddress := utils.ShardMgmtContractAddress
	txHash, err := sdk.Native.InvokeNativeContract(cfg.GasPrice, cfg.GasLimit, user, 0,
		contractAddress, method, []interface{}{buf.Bytes()})
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("shard create txHash is :%s", txHash.ToHexString())
	return nil
}

func getAccountByPassword(sdk *sdk.OntologySdk, path string, pwd string) (*sdk.Account, bool) {
	wallet, err := sdk.OpenWallet(path)
	if err != nil {
		log.Errorf("open wallet error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount([]byte(pwd))
	if err != nil {
		log.Errorf("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}
