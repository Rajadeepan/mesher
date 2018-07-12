package main

import (
	_ "net/http/pprof"

	_ "github.com/go-chassis/mesher/egress/chassis"
	_ "github.com/go-chassis/mesher/handler"
	_ "github.com/go-chassis/mesher/protocol/dubbo/client/chassis"
	_ "github.com/go-chassis/mesher/protocol/dubbo/server"
	_ "github.com/go-chassis/mesher/protocol/dubbo/simpleRegistry"
	_ "github.com/go-chassis/mesher/protocol/http"
	"github.com/go-chassis/mesher/server"

	_ "github.com/go-chassis/go-chassis/config-center" //use config center
)

func main() {
	server.Run()
}
