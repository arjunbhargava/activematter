package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

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

func main() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200.
	client, err := elastic.NewClient(elastic.SetBasicAuth("elastic", "changeme"))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	_, err = client.IndexExists("arxiv-titles").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	// Get tweet with specified ID
	get1, err := client.Get().
		Index("arxiv-titles").
		Type("Paper").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	// _, err = client.Flush().Index("twitter").Do(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// Search with a term query
	termQuery := elastic.NewTermQuery("title", "paper")
	searchResult, err := client.Search().
		Index("arxiv-titles"). // search in index "twitter"
		Query(termQuery).      // specify the query
		// Sort("authors", true). // sort by "user" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true).     // pretty print request and response JSON
		Do(ctx)           // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println(searchResult)

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Paper
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Paper); ok {
			fmt.Printf("Paper by %s: %s\n", t.Authors[0], t.Descriptions[0])
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d papers\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d papers\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Paper
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Paper by %s: %s\n", t.Authors[0], t.Title)
		}
	} else {
		// No hits
		fmt.Print("Found no papers\n")
	}
}
