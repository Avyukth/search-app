package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/clbanning/mxj/v2"
)

type Parser struct {
	errCh chan error
	wg    *sync.WaitGroup
}

func NewParser() *Parser {
	return &Parser{
		errCh: make(chan error),
		wg:    &sync.WaitGroup{},
	}
}
func (p *Parser) Parse(filePath string) (map[string]interface{}, error) {
	xmlData, err := os.ReadFile(filePath)

	if err != nil {
		p.errCh <- fmt.Errorf("Error reading XML file %s: %v", filePath, err)
		return nil, err
	}

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		p.errCh <- fmt.Errorf("Error unmarshalling XML from file %s: %v", filePath, err)
		return nil, err
	}
	mv["indexing"] = false
	return mv, nil
}

func (p *Parser) BuildPatent(data map[string]interface{}, storageID string) (*mongo.Patent, error) {
	var patent mongo.Patent

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	log.Println("JSON parsed XML file: ", jsonData)

	err = json.Unmarshal(jsonData, &patent)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}

	patent.PatentStorageID = storageID

	return &patent, nil
}
