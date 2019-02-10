/*
 * Copyright (C) 2019 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package shardmgmt

import (
	"bytes"
	"fmt"

	"github.com/ontio/ontology/common/config"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/global_params"
	"github.com/ontio/ontology/smartcontract/service/native/shardmgmt/states"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

/////////
//
// Shard management contract
//
//	. create shard
//	. config shard
//	. join shard
//	. activate shard
//
/////////

const (
	VERSION_CONTRACT_SHARD_MGMT = uint32(1)

	MAINCHAIN_SHARDID = uint64(0)

	// function names
	INIT_NAME           = "init"
	CREATE_SHARD_NAME   = "createShard"
	CONFIG_SHARD_NAME   = "configShard"
	JOIN_SHARD_NAME     = "joinShard"
	ACTIVATE_SHARD_NAME = "activateShard"

	// key prefix
	KEY_VERSION      = "version"
	KEY_GLOBAL_STATE = "globalState"
	KEY_SHARD_STATE  = "shardState"
)

func InitShardManagement() {
	native.Contracts[utils.ShardMgmtContractAddress] = RegisterShardMgmtContract
}

func RegisterShardMgmtContract(native *native.NativeService) {
	native.Register(INIT_NAME, ShardMgmtInit)
	native.Register(CREATE_SHARD_NAME, CreateShard)
	native.Register(CONFIG_SHARD_NAME, ConfigShard)
	native.Register(JOIN_SHARD_NAME, JoinShard)
	native.Register(ACTIVATE_SHARD_NAME, ActivateShard)
}

func ShardMgmtInit(native *native.NativeService) ([]byte, error) {
	// check if admin
	// get admin from database
	adminAddress, err := global_params.GetStorageRole(native,
		global_params.GenerateOperatorKey(utils.ParamContractAddress))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("getAdmin, get admin error: %v", err)
	}

	//check witness
	if err := utils.ValidateOwner(native, adminAddress); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("init shard mgmt, checkWitness error: %v", err)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress

	// check if shard-mgmt initialized
	ver, err := getVersion(native, contract)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("init shard mgmt, get version: %s", err)
	}
	if ver == 0 {
		// initialize shardmgmt version
		if err := setVersion(native, contract); err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("init shard mgmt version: %s", err)
		}

		// initialize shard mgmt
		globalState := &shardstates.ShardMgmtGlobalState{NextShardID: MAINCHAIN_SHARDID + 1}
		if err := setGlobalState(native, contract, globalState); err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("init shard mgmt global state: %s", err)
		}

		// initialize shard states
		mainShardState := &shardstates.ShardState{
			ShardID:             MAINCHAIN_SHARDID,
			ParentShardID:       config.DEFAULT_PARENT_SHARD_ID,
			GenesisParentHeight: uint64(native.Height),
			State:               shardstates.SHARD_STATE_ACTIVE,
		}
		if err := setShardState(native, contract, mainShardState); err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("init shard mgmt main shard state: %s", err)
		}
		return utils.BYTE_TRUE, nil
	}

	if ver < VERSION_CONTRACT_SHARD_MGMT {
		// make upgrade
		return utils.BYTE_FALSE, fmt.Errorf("upgrade TBD")
	} else if ver > VERSION_CONTRACT_SHARD_MGMT {
		return utils.BYTE_FALSE, fmt.Errorf("version downgrade from %d to %d", ver, VERSION_CONTRACT_SHARD_MGMT)
	}

	return utils.BYTE_TRUE, nil
}

func CreateShard(native *native.NativeService) ([]byte, error) {
	cp := new(CommonParam)
	if err := cp.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, invalid cmd param: %s", err)
	}

	params := new(CreateShardParam)
	if err := params.Deserialize(bytes.NewBuffer(cp.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, invalid param: %s", err)
	}
	if params.ParentShardID != MAINCHAIN_SHARDID {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, invalid parent shard: %d", params.ParentShardID)
	}
	if params.ParentShardID != native.ShardID {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, parent ShardID is not current shard")
	}

	if err := utils.ValidateOwner(native, params.Creator); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, invalid creator: %s", err)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	if ok, err := checkVersion(native, contract); !ok || err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, check version: %s", err)
	}

	globalState, err := getGlobalState(native, contract)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, get global state: %s", err)
	}

	shard := &shardstates.ShardState{
		ShardID:       globalState.NextShardID,
		ParentShardID: params.ParentShardID,
		Creator:       params.Creator,
		State:         shardstates.SHARD_STATE_CREATED,
		Peers:         make(map[string]*shardstates.PeerShardStakeInfo),
	}
	globalState.NextShardID += 1

	// TODO: SHARD CREATION FEE

	// update global state
	if err := setGlobalState(native, contract, globalState); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, update global state: %s", err)
	}
	// save shard
	if err := setShardState(native, contract, shard); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, set shard state: %s", err)
	}

	evt := &shardstates.CreateShardEvent{
		SourceShardID: native.ShardID,
		Height:        uint64(native.Height),
		NewShardID:    shard.ShardID,
	}
	if err := AddNotification(native, contract, evt); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, add notification: %s", err)
	}

	return utils.BYTE_TRUE, nil
}

func ConfigShard(native *native.NativeService) ([]byte, error) {
	cp := new(CommonParam)
	if err := cp.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, invalid cmd param: %s", err)
	}
	params := new(ConfigShardParam)
	if err := params.Deserialize(bytes.NewBuffer(cp.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, invalid param: %s", err)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	if ok, err := checkVersion(native, contract); !ok || err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, check version: %s", err)
	}

	shard, err := GetShardState(native, contract, params.ShardID)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, get shard: %s", err)
	}
	if shard == nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, get nil shard %d", params.ShardID)
	}

	if err := utils.ValidateOwner(native, shard.Creator); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, invalid configurator: %s", err)
	}
	if shard.ParentShardID != native.ShardID {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, not on parent shard")
	}

	if params.NetworkMin < 1 {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, invalid shard network size")
	}

	// TODO: support other stake
	if params.StakeAssetAddress.ToHexString() != utils.OntContractAddress.ToHexString() {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, only support ONT staking")
	}
	if params.GasAssetAddress.ToHexString() != utils.OngContractAddress.ToHexString() {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, only support ONG gas")
	}

	// TODO: validate input config
	shard.Config = &shardstates.ShardConfig{
		NetworkSize:       params.NetworkMin,
		StakeAssetAddress: params.StakeAssetAddress,
		GasAssetAddress:   params.GasAssetAddress,
		TestData:          params.ConfigTestData,
	}
	shard.State = shardstates.SHARD_STATE_CONFIGURED

	if err := setShardState(native, contract, shard); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, update shard state: %s", err)
	}

	return utils.BYTE_TRUE, nil
}

func JoinShard(native *native.NativeService) ([]byte, error) {
	cp := new(CommonParam)
	if err := cp.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, invalid cmd param: %s", err)
	}
	params := new(JoinShardParam)
	if err := params.Deserialize(bytes.NewBuffer(cp.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, invalid param: %s", err)
	}

	if err := utils.ValidateOwner(native, params.PeerOwner); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, invalid configurator: %s", err)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	if ok, err := checkVersion(native, contract); !ok || err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, check version: %s", err)
	}

	shard, err := GetShardState(native, contract, params.ShardID)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, get shard: %s", err)
	}
	if shard == nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, get nil shard %d", params.ShardID)
	}
	if shard.ParentShardID != native.ShardID {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, not on parent shard")
	}

	if _, present := shard.Peers[params.PeerPubKey]; present {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, peer already in shard")
	} else {
		peerStakeInfo := &shardstates.PeerShardStakeInfo{
			PeerOwner:   params.PeerOwner,
			PeerAddress: params.PeerAddress,
			StakeAmount: params.StakeAmount,
		}
		if shard.Peers == nil {
			shard.Peers = make(map[string]*shardstates.PeerShardStakeInfo)
		}
		shard.Peers[params.PeerPubKey] = peerStakeInfo
	}

	// TODO: asset transfer

	if err := setShardState(native, contract, shard); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("join shard, update shard state: %s", err)
	}

	return utils.BYTE_TRUE, nil
}

func ActivateShard(native *native.NativeService) ([]byte, error) {
	cp := new(CommonParam)
	if err := cp.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, invalid cmd param: %s", err)
	}
	params := new(ActivateShardParam)
	if err := params.Deserialize(bytes.NewBuffer(cp.Input)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, invalid param: %s", err)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	if ok, err := checkVersion(native, contract); !ok || err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, check version: %s", err)
	}

	shard, err := GetShardState(native, contract, params.ShardID)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, get shard: %s", err)
	}
	if shard == nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, get nil shard %d", params.ShardID)
	}

	if err := utils.ValidateOwner(native, shard.Creator); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, invalid configurator: %s", err)
	}
	if shard.State != shardstates.SHARD_STATE_CONFIGURED {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, invalid shard state: %d", shard.State)
	}
	if shard.ParentShardID != native.ShardID {
		return utils.BYTE_FALSE, fmt.Errorf("activate shard, not on parent shard")
	}

	// TODO: validate input config
	if uint32(len(shard.Peers)) < shard.Config.NetworkSize {
		return utils.BYTE_FALSE, fmt.Errorf("activae shard, not enough peer: %d vs %d",
			len(shard.Peers), shard.Config.NetworkSize)
	}

	shard.State = shardstates.SHARD_STATE_ACTIVE
	if err := setShardState(native, contract, shard); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("config shard, update shard state: %s", err)
	}

	evt := &shardstates.ShardActiveEvent{
		SourceShardID: native.ShardID,
		Height:        uint64(native.Height),
		ShardID:       shard.ShardID,
	}
	if err := AddNotification(native, contract, evt); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("create shard, add notification: %s", err)
	}

	return utils.BYTE_TRUE, nil
}
