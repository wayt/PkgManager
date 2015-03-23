package main

import (
	"flag"
	"fmt"
	"github.com/gohappy/happy"
	"github.com/wayt/pkgmanager/action"
	"github.com/wayt/pkgmanager/plugin"
	"log"
)

// Commandline flags
var (
	configFile = flag.String("config.file", "pkg-manager.conf.json", "Package manager configuration file name.")
)

func main() {

	if err := plugin.ReadConfig(*configFile); err != nil {
		log.Fatalln("fatal:", err)
	}

	if err := plugin.SetupMongoDB(); err != nil {
		log.Fatalln("fatal:", err)
	}

	plugin.SetupStorage()

	app := happy.NewAPI()
	app.Error404Handler = func(context *happy.Context, err interface{}) {
		context.Response.WriteHeader(404)
		context.Response.Write([]byte(`{"error": ["route not found"]}`))
	}

	app.PanicHandler = func(context *happy.Context, err interface{}) {
		context.Response.WriteHeader(500)
		context.Response.Write([]byte(`{"error": ["internal error"]}`))
	}

	action.RegisterActions(app)

	endpoint := fmt.Sprintf("%s:%d", plugin.Config.BindIp, plugin.Config.BindPort)

	log.Println("info:", "Listen on", endpoint)

	if err := app.Run(endpoint); err != nil {
		log.Fatalln("fatal:", err)
	}
}
