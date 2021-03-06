package addrindex

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// Blocks represents the cached /blocks response
type Blocks struct {
	Blocks []GetBlocksResponse `json:"blocks"`
	Length int                 `json:"length"`
	sync.Mutex
}

// GetBlocksResponse pulls new values for the blocks
func (as *AddrServer) GetBlocksResponse(limit int64) []byte {
	now := time.Now()
	blocks, err := as.GetBlockHashes(int(now.Unix()), int(now.Add(-24*time.Hour).Unix()))
	if err != nil {
		log.Println("Failed fetching block hashes")
		return []byte("")
	}
	var toQuery []string
	for i := len(blocks.Result) - 1; i >= 0; i-- {
		toQuery = append(toQuery, blocks.Result[i])
	}
	toQuery = toQuery[:limit]

	var out []GetBlocksResponse
	for _, blh := range toQuery {
		blockHash, err := chainhash.NewHashFromStr(blh)
		if err != nil {
			log.Println("Failed creating chainhash from block data")
			continue
		}
		block, err := as.Client.GetBlockVerbose(blockHash)
		if err != nil {
			log.Println("Failed fetching block data")
			return []byte("")
		}
		out = append(out, newGetBlockResponse(block))
	}
	ret := &Blocks{
		Length: int(limit),
		Blocks: out,
	}
	return ret.JSON()
}

// JSON returns the JSON representation of Blocks
func (b *Blocks) JSON() []byte {
	o, _ := json.Marshal(b)
	return o
}

func newGetBlockResponse(blk *btcjson.GetBlockVerboseResult) GetBlocksResponse {
	return GetBlocksResponse{
		Height:   blk.Height,
		Size:     blk.Weight,
		Hash:     blk.Hash,
		Time:     blk.Time,
		Txlength: len(blk.Tx),
	}
}

// GetBlocksResponse formats the response for the GetBlocks route
type GetBlocksResponse struct {
	Height   int64    `json:"height"`
	Size     int32    `json:"size"`
	Hash     string   `json:"hash"`
	Time     int64    `json:"time"`
	Txlength int      `json:"txlength"`
	PoolInfo PoolInfo `json:"poolInfo"`
}

// PoolInfo represents the mining pool information
type PoolInfo struct {
	PoolName string `json:"poolName,omitempty"`
	URL      string `json:"url,omitempty"`
}
