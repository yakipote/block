package main

import "time"

type Chain struct {
}

type Block struct {
	Index        int
	TimeStamp    time.Time
	Transactions []Transaction
	Proof        int
	PreviousHash string
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
}

type CurrentTransaction struct {
}
type BlockChain struct {
	CurrentTransactions []Transaction
	Chain               []Block
}

func (b *BlockChain) NewBlock(proof int, previousHash string) Block {
	block := Block{
		Index:        len(b.Chain) + 1,
		TimeStamp:    time.Now(),
		Transactions: b.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	b.CurrentTransactions = []Transaction{}
	b.Chain = append(b.Chain, block)
	return block
}

func (b *BlockChain) NewTransaction(sender string, recipient string, amount int) int {
	b.CurrentTransactions = append(b.CurrentTransactions, Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	return b.Chain[len(b.Chain)-1].Index + 1

}

func Hash() {

}

func (b *BlockChain) LastBlock() {

}

func main() {

}
