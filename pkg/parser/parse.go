// package parser

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"sync"

// 	"github.com/avyukth/search-app/pkg/database/mongo"
// 	"github.com/clbanning/mxj/v2"
// )

// type Parser struct {
// 	errCh chan error
// 	wg    *sync.WaitGroup
// }

// func NewParser() *Parser {
// 	return &Parser{
// 		errCh: make(chan error),
// 		wg:    &sync.WaitGroup{},
// 	}
// }
// func (p *Parser) Parse(filePath string) (map[string]interface{}, error) {
// 	xmlData, err := os.ReadFile(filePath)

// 	if err != nil {
// 		p.errCh <- fmt.Errorf("Error reading XML file %s: %v", filePath, err)
// 		return nil, err
// 	}

// 	mv, err := mxj.NewMapXml(xmlData)
// 	if err != nil {
// 		p.errCh <- fmt.Errorf("Error unmarshalling XML from file %s: %v", filePath, err)
// 		return nil, err
// 	}
// 	mv["indexing"] = false
// 	return mv, nil
// }

// // func (p *Parser) BuildPatent(data map[string]interface{}, storageID string) (*mongo.Patent, error) {
// // 	var patent mongo.Patent

// // 	jsonData, err := json.Marshal(data)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error marshalling data: %v", err)
// // 	}

// // 	log.Println("JSON parsed XML file: ", jsonData)

// // 	err = json.Unmarshal(jsonData, &patent)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error unmarshalling data: %v", err)
// // 	}

// // 	patent.PatentStorageID = storageID

// // 	return &patent, nil
// // }

// func (p *Parser) BuildPatent(data map[string]interface{}, storageID string) (*mongo.Patent, error) {
// 	var usPatentGrant mongo.UsPatentGrant
// 	// Populate usPatentGrant with actual data...
// file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Read the XML file into a byte array
// 	byteValue, err := io.ReadAll(file)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}

// 	// Initialize the UsPatentGrant struct
// 	var patentGrant UsPatentGrant

// 	// Unmarshal the byte array into the UsPatentGrant struct
// 	err = xml.Unmarshal(byteValue, &patentGrant)
// 	if err != nil {
// 		fmt.Println("Error unmarshalling XML:", err)
// 		return
// 	}

// 	// Extracting Patent struct
// 	patent := mongo.Patent{
// 		Title:          usPatentGrant.UsBibliographicDataGrant.InventionTitle.Text,
// 		Number:         usPatentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.DocNumber.Text,
// 		InventorNames:  getInventorNames(usPatentGrant.UsBibliographicDataGrant.UsParties.Inventors.Inventor),
// 		AssigneeName:   getAssigneeName(usPatentGrant.UsBibliographicDataGrant.Assignees.Assignee),
// 		ApplicationDate: usPatentGrant.UsBibliographicDataGrant.ApplicationReference.DocumentID.Date.Text,
// 		IssueDate:      usPatentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.Date.Text,
// 		DesignClass:    getClassification(usPatentGrant.UsBibliographicDataGrant.ClassificationsCpc),
// 	}

// }

// type Inventor struct {
// 	Name string
// }

// type Assignee struct {
// 	Name string
// }

// type ClassificationsCpc struct {
// 	Classification string
// }

// func getInventorNames(inventors []Inventor) string {
// 	var names []string
// 	for _, inventor := range inventors {
// 		names = append(names, inventor.Name)
// 	}
// 	return strings.Join(names, ", ")
// }

// func getAssigneeName(assignee Assignee) string {
// 	return assignee.Name
// }

// func getClassification(classificationsCpc ClassificationsCpc) string {
// 	return classificationsCpc.Classification
// }

package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/avyukth/search-app/pkg/database/mongo"
)

type Parser interface {
	Parse(filePath string) (*mongo.UsPatentGrant, error)
}

type fileParser struct{}

func NewParser() Parser {
	return &fileParser{}
}

func (p *fileParser) Parse(filePath string) (*mongo.UsPatentGrant, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	byteValue, err := fs.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var patentGrant mongo.UsPatentGrant
	err = xml.Unmarshal(byteValue, &patentGrant)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}

	return &patentGrant, nil
}

type PatentBuilder interface {
	BuildPatent(patentGrant *mongo.UsPatentGrant, storageID string) (*mongo.Patent, error)
}

type patentBuilder struct{}

func NewPatentBuilder() PatentBuilder {
	return &patentBuilder{}
}

func (pb *patentBuilder) BuildPatent(patentGrant *mongo.UsPatentGrant, storageID string) (*mongo.Patent, error) {
	if patentGrant == nil {
		return nil, errors.New("patentGrant cannot be nil")
	}

	patent := &mongo.Patent{
		PatentTitle:     patentGrant.UsBibliographicDataGrant.InventionTitle.Text,
		PatentNumber:    patentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.DocNumber.Text,
		InventorNames:   getInventorNames(patentGrant.UsBibliographicDataGrant.UsParties.Inventors.Inventor),
		AssigneeName:    getAssigneeName(patentGrant.UsBibliographicDataGrant.Assignees.Assignee),
		ApplicationDate: patentGrant.UsBibliographicDataGrant.ApplicationReference.DocumentID.Date.Text,
		IssueDate:       patentGrant.UsBibliographicDataGrant.PublicationReference.DocumentID.Date.Text,
		DesignClass:     getClassification(patentGrant.UsBibliographicDataGrant.ClassificationsCpc),
		PatentStorageID: storageID,
	}

	return patent, nil
}

func getInventorNames(inventors []struct {
	Addressbook struct {
		FirstName struct {
			Text string `xml:",chardata"`
		}
		LastName struct {
			Text string `xml:",chardata"`
		}
	}
}) []string {
	var names []string
	for _, inventor := range inventors {
		names = append(names, inventor.Addressbook.FirstName.Text+" "+inventor.Addressbook.LastName.Text)
	}
	return names
}

func getAssigneeName(assignee Assignee) string {
	return assignee.Name
}

func getClassification(classificationsCpc struct {
	MainCpc struct {
		ClassificationCpc struct {
			Section struct {
				Text string `xml:",chardata"`
			}
			Class struct {
				Text string `xml:",chardata"`
			}
			Subclass struct {
				Text string `xml:",chardata"`
			}
		} `xml:"classification-cpc"`
	} `xml:"main-cpc"`
}) Classification {
	return Classification{
		Section:  classificationsCpc.MainCpc.ClassificationCpc.Section.Text,
		Class:    classificationsCpc.MainCpc.ClassificationCpc.Class.Text,
		Subclass: classificationsCpc.MainCpc.ClassificationCpc.Subclass.Text,
	}
}

type Inventor struct {
	FirstName string
	LastName  string
}

type Assignee struct {
	Name string
}

type ClassificationsCpc struct {
	Classification string
}

type Classification struct {
	Section  string
	Class    string
	Subclass string
}
