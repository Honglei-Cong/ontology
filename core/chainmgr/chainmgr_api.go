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

package chainmgr

import (
	"math"

	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/common/log"
)

func IsRootShard(shardId uint64) bool {
	return shardId == 0
}

func GetChainManager() *ChainManager {
	return defaultChainManager
}

func GetAccount() *account.Account {
	chainmgr := GetChainManager()
	return chainmgr.account
}

func GetShardID() uint64 {
	return GetChainManager().shardID
}

func GetParentShardID() uint64 {
	chainmgr := GetChainManager()
	return chainmgr.parentShardID
}

func GetParentBlockHeight() uint64 {
	chainmgr := GetChainManager()
	chainmgr.lock.RLock()
	defer chainmgr.lock.RUnlock()

	if IsRootShard(chainmgr.shardID) {
		return 0
	}

	m := chainmgr.blockPool.Shards[chainmgr.parentShardID]
	if m == nil {
		return math.MaxUint64
	}

	h := uint64(0)
	for _, blk := range m {
		if blk.Height > h {
			h = blk.Height
		}
	}

	return h
}

func GetParentBlockHeader(height uint64) *types.Header {
	chainmgr := GetChainManager()
	chainmgr.lock.RLock()
	defer chainmgr.lock.RUnlock()
	if IsRootShard(chainmgr.shardID) {
		return nil
	}

	m := chainmgr.blockPool.Shards[chainmgr.parentShardID]
	if m == nil {
		return nil
	}
	if blk, present := m[height]; present && blk != nil {
		return blk.Header.Header
	}

	return nil
}

func GetShardTxsByParentHeight(start, end uint64) map[uint64][]*types.Transaction {
	chainmgr := GetChainManager()

	chainmgr.lock.RLock()
	defer chainmgr.lock.RUnlock()
	if IsRootShard(chainmgr.shardID) {
		return nil
	}

	parentShard := chainmgr.parentShardID
	m := chainmgr.blockPool.Shards[parentShard]
	if m == nil {
		return nil
	}
	shardTxs := make(map[uint64][]*types.Transaction)
	for ; start < end+1; start++ {
		if blk, present := m[start]; present && blk != nil {
			if shardTx, present := blk.ShardTxs[chainmgr.shardID]; present && shardTx != nil && shardTx.Tx != nil {
				if shardTxs[parentShard] == nil {
					shardTxs[parentShard] = make([]*types.Transaction, 0)
				}
				shardTxs[parentShard] = append(shardTxs[parentShard], shardTx.Tx)
				log.Infof(">>>> shard %d got remote Tx from parent %d, height: %d",
					chainmgr.shardID, parentShard, start)
			}
		}
	}

	return shardTxs
}
