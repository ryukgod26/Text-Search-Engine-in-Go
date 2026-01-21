# Text Search Engine in Go

A high-performance, full-text search engine written in Go with support for concurrent document indexing, natural language processing, and persistent index storage.

## Features

**Full-Text Search** - Search across multiple documents with multi-term queries  
**Concurrent Processing** - Multi-goroutine document loading and indexing with worker pools  
**Text Normalization** - Lowercase conversion, stopword removal, and stemming  
**Persistent Storage** - Save and load indexes using GOB encoding  
**Thread-Safe** - RWMutex-based synchronization for concurrent read/write operations  
**Efficient Intersection** - Fast set intersection algorithm for multi-term search results  

### Data Flow

```
XML File (gzip) → LoadDocuments() → Document Channel → AddStreamed() 
                                                    ↓
                                            5 Worker Goroutines
                                                    ↓
                                          Tokenizer.analyze()
                                                    ↓
                                    (lowercase → stopword → stem)
                                                    ↓
                                          Index Storage (Map)
                                                    ↓
                                     Save as GOB file / Search
```

## Usage

### Basic Search
```bash
./search-engine -q "test" -p "documents.xml.gz" -index "index.gob"
```

### Command-Line Flags
| Flag | Description | Default |
|------|-------------|---------|
| `-q` | Search query | "test" |
| `-p` | Path to gzip XML file(s), comma-separated | Required |
| `-index` | Path to save/load index file | - |

### Example
```bash
./search-engine -q "machine learning" -p "wiki_dump.xml.gz" -index "wiki_index.gob"
```

**Output:**
```
Search found 156 documents in 245ms

---------- | ---- | --------
Doc Id     | Title                        | Snippet
---------- | ---- | --------
1234       | Machine Learning Basics      | Machine learning is a subset of...
5678       | Deep Neural Networks         | Neural networks are computing...
...
---------- | ---- | --------
```

## Search Behavior

- **Multi-term queries** use AND logic: all terms must be present in a document
- **Case-insensitive** search (all text converted to lowercase)
- **Stopwords ignored**: Common words don't affect search
- **Stemming applied**: "running" finds documents with "run", "runs", "runner"

## Testing

Run the test suite:
```bash
go test -v
```

### Test Coverage
- Index creation and document addition
- Search with single and multiple terms
- Document retrieval by ID
- Index persistence (save/load)
- Intersection algorithm for multi-term queries
- Stopword filtering
- Duplicate document handling

## Performance Characteristics

- **Indexing**: 5 concurrent workers process documents
- **Search**: O(n*m) where n = number of search terms, m = average documents per term
- **Memory**: Inverted index grows with unique tokens and document count
- **Concurrency**: RWMutex allows unlimited concurrent readers, exclusive writer access

## Project Structure

```
├── main.go          # Entry point, command-line interface
├── index.go         # Core indexing and search logic
├── loader.go        # XML document loading
├── tokenizer.go     # Text analysis pipeline
├── filter.go        # Text filters (lowercase, stopwords, stemming)
├── index_test.go    # Unit tests
├── go.mod           # Go module definition
└── README.md        # This file
```
