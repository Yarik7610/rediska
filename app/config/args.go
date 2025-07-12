package config

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Args struct {
	Host       string
	Port       int
	DBDir      string
	DBFilename string
	ReplicaOf  *replicaOfConfig
}

type replicaOfConfig struct {
	Host string
	Port int
}

func (replcfg *replicaOfConfig) String() string {
	return fmt.Sprintf("%s %d", replcfg.Host, replcfg.Port)
}

func NewArgs() *Args {
	host := flag.String("host", "0.0.0.0", "The host of redis server")
	port := flag.Int("port", 6379, "The port of redis server")
	dir := flag.String("dir", "", "The path to RDB")
	filename := flag.String("dbfilename", "", "The filename of RDB")
	replicaOf := flag.String("replicaof", "", "The host and port of master server")

	flag.Parse()

	replicaOfConfig, err := configReplicaOf(replicaOf)
	if err != nil {
		log.Fatalf("wrong replicaof argument format: %v\n", err)
	}

	return &Args{
		Host:       *host,
		Port:       *port,
		DBDir:      *dir,
		DBFilename: *filename,
		ReplicaOf:  replicaOfConfig,
	}
}

func configReplicaOf(replicaOf *string) (*replicaOfConfig, error) {
	if *replicaOf == "" {
		return nil, nil
	}

	replicaOfArgs := strings.Split(*replicaOf, " ")
	if len(replicaOfArgs) != 2 {
		return nil, fmt.Errorf("provide host and port separated by space, e.g: '0.0.0.0 6379'")
	}

	masterHost := replicaOfArgs[0]
	masterPort := replicaOfArgs[1]
	if masterHost == "" || masterPort == "" {
		return nil, fmt.Errorf("host and port must be non-empty")
	}

	atoiMasterPort, err := strconv.Atoi(masterPort)
	if err != nil {
		return nil, fmt.Errorf("port should be decimal: %v", err)
	}

	return &replicaOfConfig{Host: masterHost, Port: atoiMasterPort}, nil
}
