package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Index struct {
	mu       sync.RWMutex
	index    map[string][]int
	docStore map[int]Document
}

func NewIndex() *Index {
	return &Index{
		index:    make(map[string][]int),
		docStore: make(map[int]Document),
	}
}

func (idx *Index) AddDocument(doc Document) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.docStore[doc.ID] = doc
	for _, token := range analyze(doc.Text) {
		ids := idx.index[token]
		if ids != nil && ids[len(ids)-1] == doc.ID {
			continue
		}
		idx.index[token] = append(ids, doc.ID)
	}
}

func (idx *Index) Save(filePath string) error {
	idx.mu.RLock()

	defer idx.mu.RUnlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(idx.index); err != nil {
		return err
	}

	if err := encoder.Encode(idx.docStore); err != nil {
		return err
	}

	return nil
}

func (idx *Index) AddStreamed(docChan <-chan Document) {
	var wg sync.WaitGroup
	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan {
				idx.AddDocument(doc)
			}
		}()
	}
	wg.Wait()
}

func (idx *Index) Search(text string) []int {
	var result []int
	for _, token := range analyze(text) {
		idx.mu.RLock()
		if ids, ok := idx.index[token]; ok {
			if result == nil {
				result = ids
			} else {
				result = Intersection(result, ids)
			}
		} else {
			idx.mu.RUnlock()
			return nil
		}
		idx.mu.RUnlock()
	}
	return result
}

func (idx *Index) PrintResultsTable(matchedIDs []int) {
	fmt.Printf("\n%-10s | %-40s | %s\n", "Doc Id", "Title", "Snippet")
	fmt.Println(strings.Repeat("-", 105))
	for _, id := range matchedIDs {
		if doc, found := idx.GetDocumentByID(id); found {
			snippet := doc.Text
			if len(snippet) > 50 {
				snippet = snippet[:47] + "..."
			}
			fmt.Printf("%-10d | %-40s | %s\n", doc.ID, doc.Title, snippet)
		}
	}
	fmt.Println(strings.Repeat("-", 105))
}

func Intersection(a, b []int) []int {
	var i, j int
	result := []int{}
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			result = append(result, a[i])
			i++
			j++
		}
	}
	return result
}

func (idx *Index) GetDocumentByID(id int) (Document, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	doc, exists := idx.docStore[id]
	return doc, exists
}

func (idx *Index) Load(filePath string) error{
	idx.mu.Lock()
	defer idx.mu.Unlock()

	file,err := os.Open(filePath)
	if err != nil{
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&idx.index); err != nil{
		return err
	}
	if err := decoder.Decode(&idx.docStore); err != nil{
		return nil
	}

	return nil
}