package servers

import "github.com/codecrafters-io/redis-starter-go/app/config"

type Server interface {
	Start()
}

func SpawnServer(args *config.Args) Server {
	if args.ReplicaOf == nil {
		return newMaster(args)
	} else {
		return newReplica(args)
	}
}
