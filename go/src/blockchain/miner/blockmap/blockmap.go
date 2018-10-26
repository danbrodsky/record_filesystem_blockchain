package blockmap

import (
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "math/rand"
)

type BlockMap struct {
    TailBlock Block
    GenesisBlock Block
    Map map[string]Block
}

type BM interface {
    Insert(block Block) (err error)
    GetMap() (map[string]Block)
    GetLongestChain() ([]Block)
}


type Block struct{
    PrevHash string
    Ops []string
    Nonce string
    MinerId string
    Depth int
}

func NewBlockMap(genesisBlock Block, genesisHash string) (blockmap BlockMap) {
    blockmap = BlockMap{}
    blockmap.TailBlock = genesisBlock
    blockmap.GenesisBlock = genesisBlock
    blockmap.Map = make(map[string]Block)
    blockmap.Map[genesisHash] = genesisBlock
    return blockmap
}


type PrevHashDoesNotExistError string

func (e PrevHashDoesNotExistError) Error() string {
    return fmt.Sprintf("Block-Map: Error hash does not exist in map [%s]", string(e))
}

func getHash(block Block) string{
     h := md5.New()
     h.Write([]byte(fmt.Sprintf("%v", block)))
     return hex.EncodeToString(h.Sum(nil))
}

func (bm *BlockMap) updateLongest(block Block) {
    if block.Depth == bm.TailBlock.Depth {
        if rand.Intn(2) == 1 {
            bm.TailBlock = block
        }
    }
    if block.Depth > bm.TailBlock.Depth {
        bm.TailBlock = block
    }
}

func (bm *BlockMap) Insert(block Block) (err error){
	// TODO verify block hash end of 0s
    if _, ok := bm.Map[block.PrevHash]; !ok {
	    fmt.Println("block to add:", bm.TailBlock)
        bm.Map[getHash(block)] = block
	    bm.updateLongest(block)
	    fmt.Println("tail:", bm.TailBlock)
	    return nil
    } else {
	return PrevHashDoesNotExistError(block.PrevHash)
    }
}

func (bm *BlockMap) GetMap() (map[string]Block){
	return bm.Map
}

func (bm *BlockMap) GetLongestChain() ([]Block){
    var blockChain []Block
    var currBlock = bm.TailBlock
    for bm.Map[currBlock.PrevHash].PrevHash == bm.GenesisBlock.PrevHash {
        blockChain = append(blockChain, bm.Map[currBlock.PrevHash])
    }
    blockChain = append(blockChain, bm.GenesisBlock)
    return blockChain
}


