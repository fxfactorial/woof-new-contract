package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/pkg/browser"
)

var (
	clientDial = flag.String(
		"client_dial", "ws://localhost:8545", "could be websocket or IPC",
	)
)

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Ethereum Contract Creation")
	systray.SetTooltip("new contract creations")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	fmt.Println("exit exit")
}

func etherscan(h common.Hash) string {
	return "https://etherscan.io/tx/" + h.Hex()
}

func program() error {
	flag.Parse()

	handle, err := ethclient.Dial(*clientDial)
	if err != nil {
		return err
	}

	ch := make(chan *types.Header)
	sub, err := handle.SubscribeNewHead(context.Background(), ch)

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-sub.Err():
				// deal with it later
			case h := <-ch:
				fmt.Println(h)
				//
			}
		}
	}()

	go func() {
		err := beeep.Notify(
			"new ethereum contract made", "check menu", "information.png",
		)

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Addr", "new contract link")
		ch := make(chan struct{})
		mQuit.ClickedCh = ch

		go func() {
			for range ch {
				if err := browser.OpenURL(
					"https://etherscan.io/tx/0x9a92071acb900a6a759c52d8086b2cd1252ac6a26f2b195aa8d5bc496e692349",
				); err != nil {
					fmt.Println("whatever", err)
				}
			}
		}()
		_ = mQuit
		if err != nil {
			panic(err)
		}

	}()

	systray.Run(onReady, onExit)

	return nil
}

func main() {
	if err := program(); err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}
}
