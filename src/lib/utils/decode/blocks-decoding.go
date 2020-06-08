package decode

import (
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

// Decode function checks for different kinds of modules and redirects
// the same to the respective functions to get the decoded value which
// is to be passed to the front end
func (bd *BlockDecodingBR) Decode(block tsdb.Block) interface{} {
	switch bd.Type {
	case "sys":
		return systemDecode(block.Datapoint)
	case "ping":
		return pingDecode(block.Datapoint)
	case "jitter":
		return jitterDecode(block.Datapoint)
	case "journal":
		return journalDecode(block.Datapoint)
	case "flood-ping":
		return floodPingDecode(block.Datapoint)
	case "req-res":
		return monitorDecode(block.Datapoint)
	}
	return nil
}
