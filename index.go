package main

import (
	"encoding/gob"
)

type Index struct{
	mu sync.RWMutex
	index map[string][]int
	docStore map[int]Document
}

func (idx *Index) AddDocument(doc Document){
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.docStore[doc.ID] = doc
	for _, token := range analyze(doc.Text) {
		ids := idx.index[token]
		if ids != nil && ids[len(ids -1 )] == doc.ID{
			continue
		}
		idx.index[token] = append(ids, doc.ID)
	}
}

func (idx *Index) Save(filePath string) error{
	idx.mu.RLock()

	defer idx.mu.RUnlock()

	file, err := os.Create(filePath)
	if err != nil{
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(idx.index); err != nil{
		return err
	}

	if err := encoder.Encode(idx.docStore); err != nil{
		return err
	}

	return nil
}

func (idx *Index) AddStreamed(docChan<-chan Document){
	var wg sync.WaitGroup
	numWorkers = 5

	for i:=0; i<numWorkers; i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan{
				idx.AddDocument(doc)
			}
		}()
	}
	wg.Wait()
}

func (idx *Index) Search(text string) []int{
	var result  [] int
	for _, token := range analyze(text) {
		idx.mu.RLock()
		if ids, ok := idx.index[token]; ok{
			if result == nil{
				result = ids
			}else {
				result = Intersection(result, ids)
			}
		} else{
			idx.mu.RUnlock()
			return nil
		}
		idx.mu.RUnlock()
	}
	return result
}
