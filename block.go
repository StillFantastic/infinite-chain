package main

import (
	"fmt"
	"encoding/hex"
	"crypto/sha256"
	"strings"
)

type Block struct {
	version			int
	prev_block	string	
	merkle_root string	
	target		  string	
	nonce				uint
	height			int
	hash				string	
}

type Blockchain struct {
	blocks[]	*Block
}

var merkle = "0000000000000000000000000000000000000000000000000000000000000000"
var target = "0f00000000000000000000000000000000000000000000000000000000000000"

// Generate a new block
func newBlock(prev_block string, height int) *Block {
	block := &Block{1, prev_block, merkle, target, 0, height, ""}
	return block
}

func (block *Block) proof_of_work() {
	verHex := fmt.Sprintf("%x", block.version)

	for i := 0; ; i++ {
		nonceHex := fmt.Sprintf("%x", i)
		all := fmt.Sprintf("%08v", verHex) + block.prev_block + block.merkle_root + block.target + fmt.Sprintf("%08v", nonceHex)
		hash := sha256.Sum256([]byte(all))
		hashString := hex.EncodeToString(hash[:])

		if hashString < block.target {
			block.nonce = uint(i)
			block.hash = hashString
			break
		}
	}
}

func calculateHash(block *Block) string {
	verHex := fmt.Sprintf("%x", block.version)
	nonceHex := fmt.Sprintf("%x", block.nonce)
	header := fmt.Sprintf("%08v", verHex) + block.prev_block + block.merkle_root + block.target + fmt.Sprintf("%08v", nonceHex)
	hash := sha256.New()
	hash.Write([]byte(header))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed[:])
}


func isBlockValid(newBlock, oldBlock *Block) bool {
	if oldBlock.height + 1 != newBlock.height {
		return false
	}
	
	if oldBlock.hash != newBlock.prev_block {
		return false
	}
	
	if calculateHash(newBlock) != newBlock.hash {
		return false
	}

	if strings.Compare(newBlock.hash, newBlock.target) > 0{
		return false
	}

	return true
}


func isBlockchainValid(blockArr []*Block) bool {
	for i, _ := range blockArr {
		if i == 0 {
			continue
		}
		if(isBlockValid(blockArr[i], blockArr[i-1]) == false){
			return false
		}
	}

	return true
}

func main() {
	var bc Blockchain

	genesisBlock := newBlock("", 1)
	genesisBlock.proof_of_work()

	bc.blocks = append(bc.blocks, genesisBlock)

	for i := 0; i < 10; i++ {
		prevBlock := bc.blocks[len(bc.blocks) - 1]
		var block *Block
		block = newBlock(prevBlock.hash, prevBlock.height + 1)
		block.proof_of_work()
		bc.blocks = append(bc.blocks, block)
		fmt.Println("Mined!")
		fmt.Println("Block Hash: ", block.hash)
		fmt.Println("Validity: ", isBlockValid(block, prevBlock))
		fmt.Println()
	}

	fmt.Println("Blockchain Validity: ", isBlockchainValid(bc.blocks))
}


