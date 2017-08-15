package xmlutils

/* Takes in an OAI file and generates a parsed XML object. If no valid file is
given, uses the local data variable which is a parsed version. */

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	// "net/url"
	"os"
	// "time"
)

// var example_XML = []byte(`
// 	<?xml version="1.0" encoding="UTF-8"?>
// 	<Response>
// 	  <ListRecords>
// 	    <record>
// 	      <header status="">
// 	        <identifier>oai:arXiv.org:0704.0004</identifier>
// 	        <datestamp>2007-05-23</datestamp>
// 	        <setSpec>math</setSpec>
// 	      </header>
// 				<metadata>
// 	        <oai_dc xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd">
// 	          <title>A determinant of Stirling cycle numbers counts unlabeled acyclic  single-source automata</title>
// 	          <creator>Callan, David</creator>
// 	          <subject>Mathematics - Combinatorics</subject>
// 	          <subject>05A15</subject>
// 	          <description>  We show that a determinant of Stirling cycle numbers counts unlabeled acyclicsingle-source automata. The proof involves a bijection from these automata tocertain marked lattice paths and a sign-reversing involution to evaluate thedeterminant.</description>
// 	          <description>Comment: 11 pages</description>
// 	          <date>2007-03-30</date>
// 	          <type>text</type>
// 	          <identifier>http://arxiv.org/abs/0704.0004</identifier>
// 	        </oai_dc>
// 	      </metadata>
// 	      <about>
// 	      </about>
// 	    </record>
// 			<record>
// 	      <header status="">
// 	        <identifier>oai:arXiv.org:0704.0010</identifier>
// 	        <datestamp>2007-05-23</datestamp>
// 	        <setSpec>math</setSpec>
// 	      </header>
// 	      <metadata>
// 	        <oai_dc xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd">
// 	          <title>Partial cubes: structures, characterizations, and constructions</title>
// 	          <creator>Ovchinnikov, Sergei</creator>
// 	          <subject>Mathematics - Combinatorics</subject>
// 	          <description>  Partial cubes are isometric subgraphs of hypercubes. Structures on a graphdefined by means of semicubes, and Djokovi\'{c}'s and Winkler's relations playan important role in the theory of partial cubes. These structures are employedin the paper to characterize bipartite graphs and partial cubes of arbitrarydimension. New characterizations are established and new proofs of some knownresults are given.  The operations of Cartesian product and pasting, and expansion andcontraction processes are utilized in the paper to construct new partial cubesfrom old ones. In particular, the isometric and lattice dimensions of finitepartial cubes obtained by means of these operations are calculated.</description>
// 	          <description>Comment: 36 pages, 17 figures</description>
// 	          <date>2007-03-31</date>
// 	          <type>text</type>
// 	          <identifier>http://arxiv.org/abs/0704.0010</identifier>
// 	        </oai_dc>
// 	      </metadata>
// 	      <about>
// 	      </about>
// 	    </record>
// </ListRecords>
// </Response>
// `)

type Response struct {
	XMLName    xml.Name `xml:"Response"`
	ListRecord ListRecord
}

type ListRecord struct {
	XMLName    xml.Name `xml:"ListRecords"`
	RecordList []struct {
		Metadata Metadata `xml:"oai_dc"`
	} `xml:"record>metadata"`
}

type Metadata struct {
	Title        string   `xml:"title"`
	Creator      []string `xml:"creator"`
	Descriptions []string `xml:"description"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Simple parse of OAI-PMH arXiv data.
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
	fmt.Printf("Printing preview: %s", buf[:100])

	err = xml.Unmarshal(buf, &r)

	if err != nil {
		fmt.Println("Error parsing XML file:", err)
	} else {
		fmt.Println(r)
	}

	for _, record := range r.ListRecord.RecordList {
		fmt.Println("-------------")
		fmt.Println("Title: " + record.Metadata.Title)
		fmt.Print("Author: ")
		fmt.Print(record.Metadata.Creator)
		fmt.Println()
	}

	return r
}
