package querier

import (
	"encoding/json"
	"math"
	"testing"
)

func Test_Querier_NULL_Range(t *testing.T) {
	// tests the case where the range is not provided. This should return all the blocks.
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 334 {
		t.Errorf("NULL_Range: decoded blocks are not in required numbers, %d, should actually be 334", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_complete_start_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588420763115213398, 1588231243752392414)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 334 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 334", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_infinite_start_and_finite_last_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(int64(math.MaxInt64), 1588231243752392414)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 334 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 334", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_infinite_start_and_finite_middle_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(int64(math.MaxInt64), 1588320025023992340)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 225 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 225", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_finite_first_start_and_infinite_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588420763115213398, int64(math.MinInt64))
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 334 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 334", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_finite_middle_start_and_infinite_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588326076041632077, int64(math.MinInt64))
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 234 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 234", len(decodedBlocks))
	}
}

func Test_Querier_Range_within_finite_middle_start_and_finite_middle_end(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588326076041632077, 1588319816131081362)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 140 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 140", len(decodedBlocks))
	}
}

func Test_Querier_Range_two_block_test(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588319829132430602, 1588319816131081362)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 2 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 2", len(decodedBlocks))
	}
}

func Test_Querier_Range_single_block_test(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588319829132430601, 1588319816131081362)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 1 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 1", len(decodedBlocks))
	}
}

func Test_Querier_Range_Invalid_test(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588319816131081362, 1588319829132430602)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	_, ok := resp.Value.([]interface{})
	if ok {
		t.Errorf("Range_Invalid_test: type should be string, but received %t", resp.Value)
	}
}

func Test_Querier_Range_on_point(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588319816131081362, 1588319816131081362)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 1 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 1", len(decodedBlocks))
	}
}

func Test_Querier_Range_single_block_near_line(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1588319816131081363, 1588319816131081363)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	decodedBlocks := resp.Value.([]interface{})
	if len(decodedBlocks) != 1 {
		t.Errorf("Range_within_complete_start_end: decoded blocks are not in required numbers, %d, should actually be 1", len(decodedBlocks))
	}
}

func Test_Querier_Range_no_resultant_block(t *testing.T) {
	var resp QueryResponse
	q := New("test_sample_blocks.json", "", TypeRange)
	query := q.QueryBuilder()
	query.SetRange(1688420763115213398, 1658420763115213398)
	result := query.Exec()
	if err := json.Unmarshal(result, &resp); err != nil {
		panic(err)
	}
	_, ok := resp.Value.([]interface{})
	if ok {
		t.Errorf("Range_Invalid_test: type should be string, but received %t", resp.Value)
	}
}
