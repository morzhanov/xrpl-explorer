package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

const (
	rippleAPIURL                    = "s.altnet.rippletest.net:51233"
	reqTypeLedger                   = "ledgers"
	reqTypeTransactions             = "transactions"
	reqBodyLedgersSubscription      = "{\"id\":\"Example watch for new validated ledgers\",\"command\": \"subscribe\",\"streams\": [\"ledger\"]}"
	reqBodyTransactionsSubscription = "{\"id\":\"Example watch for new validated transactions\",\"command\": \"subscribe\",\"streams\": [\"transactions\"]}"
)

func main() {
	reqType := os.Args[1]
	if reqType == reqTypeLedger {
		listenToLedgersUpdates()
	} else if reqType == reqTypeTransactions {
		listenToTransactionsUpdates()
	} else {
		log.Fatal(errors.New(fmt.Sprintf("wrong subscription type provided, valid values are %s and %s", reqTypeLedger, reqTypeTransactions)))
	}
}

func listenToLedgersUpdates() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: rippleAPIURL, Path: "/"}
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

	u := url.URL{Scheme: "wss", Host: rippleAPIURL, Path: "/"}
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
