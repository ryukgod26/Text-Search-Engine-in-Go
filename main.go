package main

import (
	"flag"
	"log"
	"strings"
	"sync"
	"time"
)

func main() {
	var query,indexPath,files string

	flag.StringVar(&query,"q","test","Search query")
	flag.Parse()

	if files == ""{
		log.Fatal("No Files Given atleast one file is required with '-p' flag")
	}

	filePaths := strings.Split(files,",")
	log.Println("Full TExt Search")
	idx := NewIndex()
	start := time.Now()

	if err := idx.Load(indexPath); err == nil{
		log.Printf("Loaded Index form %s into %v",indexPath,time.Since(start))
	}else{
		var wg sync.WaitGroup
		docChan := make(chan Document,100)

		go func(){
			idx.AddStreamed(docChan)
		}()

		for _,path := range filePaths{
			wg.Add(1)
			go func(p string){
				defer wg.Done()
				log.Printf("Loading Documents from %s",p)

				if err := LoadDocuments(p, docChan); err != nil{
					log.Printf("Failed to Load Documents form %s: %v",p, err)
				}
			}(path)
		}
		wg.Wait()
		close(docChan)
		log.Printf("Indexed Documents in %v",time.Since(start))

		if err := idx.Save(indexPath); idx != nil{
			log.Fatalf("Failed to Save Index: %v", err)
		}
		log.Printf("Saved Index to %s",indexPath)
	}

	matchedIDs := idx.Search(query)
	log.Printf("Search found %d documents in %v",len(matchedIDs),time.Since(start))
	idx.PrintResultsTable(matchedIDs)

}