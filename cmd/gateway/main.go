// starts the LLM gateway server

package main

import (
	"context"
	"fmt"

	"github.com/vaibhavvvvv/obsintel/config"
	"github.com/vaibhavvvvv/obsintel/pkg/gateway"
)

func main() {
	fmt.Println("OBSINTEL starting...")

	config.Init()
	ctx := context.Background()

	gw := gateway.New(ctx)
	gw.Run(config.C.Port)

}
