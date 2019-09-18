
package main
import(
	"log"
	"micro/server"

	
	"micro/tools/koala/demo/router"
	

	
	"micro/tools/koala/demo/generate/hello"
	
)

var routerServer = &router.RouterServer{}

func main() {

    err := server.Init("hello")
    if err != nil {
        log.Fatal("init service failed, err:%v", err)
        return
    }

		hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
		server.Run()
	}
