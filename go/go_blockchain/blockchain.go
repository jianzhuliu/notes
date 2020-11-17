package go_blockchain

import (
	"sync"
)

//区块链
type BlockChain struct {
	blocks []*Block
	mu     sync.RWMutex //加锁控制
}

//创建区块链对象
func NewBlockChain() *BlockChain {
	block := NewBlock()
	blockChain := new(BlockChain)
	blockChain.blocks = append(blockChain.blocks, block)
	return blockChain
}

//根据数据，新增区块
func (bc *BlockChain) AddBlock(data string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	preBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := preBlock.CreateBlock(data)

	bc.blocks = append(bc.blocks, newBlock)
}

//获取区块链表
func (bc *BlockChain) GetBlocks() []*Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.blocks
}
