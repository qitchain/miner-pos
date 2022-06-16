package main

import (
	"chia-miner/app"
	export "chia-miner/export"
	"chia-miner/miner"
	config2 "chia-miner/pkg/config"
	"chia-miner/pkg/log"
	"chia-miner/utils"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var config = flag.String("config", "config.yaml", "configuration file")
var exportFarmer = flag.Bool("export", false, "export farmer private")

func main() {
	flag.Parse()
	if *exportFarmer {
		fmt.Println("Please enter the chia mnemonic:")
		for {
			mnemonic := utils.GetInput()
			// 24 mnemonic words
			if len(strings.Split(strings.TrimSpace(mnemonic), " ")) != 24 {
				fmt.Println("Please enter chia mnemonics, separated by spaces")
				continue
			}
			farmerPublicKey, farmerPrivateKey, err := export.GetFarmerPrivateKeyByMnemonic(mnemonic)
			if err != nil {
				fmt.Println("export failed ~ ", err)
				return
			}
			fmt.Println("Farmer public key:", farmerPublicKey)
			fmt.Println("Farmer private key:", farmerPrivateKey)
			return
		}
	}
	fmt.Println("QitChain miner")
	fmt.Println("Version: ", app.Version, app.BuildVersion)
	fmt.Println("Build time: ", app.BuildTime)

	var cfg = &config2.Config{}
	if err := utils.LoadConfigFromFile(*config, cfg); err != nil {
		fmt.Printf("load config fail, error %v", err)
		return
	}

	if cfg.Rpc.Url == "" {
		cfg.Rpc.Url = "http://localhost:3332"
	}
	fmt.Println("Wallet address", cfg.Rpc.Url)

	if cfg.FarmerKey == nil {
		cfg.FarmerKey = make(map[string]string)
	}

	for _, v := range cfg.FarmerPrivateKey {
		list := strings.Split(v, " ")
		if len(list) == 24 {
			farmerPublicKey, farmerPrivateKey, err := export.GetFarmerPrivateKeyByMnemonic(v)
			if err != nil {
				logrus.Errorf("Failed to generate farmer private key %v", err)
				continue
			}
			cfg.FarmerKey[farmerPublicKey] = farmerPrivateKey
		} else {
			publicKey, err := export.GetFarmerPublicKey(v)
			if err != nil {
				logrus.Errorf("Wrong private key %v", err)
			} else {
				cfg.FarmerKey[publicKey] = v
			}
		}
	}

	log.InitLog(cfg.Log.Level, cfg.Log.File)
	miner.GetMiner().Start(cfg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	s := <-c
	logrus.Infof("Received signal %v", s)
}
