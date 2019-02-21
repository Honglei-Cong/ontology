package types

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestShardID_Level(t *testing.T) {
	assert.Equal(t, ShardIDFromLevels(0, 0, 0, 0).Level(), 0)
	assert.Equal(t, ShardIDFromLevels(1, 0, 0, 0).Level(), 1)
	assert.Equal(t, ShardIDFromLevels(1, 1, 0, 0).Level(), 2)
	assert.Equal(t, ShardIDFromLevels(1, 0, 1, 0).Level(), 3)
	assert.Equal(t, ShardIDFromLevels(1, 0, 0, 4).Level(), 4)
}

func TestShardID_GenSubShardID(t *testing.T) {
	_, err := ShardIDFromLevels(1, 0, 0, 0).GenSubShardID(0)
	assert.NotNil(t, err)
	_, err = ShardIDFromLevels(1, 0, 0, 0).GenSubShardID(math.MaxUint16)
	assert.Nil(t, err)
	id, _ := ShardIDFromLevels(1, 0, 0, 0).GenSubShardID(234)
	assert.Equal(t, id, ShardIDFromLevels(1, 234, 0, 0))

	id, _ = ShardIDFromLevels(1, 0, 2, 0).GenSubShardID(234)
	assert.Equal(t, id, ShardIDFromLevels(1, 0, 2, 234))
}
