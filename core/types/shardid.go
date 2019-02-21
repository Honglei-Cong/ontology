package types

import (
	"errors"
	"math"
)

type ShardID uint64

const MAX_SHARD_LEVEL = 4
const MAX_CHILD_SHARDS = math.MaxUint16 - 1

func (self ShardID) Level() int {
	l1 := self & math.MaxUint16
	l2 := (self >> 16) & math.MaxUint16
	l3 := (self >> 32) & math.MaxUint16
	l4 := (self >> 48) & math.MaxUint16
	lvs := [4]ShardID{l4, l3, l2, l1}

	level := MAX_SHARD_LEVEL
	for _, l := range lvs {
		if l != 0 {
			return level
		}
		level -= 1
	}
	return level
}

func (self ShardID) GenSubShardID(index uint16) (ShardID, error) {
	if index == 0 {
		return 0, errors.New("wrong child shard index")
	}
	level := self.Level()
	if level == MAX_SHARD_LEVEL {
		return 0, errors.New("can not generate sub shard id, max level reached")
	}

	subId := uint64(index) << (16 * uint64(level))

	return ShardID(uint64(self) + subId), nil
}

func ShardIDFromLevels(l1, l2, l3, l4 uint16) ShardID {
	id := uint64(l1)
	id += uint64(l2) << 16
	id += uint64(l3) << 32
	id += uint64(l4) << 48
	return ShardID(id)
}
