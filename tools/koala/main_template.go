package main

var main_template = `
package main
import(
	"log"
	"micro/server"

	{{if not .Prefix}}
	"router"
	{{else}}
	"{{.Prefix}}/router"
	{{end}}

	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)

var routerServer = &router.RouterServer{}

func main() {

    err := server.Init("{{.Package.Name}}")
    if err != nil {
        log.Fatal("init service failed, err:%v", err)
        return
    }

		hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
		server.Run()
	}
`
