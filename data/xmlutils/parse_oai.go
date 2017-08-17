package xmlutils

/* Takes in an OAI file and generates a parsed XML object. If no valid file is
given, returns an empty Response object */

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// Response is a high level response structure
type Response struct {
	XMLName    xml.Name `xml:"Response"`
	ListRecord ListRecord
}

//ListRecord is one level below that contains the array of records
type ListRecord struct {
	XMLName    xml.Name `xml:"ListRecords"`
	RecordList []struct {
		Metadata Metadata `xml:"oai_dc"`
	} `xml:"record>metadata"`
	Headers []Header `xml:"record>header"`
}

//Header grabs the header metadata
type Header struct {
	Identifier string `xml:"identifier"`
	Datestamp  string `xml:"datestamp"`
	SetSpec    string `xml:"setSpec"`
}

//Metadata object is each record is composed of other stuff
type Metadata struct {
	Title        string   `xml:"title"`
	Creator      []string `xml:"creator"`
	Descriptions []string `xml:"description"`
	Subjects     []string `xml:"subject"`
	Date         string   `xml:"date"`
	Identifier   string   `xml:"identifier"`
}

// ParseOAIXML is for a simple parse of OAI-PMH arXiv data.
func ParseOAIXML(fileName string) Response {

	fmt.Println("Parsing " + fileName + " ...")
	xmlFile, err := os.Open(fileName)

	var r Response

	if err != nil {
		fmt.Println("Error opening file:", err)
		return r
	}

	defer xmlFile.Close()

	buf, _ := ioutil.ReadAll(xmlFile)
	fmt.Printf("Printing preview: %s", buf[:50])

	err = xml.Unmarshal(buf, &r)

	if err != nil {
		fmt.Println("Error parsing XML file:", err)
	} else {
		fmt.Println(r)
	}

	for i, record := range r.ListRecord.RecordList {

		fmt.Println("-------------")
		fmt.Println(r.ListRecord.Headers[i].SetSpec)
		fmt.Println("Title: " + record.Metadata.Title)
		fmt.Print("Author: ")
		fmt.Print(record.Metadata.Creator)
		fmt.Println()
	}

	return r
}
