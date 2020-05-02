package utils

import (
	"github.com/zairza-cetb/bench-routes/src/metrics/system"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// BlockDecodingBR implements the decoding of tsdb blocks into the respective types.
type BlockDecodingBR struct {
	Type string
}

// NewBlockDecoding returns the new BlockDecodingBR type.
func NewBlockDecoding(Type string) *BlockDecodingBR {
	return &BlockDecodingBR{
		Type: Type,
	}
}

func (bd *BlockDecodingBR) Decode(block tsdb.Block) interface{} {
	switch bd.Type {
	case "sys":
		return system.Decode(block.Datapoint)
	}
	return nil
}
