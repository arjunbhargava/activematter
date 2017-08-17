package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/arjunbhargava/activematter/data/xmlutils"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Paper is the structure for serializing/deserializing data
type Paper struct {
	Title          string    `json:"title"`
	Authors        []string  `json:"authors,omitempty"`
	Subjects       []string  `json:"subjects,omitempty"`
	Descriptions   []string  `json:"descriptions,omitempty"`
	SubmissionDate time.Time `json:"submitted,omitempty"`
	Identifier     url.URL   `json:"identifier,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 5,
		"number_of_replicas": 1
	},
	"mappings":{
		"paper":{
			"properties":{
				"title":{
					"type":"text",
					"fielddata": true,
				},
				"authors":{
					"type":"text",
					"store": true,
				},
				"subjects":{
					"type":"keyword"
				},
				"description":{
					"type":"text",
					"fielddata": true,
					"store": true,
				},
				"submitted":{
					"type":"date"
				},
				"identifier":{
					"type": "keyword"
				}
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func clientHelper(ctx context.Context) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetBasicAuth("elastic", "changeme"))
	if err != nil {
		// Handle error
		panic(err)
	}
	//Confirm the connection and check versioning
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return client
}

func insertElastic(ctx context.Context, client *elastic.Client, r xmlutils.Response) error {

	var paperToinsert Paper
	for i, record := range r.ListRecord.RecordList {
		fmt.Printf("Inserting paper %d", i+1)
		paperToinsert = Paper{Title: record.Metadata.Title,
			Authors:      record.Metadata.Creator,
			Descriptions: record.Metadata.Descriptions}

		put1, err := client.Index().
			Index("arxiv-titles").
			Type("Paper").
			Id(strconv.Itoa(i)).
			BodyJson(paperToinsert).
			Do(ctx)
		if err != nil {
			// Handle error
			return err
		}
		fmt.Printf("Indexed paper %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}

	return nil
}

func parseDirectory(ctx context.Context, s string) {

	fmt.Println("Directory name is: " + s)
	fileList, err := filepath.Glob(s + "*.gz")
	if err != nil {
		fmt.Println(err)
	} else {
		client := clientHelper(ctx)
		for _, record := range fileList {
			parseFile(ctx, client, record)
		}
	}
}

func parseFile(ctx context.Context, client *elastic.Client, record string) {
	r := xmlutils.ParseOAIXML(record)
	insertElastic(ctx, client, r)
}

/* Entry point /wrapper for parsing records. Takes either a directory
or a single file and indexes it in an Elastic instance. */

func main() {

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	//Go through files and index the records
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_records XML_dir_name or parse_recrds XML_file_name")
		return
	}

	fi, err := os.Stat(os.Args[1])

	if err != nil {
		fmt.Printf("Error: File not found")
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		fmt.Println("Parsing directory: " + os.Args[1])
		parseDirectory(ctx, os.Args[1])
	case mode.IsRegular():
		fmt.Println("Parsing file: " + os.Args[1])
		parseFile(ctx, clientHelper(ctx), os.Args[1])
	}
}
