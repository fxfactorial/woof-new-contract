package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	fmt.Println("exit exit")
	// clean up here
}

func program() error {
	flag.Parse()

	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("speak speak")
		err := beeep.Notify(
			"new ethereum contract made", "check menu", "information.png",
		)

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Addr", "new contract link")
		ch := make(chan struct{})
		mQuit.ClickedCh = ch
		go func() {
			for range ch {
				fmt.Println("click happened!")
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
	fmt.Println("after systray worked")
	handle, err := ethclient.Dial(*clientDial)

	_ = handle

	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := program(); err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}
}
