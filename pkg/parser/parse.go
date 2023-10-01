package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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

func (p *Parser) ParseToStruct(filePath string, storageID string) (*mongo.Patent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var patentGrant mongo.UsPatentGrant
	err = xml.Unmarshal(byteValue, &patentGrant)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}

	return p.BuildPatent(&patentGrant, storageID)
}

func (p *Parser) BuildPatent(patentGrant *mongo.UsPatentGrant, storageID string) (*mongo.Patent, error) {
	if patentGrant == nil {
		return nil, errors.New("patentGrant cannot be nil")
	}

	var inventorNames []string
	for _, inventor := range patentGrant.UsBibliographicDataGrant.UsParties.Inventors.Inventor {
		inventorNames = append(inventorNames, inventor.Addressbook.FirstName.Text+" "+inventor.Addressbook.LastName.Text)
	}

	patent := mongo.Patent{
		PatentTitle:     patentGrant.UsBibliographicDataGrant.InventionTitle.Text,
		PatentNumber:    patentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.DocNumber.Text,
		InventorNames:   inventorNames,
		AssigneeName:    patentGrant.UsBibliographicDataGrant.Assignees.Assignee.Addressbook.Orgname.Text,
		ApplicationDate: patentGrant.UsBibliographicDataGrant.ApplicationReference.DocumentID.Date.Text,
		IssueDate:       patentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.Date.Text,
		DesignClass:     patentGrant.UsBibliographicDataGrant.ClassificationsCpc.MainCpc.ClassificationCpc.Section.Text,
		PatentStorageID: storageID,
	}

	return &patent, nil
}
