package elasticsearch

import (
	"app/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

const batch = 1000

func conn() (*elasticsearch.Client, error){

	elasticUser := os.Getenv("ELASTIC_USER")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	elasticUrl := os.Getenv("ELASTIC_URL")
	elasticHost := os.Getenv("ELASTIC_HOST")
	elasticPort := os.Getenv("ELASTIC_PORT")
	elasticPortAux := os.Getenv("ELASTIC_PORTAUX")

	cfg := elasticsearch.Config{Addresses: []string{
		elasticUrl+elasticHost+":"+elasticPort,
		elasticUrl+elasticHost+":"+elasticPortAux,
	},
		Username: elasticUser,
		Password: elasticPassword,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil{
		return nil, err
	}

	_ , err = client.Info()
	if err != nil {
		return nil, err
	}

	return client, err
}

func BulkInsertion(data []domain.Base){

	count := len(data)
	index := data[0].Index()

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error
		raw map[string]interface{}
		blk *bulkResponse

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
	)

	client, err := conn()
	if err != nil{
		log.Fatalf("Error from elastic connection: %v", err)
	}

	createIndex(client, index)

	if count%batch == 0 {
		numBatches = count / batch
	} else {
		numBatches = (count / batch) + 1
	}

	start := time.Now().UTC()

	for i, a := range data {

		numItems++

		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		meta := []byte(fmt.Sprintf(`{ "index" : { } }%s`, "\n"))

		data, err := json.Marshal(a.Data())
		if err != nil {
			log.Fatalf("Cannot encode article : %s", err)
		}

		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		if i > 0 && i%batch == 0 || i == count-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			res, err = client.Bulk(bytes.NewReader(buf.Bytes()), client.Bulk.WithIndex(index))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
			}

			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}

			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {

						if d.Index.Status > 201 {

							numErrors++

							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {

							numIndexed++
						}
					}
				}
			}

			res.Body.Close()

			buf.Reset()
			numItems = 0
		}
	}

	result(start, numErrors, numIndexed)
}

func createIndex(client *elasticsearch.Client, index string) {

	if indexExists(client, index) == false{
		res, err := client.Indices.Create(index)
		if err != nil {
			log.Fatalf("Cannot create index: %s", err)
		}
		if res.IsError() {
			log.Fatalf("Cannot create index: %s", res)
		}
	}
}

func indexExists(client *elasticsearch.Client, index string) bool{

	if _, exist := client.Indices.Exists([]string{index}); exist != nil {
		return true
	}

	return false
}

func result(start time.Time, numErrors int, numIndexed int) {
	fmt.Print("\n")
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}

	log.Printf(
		"[%s] documents in %s (%s docs/sec)",
		humanize.Comma(int64(numIndexed)),
		dur.Truncate(time.Millisecond),
		humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
	)

}

func Deletion(filter string, value string, index string, done chan<- bool) {

	client, err := conn()
	if err != nil{
		log.Fatalf("Error from elastic connection: %v", err)
	}

	if indexExists(client, index){
		ids := search(filter, value, index, client)
		bulkDelete(index, ids, client)
	}

	done<-true
}

func bulkDelete(index string, ids []string, client *elasticsearch.Client) {

	var buf bytes.Buffer
	var mapResp interface{}
	var numItems int

	start := time.Now().UTC()

	for _, value := range ids {
		numItems++

		meta := []byte(fmt.Sprintf(`{ "delete" : { "_index" : "%s", "_id" : "%s" } }%s`, index, value, "\n"))
		buf.Grow(len(meta))
		buf.Write(meta)
	}
	res, err := client.Bulk(bytes.NewReader(buf.Bytes()), client.Bulk.WithIndex(index))
	if err != nil {
		log.Fatal( err)
	}

	if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		log.Fatalf("json.NewEncoder() ERROR: %v", err)
	}
	//fmt.Println(mapResp)
	result(start, 0, numItems)
}


func search(filter string, value string, index string, client *elasticsearch.Client) []string{

	ctx := context.Background()
	read := constructQuery(filter, value)

	type Response struct {
		Hits struct {
			Values []struct {
				Id string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var mapResp Response
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("json.NewEncoder() ERROR: %v", err)
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(read),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("client.Search ERROR: %v", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		log.Fatalf("json.NewEncoder() ERROR: %v", err)
	}

	var ids []string
	for _, v := range mapResp.Hits.Values {
		ids = append(ids, v.Id)
	}

	return ids
}

func constructQuery(filter string, value string) *strings.Reader {

	var query = `{"query": {`

	var q = `
		"bool": {
			"filter": {
				"term": {
					"`+filter+`" : "`+value+`"
				}
			}
		}`

	query = query + q

	query = query + `}, "size": 10000}`

	isValid := json.Valid([]byte(query))

	if isValid == false {
		log.Fatalf("constructQuery() ERROR: query string not valid: %v", query)
	}

	var b strings.Builder
	b.WriteString(query)

	read := strings.NewReader(b.String())

	return read
}