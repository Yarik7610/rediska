package config

import "flag"

type Args struct {
	Port       int
	DBDir      string
	DBFilename string
}

func NewArgs() *Args {
	port := flag.Int("port", 6379, "The port of redis server")
	dir := flag.String("dir", "/tmp/redis-files", "The path to RDB")
	filename := flag.String("dbfilename", "rdbfile", "The filename of RDB")

	flag.Parse()

	return &Args{
		Port:       *port,
		DBDir:      *dir,
		DBFilename: *filename,
	}
}
