package main

import (
	"fmt"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal"
)

func main() {
	singleton := boxed.GetInstance()
	defer singleton.DbConn.Close()

	// It setups the controllers and then start the server
	server := internal.SetupControllers()
	server.Start(fmt.Sprintf(":%v", singleton.BackendPort))
}
