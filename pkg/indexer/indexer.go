package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/blevesearch/bleve/v2"
)

type SearchEngine struct {
	index bleve.Index
}

func NewSearchEngine(indexDir string) (*SearchEngine, error) {
	var index bleve.Index

	// Check if the index already exists
	if _, err := os.Stat(indexDir); errors.Is(err, os.ErrNotExist) {
		// Create a new index
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexDir, mapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		// Handle other os.Stat errors
		return nil, err
	} else {
		// Open the existing index
		index, err = bleve.Open(indexDir)
		if err != nil {
			return nil, err
		}
	}

	return &SearchEngine{index: index}, nil
}

func (se *SearchEngine) IndexPatent(patent *mongo.Patent) error {
	patentBytes, err := json.Marshal(patent)
	if err != nil {
		return fmt.Errorf("error marshalling patent: %v", err)
	}

	err = se.index.Index(patent.PatentStorageID, patent)
	if err != nil {
		return fmt.Errorf("error adding patent to index: %v", err)
	}

	err = se.index.SetInternal([]byte(patent.PatentStorageID), patentBytes)
	if err != nil {
		return fmt.Errorf("error setting internal patent: %v", err)
	}

	return nil
}

func (se *SearchEngine) SearchAndRetrievePatents(searchTerm string) ([]mongo.Patent, error) {
	query := bleve.NewQueryStringQuery(searchTerm)
	search := bleve.NewSearchRequest(query)
	searchResults, err := se.index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %v", err)
	}

	var patents []mongo.Patent
	for _, hit := range searchResults.Hits {
		id := hit.ID
		originalPatentBytes, err := se.index.GetInternal([]byte(id))
		if err != nil {
			return nil, fmt.Errorf("error getting internal patent: %v", err)
		}

		var originalPatent mongo.Patent
		err = json.Unmarshal(originalPatentBytes, &originalPatent)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling patent: %v", err)
		}

		patents = append(patents, originalPatent)
	}

	if len(patents) == 0 {
		return nil, fmt.Errorf("no patents found for search term: %s", searchTerm)
	}

	return patents, nil
}
