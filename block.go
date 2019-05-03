package main

import (
	"fmt"
	"encoding/hex"
	"crypto/sha256"
	"strings"
	"bytes"
	"io/ioutil"
	"encoding/gob"
	"log"
//	"os"
)

type Block struct {
	Version			int
	Prev_block	string	
	Merkle_root string	
	Target		  string	
	Nonce				uint
	Height			int
	Hash				string	
}

type Blockchain struct {
	Blocks[]	*Block
}

var merkle = "0000000000000000000000000000000000000000000000000000000000000000"
var target = "000000f000000000000000000000000000000000000000000000000000000000"

// Generate a new block
func newBlock(prev_block string, height int) *Block {
	block := &Block{1, prev_block, merkle, target, 0, height, ""}
	return block
}

func (block *Block) proof_of_work() {
	verHex := fmt.Sprintf("%x", block.Version)

	for i := 0; ; i++ {
		nonceHex := fmt.Sprintf("%x", i)
		all := fmt.Sprintf("%08v", verHex) + block.Prev_block + block.Merkle_root + block.Target + fmt.Sprintf("%08v", nonceHex)
		hash := sha256.Sum256([]byte(all))
		hashString := hex.EncodeToString(hash[:])

		if hashString < block.Target {
			block.Nonce = uint(i)
			block.Hash = hashString
			break
		}
	}
}

func calculateHash(block *Block) string {
	verHex := fmt.Sprintf("%x", block.Version)
	nonceHex := fmt.Sprintf("%x", block.Nonce)
	header := fmt.Sprintf("%08v", verHex) + block.Prev_block + block.Merkle_root + block.Target + fmt.Sprintf("%08v", nonceHex)
	hash := sha256.New()
	hash.Write([]byte(header))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed[:])
}


func isBlockValid(newBlock, oldBlock *Block) bool {
	if oldBlock.Height + 1 != newBlock.Height {
		return false
	}
	
	if oldBlock.Hash != newBlock.Prev_block {
		return false
	}
	
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	if strings.Compare(newBlock.Hash, newBlock.Target) > 0{
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

func (bc *Blockchain) saveToFile() {
	var content  bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(bc)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile("blockchain_storage", content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) loadFromFile() {
	fileContent, err := ioutil.ReadFile("blockchain_storage")
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(bc)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	/*
	var bc Blockchain

	if _, err := os.Stat("blockchain_storage"); os.IsNotExist(err) {
		genesisBlock := newBlock("0000000000000000000000000000000000000000000000000000000000000000", 1)
		genesisBlock.proof_of_work()

		bc.Blocks = append(bc.Blocks, genesisBlock)
		bc.saveToFile()
	} else {
		bc.loadFromFile()
	}

	fmt.Println("Block height: ", bc.Blocks[len(bc.Blocks) - 1].Height)

	for i := 0; ; i++ {
		prevBlock := bc.Blocks[len(bc.Blocks) - 1]
		var block *Block
		block = newBlock(prevBlock.Hash, prevBlock.Height + 1)
		block.proof_of_work()
		bc.Blocks = append(bc.Blocks, block)
		bc.saveToFile()
		fmt.Println("Mined!")
		fmt.Println("Block Hash: ", block.Hash)
		fmt.Println("Validity: ", isBlockValid(block, prevBlock))
		fmt.Println()
	}

	// fmt.Println("Blockchain Validity: ", isBlockchainValid(bc.Blocks))
	*/
	startServer()
}


