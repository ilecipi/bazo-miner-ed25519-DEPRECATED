package miner

import (
	"errors"
	"github.com/bazo-blockchain/bazo-miner/protocol"
	"github.com/bazo-blockchain/bazo-miner/storage"
	"github.com/bazo-blockchain/bazo-miner/p2p"
	"fmt"
)

type SlashingProof struct {
	ConflictingBlockHash1 [32]byte
	ConflictingBlockHash2 [32]byte
}

//find a proof where a validator votes on two different chains within the slashing window
func seekSlashingProof(block *protocol.Block) error {
	//check if block is being added to your chain
	lastClosedBlock := storage.ReadLastClosedBlock()
	if lastClosedBlock == nil {
		return errors.New("Latest block not found.")
	}

	//fmt.Printf("lastClosedBlock, block\nlastClosedBlock:\n%v\n, block:\n%v\n", lastClosedBlock, block)

	//when the block is added ontop of your chain then there is no slashing needed
	if lastClosedBlock.Hash == block.PrevHash {
		//fmt.Printf("lastClosedBlock.Hash == block.PrevHash - %x == %x", lastClosedBlock.Hash, block.PrevHash)
		return nil
	} else {
		//get the latest blocks and check if there is proof for multivoting within the slashing window
		prevBlocks := storage.ReadAllClosedBlocks()

		if prevBlocks == nil {
			return nil
		}

		for _, prevBlock := range prevBlocks {
			if prevBlock == nil {
				continue
			}
			if IsInSameChain(prevBlock, block) {
				return nil
			}
			if prevBlock.Beneficiary == block.Beneficiary &&
				(uint64(prevBlock.Height) < uint64(block.Height)+activeParameters.Slashing_window_size ||
					uint64(block.Height) < uint64(prevBlock.Height)+activeParameters.Slashing_window_size) {
						fmt.Printf("block.Beneficiary:\n%v\n", block)
						fmt.Printf("block.Beneficiary:\n%v\n", block.Beneficiary)
						fmt.Printf("prevblock.Beneficiary:\n%v\n", prevBlock)
						fmt.Printf("slashingDict:\n%v\n", slashingDict)
				slashingDict[block.Beneficiary] = SlashingProof{ConflictingBlockHash1: block.Hash, ConflictingBlockHash2: prevBlock.Hash}
			}
		}
	}
	return nil
}

//Check if two blocks are part of the same chain or if they appear in two competing chains
func IsInSameChain(b1, b2 *protocol.Block) bool {
	var higherBlock *protocol.Block
	var lowerBlock *protocol.Block
	if b1.Height == b2.Height {
		return false
	}
	if b1.Height > b2.Height {
		higherBlock = b1
		lowerBlock = b2
	} else {
		higherBlock = b2
		lowerBlock = b1
	}

	// TODO: prove that this is correct
	if higherBlock.Height > 0 && higherBlock.NrConsolidationTx > 0{
		txHash := higherBlock.ConsolidationTxData[0]
		tx := getTransaction(p2p.CONSOLIDATIONTX_REQ, txHash)
		consolidationTx := tx.(*protocol.ConsolidationTx)
		if consolidationTx.PreviousConsHash == lowerBlock.Hash {
			return true
		}
	}

	for higherBlock != nil && higherBlock.Height > 0 {
		higherBlock = storage.ReadClosedBlock(higherBlock.PrevHash)

		if higherBlock!= nil && higherBlock.Hash == lowerBlock.Hash {
			return true
		}
	}
	return false
}
