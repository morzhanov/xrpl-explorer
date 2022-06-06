package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const (
	rippleTestAPIURL                = "s.altnet.rippletest.net:51233"
	rippleMainAPIURL                = "https://s1.ripple.com:51234/"
	reqTypeLedger                   = "ledgers"
	reqTypeTransactions             = "transactions"
	reqTypeAccountTransactions      = "acc-tx"
	reqBodyLedgersSubscription      = "{\"id\":\"Example watch for new validated ledgers\",\"command\": \"subscribe\",\"streams\": [\"ledger\"]}"
	reqBodyTransactionsSubscription = "{\"id\":\"Example watch for new validated transactions\",\"command\": \"subscribe\",\"streams\": [\"transactions\"]}"
)

func main() {
	reqType := os.Args[1]
	if reqType == reqTypeLedger {
		listenToLedgersUpdates()
	} else if reqType == reqTypeTransactions {
		listenToTransactionsUpdates()
	} else if reqType == reqTypeAccountTransactions {
		fetchAccountTransactions()
	} else {
		log.Fatal(errors.New(fmt.Sprintf("wrong subscription type provided, valid values are %s and %s", reqTypeLedger, reqTypeTransactions)))
	}
}

func listenToLedgersUpdates() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: rippleTestAPIURL, Path: "/"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(reqBodyLedgersSubscription))
	if err != nil {
		log.Fatal(err)
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	var result LedgerSubscriptionStartMessage
	if err = json.Unmarshal(message, &result); err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"STATUS = %s, TYPE = %s, HASH = %s, INDEX = %d, VALIDATED_LEDGERS: %s",
		result.Status,
		result.Type,
		result.Result.LedgerHash,
		result.Result.LedgerIndex,
		result.Result.ValidatedLedgers,
	)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			return
		default:
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg LedgerSubscriptionRes
			if err = json.Unmarshal(message, &msg); err != nil {
				log.Fatal(err)
			}
			log.Printf(
				"HASH = %s, INDEX = %d, VALIDATED_LEDGERS: %s",
				msg.LedgerHash,
				msg.LedgerIndex,
				msg.ValidatedLedgers,
			)
		}
	}
}

func listenToTransactionsUpdates() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: rippleTestAPIURL, Path: "/"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(reqBodyTransactionsSubscription))
	if err != nil {
		log.Fatal(err)
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	var result TransactionSubscriptionStartMessage
	if err = json.Unmarshal(message, &result); err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"STATUS = %s, TYPE = %s",
		result.Status,
		result.Type,
	)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			return
		default:
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg TransactionSubscriptionRes
			if err = json.Unmarshal(message, &msg); err != nil {
				log.Fatal(err)
			}
			log.Printf(
				"STATUS = %s, TYPE = %s, HASH = %s, INDEX = %d, TX_ACC: %s, TX_TYPE: %s, TX_HASH: %s",
				msg.Status,
				msg.Type,
				msg.LedgerHash,
				msg.LedgerIndex,
				msg.Transaction.Account,
				msg.Transaction.TransactionType,
				msg.Transaction.Hash,
			)
		}
	}
}

func fetchAccountTransactions() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	t := time.NewTicker(time.Second * 5)
	txs := make(map[string]bool, 0)
	txCount := 0

	req := FetchAccTransactionsReqBody{
		Method: "account_tx",
		Params: []FetchAccTransactionsParams{{
			Account:        "rLNaPoKeeBjZe2qs6x52yVPZpZ8td4dc6w",
			Binary:         false,
			Forward:        false,
			LedgerIndexMax: -1,
			LedgerIndexMin: -1,
		}},
	}
	reqBodyBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			return
		case <-t.C:
			reqBody := bytes.NewReader(reqBodyBytes)
			response, err := http.Post(rippleMainAPIURL, "application/json", reqBody)
			if err != nil {
				log.Fatal(err)
			}
			b, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			if response.StatusCode >= 400 {
				log.Fatal(string(b))
			}
			var res FetchAccTransactions
			if err = json.Unmarshal(b, &res); err != nil {
				log.Fatal(err)
			}
			for _, tx := range res.Result.Transactions {
				if _, ok := txs[tx.Tx.Hash]; ok {
					continue
				}
				txCount++
				txs[tx.Tx.Hash] = true
				log.Printf(
					"New Transaction: ACC = %s, HASH = %s, LEDGER_IDX = %s, FEE = %d, AMOUNT: %s, TX_TYPE: %s, TX_SIGN: %s",
					tx.Tx.Account,
					tx.Tx.Hash,
					tx.Tx.LedgerIndex,
					tx.Tx.Fee,
					tx.Tx.Amount,
					tx.Tx.TransactionType,
					tx.Tx.TxnSignature,
				)
				log.Println("TX Count = ", txCount)
			}
		}
	}
}
