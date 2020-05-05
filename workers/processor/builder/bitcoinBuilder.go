package builder

import (
	"../../../api/blockatlas"
	m "../../../models"
	"encoding/json"
	"log"
	"math/big"
)

type InOut struct {
	Amount  *big.Int
	Address string
	Spent   bool
}

func BuiltBitcoinTransactions(orders []m.Order, addresses []m.Address, transactions []m.Transaction) {
	ins := map[string]map[uint]InOut{}
	mAddresses := m.Addresses(addresses)

	for _, t := range transactions {
		if t.Status != m.TransactionStatusConfirmed {
			continue
		}
		var tx = new(blockatlas.Tx)
		json.Unmarshal([]byte(t.TxData), tx)
		for outputIdx, curOutput := range tx.Outputs {
			if mAddresses.FindAddress(curOutput.Address) == nil {
				continue
			}
			if insTx, exists := ins[tx.ID]; exists {
				if _, exists := insTx[curOutput.N]; exists {
					log.Println("Double In", tx.ID, outputIdx)
				} else {
					Amount, _ := new(big.Int).SetString(string(curOutput.Value), 10)
					ins[tx.ID][curOutput.N] = InOut{Amount: Amount, Address: curOutput.Address, Spent: false}
				}
			} else {
				ins[tx.ID] = map[uint]InOut{}
				Amount, _ := new(big.Int).SetString(string(curOutput.Value), 10)
				ins[tx.ID][curOutput.N] = InOut{Amount: Amount, Address: curOutput.Address, Spent: false}
			}
		}
	}

	for _, t := range transactions {
		if t.Status != m.TransactionStatusConfirmed {
			continue
		}
		var tx = new(blockatlas.Tx)
		json.Unmarshal([]byte(t.TxData), tx)
		for _, curInput := range tx.Inputs {
			if mAddresses.FindAddress(curInput.Address) == nil {
				continue
			}
			if insTx, exists := ins[curInput.TxId]; exists {
				if _, exists := insTx[curInput.N]; exists {
					curValue, _ := new(big.Int).SetString(string(curInput.Value), 10)
					ins[curInput.TxId][curInput.N] = InOut{Amount: curValue, Address: curInput.Address, Spent: true}
				}
			}
		}
	}
	log.Println(ins)

}
