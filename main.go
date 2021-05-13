package main

import (
	log "github.com/EntropyPool/entropy-logger"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"os"
)

func main() {
	app := &cli.App{
		Name:                 "filecoin-wallet",
		Usage:                "Internal wallet to manage filecoin accounts",
		Version:              "0.1.0",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./wallet-service.conf",
			},
		},
		Action: func(cctx *cli.Context) error {
			configFile := cctx.String("config")
			server := NewWalletServer(configFile)
			if server == nil {
				return xerrors.Errorf("cannot create wallet server")
			}
			err := server.Run()
			if err != nil {
				return xerrors.Errorf("fail to run wallet server: %v", err)
			}

			ch := make(chan int)
			<-ch

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(log.Fields{}, "fail to run %v: %v", app.Name, err)
	}
}
