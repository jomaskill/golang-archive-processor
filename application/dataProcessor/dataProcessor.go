package dataProcessor

import (
	"app/application/archiveProcessor"
	"app/domain"
	"app/framework/elasticsearch"
	"errors"
	"log"
	"os"
	"sync"
)

func Init(msg map[string]string, wg *sync.WaitGroup)  {

	file, err := os.Open("./archives/"+ msg["file"])
	if err != nil{
		log.Fatal(err)
	}

	content, err := archiveProcessor.Process(file)
	if err != nil{
		log.Fatal(err)
	}

	data := parse(content, msg)

	err = prepare(data)
	if err != nil{
		log.Fatal(err)
	}

	wg.Done()
}

func parse(content []string, msg map[string]string)  []domain.Base{

	if _, ok := msg["data_type"]; ok {

		var blacklists []domain.Base

		for _, value := range content{
			var blacklist domain.Blacklist
			blacklist.New(value, msg)
			blacklists = append(blacklists, blacklist)
		}
		return blacklists
	}

	var whitelists []domain.Base

	for _, value := range content{
		var whitelist domain.Whitelist
		whitelist.New(value, msg)
		whitelists = append(whitelists, whitelist)
	}
	return whitelists
}

func prepare(data []domain.Base)  error{

	switch data[0].(type) {

	case domain.Blacklist:

		var done = make(chan bool, 1)
		if data[0].(domain.Blacklist).DataType == "staff" {
			go elasticsearch.Deletion("data_type", "staff", "blacklist", done)
			<-done
		}
		elasticsearch.BulkInsertion(data)

	case domain.Whitelist:

		elasticsearch.BulkInsertion(data)

	default:
		return errors.New("unknown type")
	}

	return nil
}
