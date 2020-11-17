package go_blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

//区块结构体
type Block struct {
	Index        int    `json:"index"`        //区块编号
	Data         string `json:"data"`         //区块保存的数据
	Hash         string `json:"hash"`         //当前对应的hash值,32位
	PreBlockHash string `json:"preBlockHash"` //上一个区块对应的hash值
	Timestamp    int64  `json:"timestamp"`    //创建时间戳
}

//创建初始区块
func NewBlock() *Block {
	return &Block{
		Index:     -1,
		Hash:      "",
		Timestamp: time.Now().Unix(),
	}
}

//创建新区块
func (b *Block) CreateBlock(data string) *Block {
	block := new(Block)
	block.Index = b.Index + 1
	block.Data = data
	block.PreBlockHash = b.Hash
	block.Timestamp = time.Now().Unix()
	block.Hash = genHash(block)

	return block
}

//生成hash
func genHash(b *Block) string {
	sumData := strconv.Itoa(b.Index) + b.Data + strconv.FormatInt(b.Timestamp, 10) + b.PreBlockHash
	hashBytes := sha256.Sum256([]byte(sumData))
	return hex.EncodeToString(hashBytes[:])
}
