package worker

import (
	"app/application/dataProcessor"
	"app/framework/rabbitmq"
	"fmt"
	"log"
	"os"
	"sync"
)

func Listen()  {

	blackMsgs, err := rabbitmq.Consume(os.Getenv("RABBIT_QUEUE_BLACKLIST"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("blackMsgs ", blackMsgs)
	var blackChan sync.WaitGroup
	if len(blackMsgs) > 0 {
		blackChan.Add(len(blackMsgs))
		for _,value := range blackMsgs{
			go dataProcessor.Init(value, &blackChan)
		}
	}

	whiteMsgs, err:= rabbitmq.Consume(os.Getenv("RABBIT_QUEUE_WHITELIST"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("whiteMsgs ", whiteMsgs)
	var whiteChan sync.WaitGroup
	if len(whiteMsgs) > 0 {
		whiteChan.Add(len(whiteMsgs))
		for _,value := range whiteMsgs{
			go dataProcessor.Init(value, &whiteChan)
		}
	}

	blackChan.Wait()
	whiteChan.Wait()
}
