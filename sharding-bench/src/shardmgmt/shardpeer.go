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
	"github.com/ontspace/ontology-bench/config"
)

type ShardPeerJoinParam struct {
	Path        string `json:"path"`
	ShardID     uint64 `json:"shard_id"`
	PeerAddress string `json:"peer_address"`
	PeerPubKey  string `json:"peer_pub_key"`
	StakeAmount uint64 `json:"stake_amount"`
}

func ShardPeerJoin(sdk *sdk.OntologySdk, cfg *config.Config, configFile string) error {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config from %s: %s", configFile, err)
	}

	param := &ShardPeerJoinParam{}
	if err := json.Unmarshal(data, param); err != nil {
		return fmt.Errorf("unmarshal shard peer join param: %s", err)
	}

	user, ok := getAccountByPassword(sdk, param.Path, cfg.Password)
	if !ok {
		return fmt.Errorf("get account failed")
	}

	tShardId, _ := types.NewShardID(param.ShardID)
	joinParam := &shardmgmt.JoinShardParam{
		ShardID:     tShardId,
		PeerOwner:   user.Address,
		PeerAddress: param.PeerAddress,
		PeerPubKey:  param.PeerPubKey,
		StakeAmount: param.StakeAmount,
	}

	buf := new(bytes.Buffer)
	if err := joinParam.Serialize(buf); err != nil {
		return fmt.Errorf("failed to ser join shard param: %s", err)
	}

	method := shardmgmt.JOIN_SHARD_NAME
	contractAddress := utils.ShardMgmtContractAddress
	txHash, err := sdk.Native.InvokeNativeContract(cfg.GasPrice, cfg.GasLimit, user, 0,
		contractAddress, method, []interface{}{buf.Bytes()})
	if err != nil {
		return fmt.Errorf("invokeNativeContract error :", err)
	}
	log.Infof("join shard txHash is :%s", txHash.ToHexString())
	return nil
}
