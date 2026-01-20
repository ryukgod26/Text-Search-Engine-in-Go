package main

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
