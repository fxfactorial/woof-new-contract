package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

const t = "ws://localhost:8545"

var (
	clientDial = flag.String(
		"client_dial", t, "could be websocket or IPC",
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

func kovanEtherscan(h common.Hash) string {
	return "https://kovan.etherscan.io/tx/" + h.Hex()
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

	errCh := make(chan error)
	_ = errCh
	// come back to it , error thread reader ?,
	// does systray need to be on mainthread?

	go func() {
		for {
			select {
			case e := <-sub.Err():
				log.Fatal(e)
				// deal with it later
			case h := <-ch:
				block, err := handle.BlockByNumber(context.Background(), h.Number)
				if err != nil {
					log.Fatal(errors.Wrapf(err, "block by hash issue"))
				}

				for _, tx := range block.Transactions() {
					// new contract?
					if t := tx.To(); t == nil && len(tx.Data()) >= 4 {
						if err := beeep.Notify(
							"new ethereum contract made", "check menu", "information.png",
						); err != nil {
							log.Fatal(err)
						}

						systray.AddSeparator()
						newContract := systray.AddMenuItem(
							fmt.Sprintf("at block %d", h.Number), "new contract link",
						)

						ch := make(chan struct{})
						newContract.ClickedCh = ch
						hash := tx.Hash()
						go func() {
							for range newContract.ClickedCh {
								if err := browser.OpenURL(kovanEtherscan(hash)); err != nil {
									log.Fatal(errors.Wrapf(err, "browser fucked up?"))
								}
							}
						}()
					}
				}
			}
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
