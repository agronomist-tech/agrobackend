package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Clickhouse struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Solana struct {
		Servers []string `yaml:"servers"`
	}

	Port int
}

func LoadConfig() *config {
	log.Println("Load config from command line")
	cfg := config{}

	filename, _ := filepath.Abs("config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println("Problem on open config file: ", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)

	if err != nil {
		log.Println("Problem on parse config file: ", err)
		os.Exit(1)
	}

	port := flag.Int("port", 8090, "HTTP Port")
	cmdClickHost := flag.String("click-host", "", "Host to clickhouse database (Required)")
	cmdClickPort := flag.String("click-port", "", "Port to clickhouse database  (Required)")


	flag.Parse()

	cfg.Port = *port

	if len(*cmdClickHost) > 0 {
		cfg.Clickhouse.Host = *cmdClickHost
	}

	if len(*cmdClickPort) > 0 {
		cfg.Clickhouse.Port = *cmdClickPort
	}

	if len(cfg.Clickhouse.Host) == 0 || len(cfg.Clickhouse.Port) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return &cfg
}
