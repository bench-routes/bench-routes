package file

import "encoding/json"

// Block use case block for the TSDB chain
type Block struct {
	Datapoint      string `json:"datapoint"`       // complex data would be decoded by using a blockSeparator
	NormalizedTime int64  `json:"normalized-time"` // based on time.UnixNano()
	Type           string `json:"type"`            // would be used to decide the marshalling struct
	Timestamp      string `json:"timestamp"`
}

// NewBlock creates and returns a new block with the specified type.
func NewBlock(blockType, value string) Block {
	return Block{
		Timestamp:      GetTimeStampCalc(),
		NormalizedTime: GetNormalizedTimeCalc(),
		Datapoint:      value,
		Type:           blockType,
	}
}

// Encode decodes the structure and marshals into a string
func (b Block) Encode() []byte {
	bbyte, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return bbyte
}

// GetType returns the type of the block
func (b Block) GetType() string {
	return b.Type
}

// GetDatapointEnc returns the data point to the caller.
// The encoded refers to the combined _(containing *|*)_ values in the string
// form.
func (b Block) GetDatapointEnc() string {
	return b.Datapoint
}

// GetNormalizedTimeStringified returns the normalized time of the block.
func (b Block) GetNormalizedTimeStringified() string {
	return string(rune(b.NormalizedTime))
}

// GetNormalizedTime returns the normalized time of the block.
func (b Block) GetNormalizedTime() int64 {
	return b.NormalizedTime
}

// GetTimeStamp returns the timestamp of the block.
func (b Block) GetTimeStamp() string {
	return b.Timestamp
}

// mergeBlocksSlice merges the slices of two blocks into one.
func mergeBlocksSlice(oldSlice, newSlice []Block) []Block {
	oldSlice = append(oldSlice, newSlice...)
	return oldSlice
}
