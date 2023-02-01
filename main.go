package main

import (
	sdk "github.com/UnUniFi/notificator-sdk"
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

	// notificator.RegisterEventHandler()

	notificator.Start()
}
