package main

import (
	"fmt"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal"
)

func main() {
	singleton := boxed.GetInstance()

	server := internal.SetupControllers()
	server.Start(fmt.Sprintf(":%v", singleton.BackendPort))
}
