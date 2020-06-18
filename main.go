package main

import (
	"app/application/worker"
	"app/framework"
	"app/tests/seed"
	"fmt"
	"github.com/joho/godotenv"
)

func main()  {

	godotenv.Load()

	framework.StartConnexions()
	defer framework.CloseConnexions()


	fmt.Println("starting")
	seed.RabbitMsgSeed() //populando queue com msg
	worker.Listen()
	fmt.Println("ending")

}
