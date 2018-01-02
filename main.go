package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type Block struct {
	Index        int           `json:"index"`
	TimeStamp    time.Time     `json:"timeStamp"`
	Transactions []Transaction `json:transactions`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previousHash"`
}

type Transaction struct {
	Sender    string `json:'sender'`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

type BlockChain struct {
	CurrentTransactions []Transaction `json:"transactions"`
	Chain               []Block       `json:"chain"`
}

func NewBlockChain() *BlockChain {
	blockChain := new(BlockChain)

	// create the genesis block
	blockChain.NewBlock(100, "1")
	return blockChain
}

func ProofOfWork(lastProof int) int {
	proof := 0
	for {
		if ValidProof(lastProof, proof) {
			return proof
		}
		proof += 1
	}
}

func ValidProof(lastProof int, proof int) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(proof)
	guessHash := sha256.Sum256([]byte(guess))
	guessHashStr := hex.EncodeToString(guessHash[:])
	return guessHashStr[:4] == "0000"
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

func Hash(block Block) string {
	bytes, err := json.Marshal(block)
	if err != nil {
	}
	hash := sha256.Sum256(bytes)
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func (b *BlockChain) LastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

var blockChain *BlockChain
var NodeIdentifier string

func main() {
	e := echo.New()
	blockChain = NewBlockChain()
	// nodeアドレス
	NodeIdentifier = uuid.New().String()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/transactions/new", newTransactions)
	e.GET("/mine", mine)
	e.GET("/chain", chain)
	e.Logger.Fatal(e.Start(":1323"))
}

type TransactionRequest struct {
	Sender    string
	Recipient string
	Amount    int
}

// api
func mine(c echo.Context) error {
	lastBlock := blockChain.LastBlock()
	lastProof := lastBlock.Proof
	proof := ProofOfWork(lastProof)

	blockChain.NewTransaction(
		"0",
		NodeIdentifier,
		1,
	)

	previousHash := Hash(lastBlock)
	block := blockChain.NewBlock(proof, previousHash)

	c.JSON(http.StatusOK, block)
	return nil
}

func newTransactions(c echo.Context) error {
	transactionRequest := new(TransactionRequest)
	c.Bind(transactionRequest)
	index := blockChain.NewTransaction(transactionRequest.Sender, transactionRequest.Recipient, transactionRequest.Amount)
	c.String(http.StatusOK, "Transaction will be added to Block "+strconv.Itoa(index))
	return nil

}

func chain(c echo.Context) error {
	c.JSON(http.StatusOK, blockChain.Chain)
	return nil
}
