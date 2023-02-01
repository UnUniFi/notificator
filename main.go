package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Notificator struct {
	TmProvider string
	Logger     tmLog.Logger
	DB         *leveldb.DB
}

func NewNotificator(tmProvider string, logger tmLog.Logger) (*Notificator, error) {
	path := os.ExpandEnv("$HOME/monitor/db")
	db, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		return nil, err
	}

	return &Notificator{
		TmProvider: tmProvider,
		Logger:     logger,
	}, nil
}

func (notificator Notificator) Close() {
	notificator.DB.Close()
}

func (notificator Notificator) Start() {
	client, err := tmClient.New(notificator.TmProvider)
	if err != nil {
		notificator.Logger.Error("failed to initialize a client", "err", err)
		os.Exit(1)
	}
	client.SetLogger(notificator.Logger)

	if err := client.Start(); err != nil {
		notificator.Logger.Error("failed to start a client", "err", err)
		os.Exit(1)
	}

	defer client.Stop() //nolint:errcheck

	// Subscribe to all tendermint transactions
	query := "tm.event = 'Tx'"
	out, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		notificator.Logger.Error("failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-out:
			tx, ok := result.Data.(tmTypes.EventDataTx)
			if !ok {
				notificator.Logger.Error("new tx: error while extracting event data from new tx")
			}
			notificator.Logger.Info("New transaction witnessed")

			// Iterate over each event in the transaction
			for _, event := range tx.Result.Events {
				eventType := event.GetType()

				switch eventType {
				case "notify":
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					err := notificator.notify(event.GetAttributes())
					if err != nil {
						notificator.Logger.Error(err.Error())
					}
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}

func (notificator Notificator) notify(attributes []tmKv.Pair) error {
	return nil
}
