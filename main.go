package main

import (
	sdk "github.com/UnUniFi/notificator-sdk"
	"github.com/tendermint/tendermint/abci/types"
)

var appName = "ununifi-notificator"

func main() {
	config, err := sdk.LoadConfig(appName)
	if err != nil {
		panic(err)
	}
	notificator, err := sdk.NewNotificator(appName, *config)
	if err != nil {
		panic(err)
	}
	defer notificator.Close()

	notificator.RegisterEventHandler(EventHogeType, handleEventHoge)

	notificator.Start()
}

const EventHogeType = ""

func handleEventHoge(attributes []types.EventAttribute) error {
	return nil
}
