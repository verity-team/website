package ethereum

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"github.com/verity-team/dws/internal/common"
)

func fetchBlock(ctxt common.Context, bn uint64) ([]byte, error) {
	request := EthGetBlockByNumberRequest{
		JsonRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%x", bn), true},
		ID:      1,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	params := common.HTTPParams{
		URL:              ctxt.ETHRPCURL,
		RequestBody:      requestBytes,
		MaxWaitInSeconds: ctxt.MaxWaitInSeconds,
	}
	body, err := common.HTTPPost(params)
	if err != nil {
		return nil, err
	}

	log.Infof("fetched block %d", bn)

	err = writeBlockToFile(ctxt, bn, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBlock(ctxt common.Context, bn uint64) (*common.Block, error) {
	var (
		body      []byte
		err       error
		freshCopy bool
	)
	// try getting the finalized block from the cache
	body, err = getFinalizedBlockFromCache(ctxt, bn)
	if err != nil {
		log.Debugf("block %d not in cache", bn)
	} else {
		log.Infof("****** block %d served from cache", bn)
	}

	if body == nil {
		// not found in cache -- get it from the ethereum jsonrpc API provider
		body, err = fetchBlock(ctxt, bn)
		if err != nil {
			return nil, err
		}
		freshCopy = true
	}
	block, err := parseBlock(body)
	if err != nil {
		err = fmt.Errorf("failed to parse block #%d, %w", bn, err)
		log.Error(err)
		if freshCopy {
			return nil, err
		}
		// hmm .. maybe the cached copy was corrupted .. fetch again and retry
		body, err = fetchBlock(ctxt, bn)
		if err != nil {
			return nil, err
		}
		block, err = parseBlock(body)
		if err != nil {
			err = fmt.Errorf("failed to parse block #%d, %w", bn, err)
			return nil, err
		}
	}
	return block, nil
}

func parseBlock(body []byte) (*common.Block, error) {
	type Response struct {
		Block common.Block `json:"result"`
	}

	var resp Response
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	b := resp.Block
	log.Infof("parsed block %d, %s -- %d transactions", b.Number, b.Hash, len(b.Transactions))
	return &b, nil
}

func GetFinalizedBlock(ctxt common.Context, blockNumber uint64) (*common.FinalizedBlock, error) {
	request := EthGetBlockByNumberRequest{
		JsonRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%x", blockNumber), false},
		ID:      1,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	params := common.HTTPParams{
		URL:              ctxt.ETHRPCURL,
		RequestBody:      requestBytes,
		MaxWaitInSeconds: ctxt.MaxWaitInSeconds,
	}
	body, err := common.HTTPPost(params)
	if err != nil {
		return nil, err
	}

	err = writeBlockToFile(ctxt, blockNumber, body)
	if err != nil {
		return nil, err
	}

	fb, err := parseFinalizedBlock(body)
	if err != nil {
		err = fmt.Errorf("failed to parse finalized block #%d, %w", blockNumber, err)
		log.Error(err)
		return nil, err
	}
	return fb, nil
}

func parseFinalizedBlock(body []byte) (*common.FinalizedBlock, error) {
	type fblock struct {
		common.FinalizedBlock
		HexSeconds string `json:"timestamp"`
		HexNumber  string `json:"number"`
	}
	type Response struct {
		Block fblock `json:"result"`
	}

	var resp Response
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	seconds, err := common.HexStringToDecimal(resp.Block.HexSeconds)
	if err != nil {
		err = fmt.Errorf("failed to convert block timestamp, %w", err)
		return nil, err
	}
	number, err := common.HexStringToDecimal(resp.Block.HexNumber)
	if err != nil {
		err = fmt.Errorf("failed to convert block number, %w", err)
		return nil, err
	}
	ts := time.Unix(seconds.IntPart(), 0)
	result := common.FinalizedBlock{
		BaseFeePerGas: resp.Block.BaseFeePerGas,
		GasLimit:      resp.Block.GasLimit,
		GasUsed:       resp.Block.GasUsed,
		Hash:          resp.Block.Hash,
		Number:        uint64(number.IntPart()),
		ReceiptsRoot:  resp.Block.ReceiptsRoot,
		Size:          resp.Block.Size,
		StateRoot:     resp.Block.StateRoot,
		Timestamp:     ts.UTC(),
		Transactions:  resp.Block.Transactions,
	}
	return &result, nil
}
