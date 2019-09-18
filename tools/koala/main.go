package main

import (
	"log"
	"os"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

func main() {
	var opt Option

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Value:       "./test.proto",
			Usage:       "idl filename",
			Destination: &opt.Proto3Filename,
		},
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc client code",
			Destination: &opt.GenServerCode,
		},
		cli.StringFlag{
			Name:        "p",
			Value:       "",
			Usage:       "prefix of package",
			Destination: &opt.Prefix,
		},
	}

	app.Action = func(c *cli.Context) error {
		err := genMgr.Run(&opt)
		if err != nil {
			glog.Infof("code generator failed, err:%v\n", err)
			return err
		}

		glog.Info("code generate succ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
