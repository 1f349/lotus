package main

import (
	"flag"
	"github.com/1f349/lotus/api"
	"github.com/1f349/mjwt"
	"github.com/1f349/violet/utils"
	exitReload "github.com/MrMelon54/exit-reload"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "conf", "", "/path/to/config.yml : path to the config file")
	flag.Parse()

	if configPath == "" {
		log.Println("[Lotus] Error: config flag is missing")
		return
	}

	openConf, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("[Lotus] Error: missing config file")
		} else {
			log.Println("[Lotus] Error: open config file: ", err)
		}
		return
	}

	var conf Conf
	err = yaml.NewDecoder(openConf).Decode(&conf)
	if err != nil {
		log.Println("[Lotus] Error: invalid config file: ", err)
		return
	}

	wd := filepath.Dir(configPath)

	verify, err := mjwt.NewMJwtVerifierFromFile(filepath.Join(wd, "signer.public.pem"))
	if err != nil {
		log.Fatalf("[Lotus] Failed to load MJWT verifier public key from file '%s': %s", filepath.Join(wd, "signer.public.pem"), err)
	}

	userAuth := &api.AuthChecker{Verify: verify, Aud: conf.Audience}
	srv := api.SetupApiServer(conf.Listen, userAuth, &conf.SendMail, &conf.Imap)
	log.Printf("[Lotus] Starting API server on: '%s'\n", srv.Addr)
	go utils.RunBackgroundHttp("Lotus", srv)

	exitReload.ExitReload("Lotus", func() {}, func() {
		// stop server
		srv.Close()
	})
}
