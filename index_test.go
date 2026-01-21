//Most of this file is Created by AI
package main

import (
    "os"
    "testing"
)

func TestNewIndex(t *testing.T) {
    idx := NewIndex()
    if idx == nil {
        t.Fatal("NewIndex returned nil")
    }
    if idx.index == nil {
        t.Error("index map not initialized")
    }
    if idx.docStore == nil {
        t.Error("docStore map not initialized")
    }
}

func TestAddDocument(t *testing.T) {
    idx := NewIndex()
    doc := Document{
        ID:    1,
        Title: "Test Document",
        Text:  "this is a test document with some content",
        URL:   "http://example.com",
    }

    idx.AddDocument(doc)

    if len(idx.docStore) != 1 {
        t.Errorf("expected 1 document in docStore, got %d", len(idx.docStore))
    }

    if len(idx.index) == 0 {
        t.Error("index should contain tokens")
    }
}

func TestAddDocumentDuplicate(t *testing.T) {
    idx := NewIndex()
    doc := Document{
        ID:    1,
        Title: "Test",
        Text:  "test content",
        URL:   "http://example.com",
    }

    idx.AddDocument(doc)
    idx.AddDocument(doc)

    // Should not duplicate entries
    if len(idx.docStore) != 1 {
        t.Errorf("expected 1 document, got %d", len(idx.docStore))
    }
}

func TestSearch(t *testing.T) {
    idx := NewIndex()
    
    docs := []Document{
        {ID: 1, Title: "Go Programming", Text: "learn golang programming language", URL: "http://example.com/1"},
        {ID: 2, Title: "Python Guide", Text: "learn python programming language", URL: "http://example.com/2"},
        {ID: 3, Title: "Web Dev", Text: "golang web development framework", URL: "http://example.com/3"},
    }

    for _, doc := range docs {
        idx.AddDocument(doc)
    }

    results := idx.Search("golang")
    if len(results) != 2 {
        t.Errorf("expected 2 results for 'golang', got %d", len(results))
    }

    results = idx.Search("programming golang")
    if len(results) != 1 {
        t.Errorf("expected 1 result for 'programming golang', got %d", len(results))
    }

    results = idx.Search("nonexistent")
    if len(results) != 0 {
        t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
    }
}

func TestGetDocumentByID(t *testing.T) {
    idx := NewIndex()
    doc := Document{
        ID:    42,
        Title: "Test Doc",
        Text:  "test content",
        URL:   "http://example.com",
    }

    idx.AddDocument(doc)

    retrieved, found := idx.GetDocumentByID(42)
    if !found {
        t.Error("document not found")
    }
    if retrieved.Title != "Test Doc" {
        t.Errorf("expected title 'Test Doc', got '%s'", retrieved.Title)
    }

    _, found = idx.GetDocumentByID(999)
    if found {
        t.Error("should not find non-existent document")
    }
}

func TestSaveAndLoad(t *testing.T) {
    idx := NewIndex()
    docs := []Document{
        {ID: 1, Title: "Doc1", Text: "golang programming", URL: "http://example.com/1"},
        {ID: 2, Title: "Doc2", Text: "python programming", URL: "http://example.com/2"},
    }

    for _, doc := range docs {
        idx.AddDocument(doc)
    }

    tempFile := "test_index.gob"
    defer os.Remove(tempFile)

    err := idx.Save(tempFile)
    if err != nil {
        t.Fatalf("failed to save index: %v", err)
    }

    idx2 := NewIndex()
    err = idx2.Load(tempFile)
    if err != nil {
        t.Fatalf("failed to load index: %v", err)
    }

    if len(idx2.docStore) != 2 {
        t.Errorf("expected 2 documents after load, got %d", len(idx2.docStore))
    }

    results := idx2.Search("golang")
    if len(results) != 1 {
        t.Errorf("expected 1 result after load, got %d", len(results))
    }
}

func TestIntersection(t *testing.T) {
    tests := []struct {
        name     string
        a        []int
        b        []int
        expected []int
    }{
        {"both empty", []int{}, []int{}, []int{}},
        {"one empty", []int{1, 2}, []int{}, []int{}},
        {"no overlap", []int{1, 2}, []int{3, 4}, []int{}},
        {"full overlap", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
        {"partial overlap", []int{1, 2, 3, 4}, []int{2, 3, 5}, []int{2, 3}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Intersection(tt.a, tt.b)
            if len(result) != len(tt.expected) {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
            for i, v := range result {
                if i >= len(tt.expected) || v != tt.expected[i] {
                    t.Errorf("expected %v, got %v", tt.expected, result)
                    break
                }
            }
        })
    }
}