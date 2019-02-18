package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danitello/go-blockchain/common/errutil"
)

const walletFile = "./tmp/wallets.dat"

/*Wallets keeps track of all current Wallet structs
@param Wallets - map of addresses to Wallet structs
*/
type Wallets struct {
	Wallets map[string]*Wallet
}

/*InitWallets makes a new Wallets struct and loads it with previous Wallets data if possible
@return the new Wallet
@return any error
*/
func InitWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

/*CreateWallet makes a new wallet and adds it to the Wallets
@return the new wallet address
*/
func (ws *Wallets) CreateWallet() string {
	wallet := InitWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

/*GetAddresses retrieves all of the address from the Wallets
@return a slice of address strings
*/
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

/*GetWallet retrieves a specific wallet by address
@param address - the address of the desired wallet
*/
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

/*LoadFromFile loads Wallets data from disk
@return any error
*/
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	data, err := ioutil.ReadFile(walletFile)
	errutil.HandleErr(err)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&wallets)
	errutil.HandleErr(err)

	ws.Wallets = wallets.Wallets

	return nil
}

/*SaveToFile writes the Wallets data to disk */
func (ws *Wallets) SaveToFile() {
	var data bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(ws)
	errutil.HandleErr(err)

	err = ioutil.WriteFile(walletFile, data.Bytes(), 0644)
	errutil.HandleErr(err)
}