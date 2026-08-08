// starts the LLM gateway server

package main

import (
	"context"
	"fmt"

	"github.com/vaibhavvvvv/obsintel/config"
	"github.com/vaibhavvvvv/obsintel/pkg/gateway"
	"github.com/vaibhavvvvv/obsintel/store"
)

func main() {
	fmt.Println("OBSINTEL starting...")
	config.Init()

	ctx := context.Background()

	//initializing db
	store.Init(ctx)
	defer store.Close()

	gateway.New(ctx).Run(config.C.Port)

}
