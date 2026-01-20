package main

import (
	"compress/gzip"
	"encoding/xml"
	"os"
)

type Document struct{
	Title string `xml: "title"`
	URL string `xml: "url"`
	Text string `xml: "abstract"`
	ID int
}

func LoadDocuments(path string, docChan chan<- Document) error {
	f, err := os.Open(path)
	if err != nil{
		return err
	}

	gz, err := gzip.NewReader(f)
	if err !=nil{
		return err
	}

	defer gz.Close()
	
	dec := xml.NewDecoder(gz)
	dump := struct{
		Documents []Document `xml:"doc"`
	}{}

	if err := dec.Decode(&dump); err != nil{
		return err
	}

	for i, doc := range dump.Documents{
		doc.ID = i
		docChan <- doc
	}
	return nil
}

