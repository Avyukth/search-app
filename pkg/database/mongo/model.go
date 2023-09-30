package mongo

import (
	"encoding/xml"
	"time"

	"github.com/avyukth/search-app/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	LinkHash  string             `bson:"linkHash"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Patent struct {
	PatentTitle     string   `bson:"patentTitle"`
	PatentNumber    string   `bson:"patentNumber"`
	InventorNames   []string `bson:"inventorNames"`
	AssigneeName    string   `bson:"assigneeName"`
	ApplicationDate string   `bson:"applicationDate"`
	IssueDate       string   `bson:"issueDate"`
	DesignClass     string   `bson:"designClass,omitempty"`
	PatentStorageID string   `bson:"patentStorageID"`
}

type Index struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PatentObj Patent             `bson:"patentObj"`
}

type Database struct {
	Client     *mongo.Client
	Config     *config.Config
	Collection *mongo.Collection
}

type UsPatentGrant struct {
	XMLName                  xml.Name `xml:"us-patent-grant" json:"us-patent-grant,omitempty"`
	Text                     string   `xml:",chardata" json:"text,omitempty"`
	Lang                     string   `xml:"lang,attr" json:"lang,omitempty"`
	DtdVersion               string   `xml:"dtd-version,attr" json:"dtd-version,omitempty"`
	File                     string   `xml:"file,attr" json:"file,omitempty"`
	Status                   string   `xml:"status,attr" json:"status,omitempty"`
	ID                       string   `xml:"id,attr" json:"id,omitempty"`
	Country                  string   `xml:"country,attr" json:"country,omitempty"`
	DateProduced             string   `xml:"date-produced,attr" json:"date-produced,omitempty"`
	DatePubl                 string   `xml:"date-publ,attr" json:"date-publ,omitempty"`
	UsBibliographicDataGrant struct {
		Text                 string `xml:",chardata" json:"text,omitempty"`
		PublicationReference struct {
			Text       string `xml:",chardata" json:"text,omitempty"`
			DocumentID struct {
				Text    string `xml:",chardata" json:"text,omitempty"`
				Country struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"country" json:"country,omitempty"`
				DocNumber struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"doc-number" json:"doc-number,omitempty"`
				Kind struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"kind" json:"kind,omitempty"`
				Date struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"date" json:"date,omitempty"`
			} `xml:"document-id" json:"document-id,omitempty"`
		} `xml:"publication-reference" json:"publication-reference,omitempty"`
		ApplicationReference struct {
			Text       string `xml:",chardata" json:"text,omitempty"`
			ApplType   string `xml:"appl-type,attr" json:"appl-type,omitempty"`
			DocumentID struct {
				Text    string `xml:",chardata" json:"text,omitempty"`
				Country struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"country" json:"country,omitempty"`
				DocNumber struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"doc-number" json:"doc-number,omitempty"`
				Date struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"date" json:"date,omitempty"`
			} `xml:"document-id" json:"document-id,omitempty"`
		} `xml:"application-reference" json:"application-reference,omitempty"`
		UsApplicationSeriesCode struct {
			Text string `xml:",chardata" json:"text,omitempty"`
		} `xml:"us-application-series-code" json:"us-application-series-code,omitempty"`
		UsTermOfGrant struct {
			Text            string `xml:",chardata" json:"text,omitempty"`
			UsTermExtension struct {
				Text string `xml:",chardata" json:"text,omitempty"`
			} `xml:"us-term-extension" json:"us-term-extension,omitempty"`
		} `xml:"us-term-of-grant" json:"us-term-of-grant,omitempty"`
		ClassificationsIpcr struct {
			Text               string `xml:",chardata" json:"text,omitempty"`
			ClassificationIpcr []struct {
				Text                string `xml:",chardata" json:"text,omitempty"`
				IpcVersionIndicator struct {
					Text string `xml:",chardata" json:"text,omitempty"`
					Date struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"date" json:"date,omitempty"`
				} `xml:"ipc-version-indicator" json:"ipc-version-indicator,omitempty"`
				ClassificationLevel struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"classification-level" json:"classification-level,omitempty"`
				Section struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"section" json:"section,omitempty"`
				Class struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"class" json:"class,omitempty"`
				Subclass struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"subclass" json:"subclass,omitempty"`
				MainGroup struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"main-group" json:"main-group,omitempty"`
				Subgroup struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"subgroup" json:"subgroup,omitempty"`
				SymbolPosition struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"symbol-position" json:"symbol-position,omitempty"`
				ClassificationValue struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"classification-value" json:"classification-value,omitempty"`
				ActionDate struct {
					Text string `xml:",chardata" json:"text,omitempty"`
					Date struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"date" json:"date,omitempty"`
				} `xml:"action-date" json:"action-date,omitempty"`
				GeneratingOffice struct {
					Text    string `xml:",chardata" json:"text,omitempty"`
					Country struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"country" json:"country,omitempty"`
				} `xml:"generating-office" json:"generating-office,omitempty"`
				ClassificationStatus struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"classification-status" json:"classification-status,omitempty"`
				ClassificationDataSource struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"classification-data-source" json:"classification-data-source,omitempty"`
			} `xml:"classification-ipcr" json:"classification-ipcr,omitempty"`
		} `xml:"classifications-ipcr" json:"classifications-ipcr,omitempty"`
		ClassificationsCpc struct {
			Text    string `xml:",chardata" json:"text,omitempty"`
			MainCpc struct {
				Text              string `xml:",chardata" json:"text,omitempty"`
				ClassificationCpc struct {
					Text                string `xml:",chardata" json:"text,omitempty"`
					CpcVersionIndicator struct {
						Text string `xml:",chardata" json:"text,omitempty"`
						Date struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"date" json:"date,omitempty"`
					} `xml:"cpc-version-indicator" json:"cpc-version-indicator,omitempty"`
					Section struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"section" json:"section,omitempty"`
					Class struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"class" json:"class,omitempty"`
					Subclass struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"subclass" json:"subclass,omitempty"`
					MainGroup struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"main-group" json:"main-group,omitempty"`
					Subgroup struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"subgroup" json:"subgroup,omitempty"`
					SymbolPosition struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"symbol-position" json:"symbol-position,omitempty"`
					ClassificationValue struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-value" json:"classification-value,omitempty"`
					ActionDate struct {
						Text string `xml:",chardata" json:"text,omitempty"`
						Date struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"date" json:"date,omitempty"`
					} `xml:"action-date" json:"action-date,omitempty"`
					GeneratingOffice struct {
						Text    string `xml:",chardata" json:"text,omitempty"`
						Country struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"country" json:"country,omitempty"`
					} `xml:"generating-office" json:"generating-office,omitempty"`
					ClassificationStatus struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-status" json:"classification-status,omitempty"`
					ClassificationDataSource struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-data-source" json:"classification-data-source,omitempty"`
					SchemeOriginationCode struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"scheme-origination-code" json:"scheme-origination-code,omitempty"`
				} `xml:"classification-cpc" json:"classification-cpc,omitempty"`
			} `xml:"main-cpc" json:"main-cpc,omitempty"`
			FurtherCpc struct {
				Text              string `xml:",chardata" json:"text,omitempty"`
				ClassificationCpc struct {
					Text                string `xml:",chardata" json:"text,omitempty"`
					CpcVersionIndicator struct {
						Text string `xml:",chardata" json:"text,omitempty"`
						Date struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"date" json:"date,omitempty"`
					} `xml:"cpc-version-indicator" json:"cpc-version-indicator,omitempty"`
					Section struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"section" json:"section,omitempty"`
					Class struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"class" json:"class,omitempty"`
					Subclass struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"subclass" json:"subclass,omitempty"`
					MainGroup struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"main-group" json:"main-group,omitempty"`
					Subgroup struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"subgroup" json:"subgroup,omitempty"`
					SymbolPosition struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"symbol-position" json:"symbol-position,omitempty"`
					ClassificationValue struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-value" json:"classification-value,omitempty"`
					ActionDate struct {
						Text string `xml:",chardata" json:"text,omitempty"`
						Date struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"date" json:"date,omitempty"`
					} `xml:"action-date" json:"action-date,omitempty"`
					GeneratingOffice struct {
						Text    string `xml:",chardata" json:"text,omitempty"`
						Country struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"country" json:"country,omitempty"`
					} `xml:"generating-office" json:"generating-office,omitempty"`
					ClassificationStatus struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-status" json:"classification-status,omitempty"`
					ClassificationDataSource struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"classification-data-source" json:"classification-data-source,omitempty"`
					SchemeOriginationCode struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"scheme-origination-code" json:"scheme-origination-code,omitempty"`
				} `xml:"classification-cpc" json:"classification-cpc,omitempty"`
			} `xml:"further-cpc" json:"further-cpc,omitempty"`
		} `xml:"classifications-cpc" json:"classifications-cpc,omitempty"`
		InventionTitle struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			ID   string `xml:"id,attr" json:"id,omitempty"`
		} `xml:"invention-title" json:"invention-title,omitempty"`
		UsReferencesCited struct {
			Text       string `xml:",chardata" json:"text,omitempty"`
			UsCitation []struct {
				Text   string `xml:",chardata" json:"text,omitempty"`
				Patcit struct {
					Text       string `xml:",chardata" json:"text,omitempty"`
					Num        string `xml:"num,attr" json:"num,omitempty"`
					DocumentID struct {
						Text    string `xml:",chardata" json:"text,omitempty"`
						Country struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"country" json:"country,omitempty"`
						DocNumber struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"doc-number" json:"doc-number,omitempty"`
						Kind struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"kind" json:"kind,omitempty"`
						Name struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"name" json:"name,omitempty"`
						Date struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"date" json:"date,omitempty"`
					} `xml:"document-id" json:"document-id,omitempty"`
				} `xml:"patcit" json:"patcit,omitempty"`
				Category struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"category" json:"category,omitempty"`
				ClassificationCpcText struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"classification-cpc-text" json:"classification-cpc-text,omitempty"`
				ClassificationNational struct {
					Text    string `xml:",chardata" json:"text,omitempty"`
					Country struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"country" json:"country,omitempty"`
					MainClassification struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"main-classification" json:"main-classification,omitempty"`
				} `xml:"classification-national" json:"classification-national,omitempty"`
			} `xml:"us-citation" json:"us-citation,omitempty"`
		} `xml:"us-references-cited" json:"us-references-cited,omitempty"`
		NumberOfClaims struct {
			Text string `xml:",chardata" json:"text,omitempty"`
		} `xml:"number-of-claims" json:"number-of-claims,omitempty"`
		UsExemplaryClaim struct {
			Text string `xml:",chardata" json:"text,omitempty"`
		} `xml:"us-exemplary-claim" json:"us-exemplary-claim,omitempty"`
		UsFieldOfClassificationSearch struct {
			Text                  string `xml:",chardata" json:"text,omitempty"`
			ClassificationCpcText []struct {
				Text string `xml:",chardata" json:"text,omitempty"`
			} `xml:"classification-cpc-text" json:"classification-cpc-text,omitempty"`
		} `xml:"us-field-of-classification-search" json:"us-field-of-classification-search,omitempty"`
		Figures struct {
			Text                  string `xml:",chardata" json:"text,omitempty"`
			NumberOfDrawingSheets struct {
				Text string `xml:",chardata" json:"text,omitempty"`
			} `xml:"number-of-drawing-sheets" json:"number-of-drawing-sheets,omitempty"`
			NumberOfFigures struct {
				Text string `xml:",chardata" json:"text,omitempty"`
			} `xml:"number-of-figures" json:"number-of-figures,omitempty"`
		} `xml:"figures" json:"figures,omitempty"`
		UsRelatedDocuments struct {
			Text               string `xml:",chardata" json:"text,omitempty"`
			RelatedPublication struct {
				Text       string `xml:",chardata" json:"text,omitempty"`
				DocumentID struct {
					Text    string `xml:",chardata" json:"text,omitempty"`
					Country struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"country" json:"country,omitempty"`
					DocNumber struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"doc-number" json:"doc-number,omitempty"`
					Kind struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"kind" json:"kind,omitempty"`
					Date struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"date" json:"date,omitempty"`
				} `xml:"document-id" json:"document-id,omitempty"`
			} `xml:"related-publication" json:"related-publication,omitempty"`
		} `xml:"us-related-documents" json:"us-related-documents,omitempty"`
		UsParties struct {
			Text         string `xml:",chardata" json:"text,omitempty"`
			UsApplicants struct {
				Text        string `xml:",chardata" json:"text,omitempty"`
				UsApplicant struct {
					Text                       string `xml:",chardata" json:"text,omitempty"`
					Sequence                   string `xml:"sequence,attr" json:"sequence,omitempty"`
					AppType                    string `xml:"app-type,attr" json:"app-type,omitempty"`
					Designation                string `xml:"designation,attr" json:"designation,omitempty"`
					ApplicantAuthorityCategory string `xml:"applicant-authority-category,attr" json:"applicant-authority-category,omitempty"`
					Addressbook                struct {
						Text    string `xml:",chardata" json:"text,omitempty"`
						Orgname struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"orgname" json:"orgname,omitempty"`
						Address struct {
							Text string `xml:",chardata" json:"text,omitempty"`
							City struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"city" json:"city,omitempty"`
							State struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"state" json:"state,omitempty"`
							Country struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"country" json:"country,omitempty"`
						} `xml:"address" json:"address,omitempty"`
					} `xml:"addressbook" json:"addressbook,omitempty"`
					Residence struct {
						Text    string `xml:",chardata" json:"text,omitempty"`
						Country struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"country" json:"country,omitempty"`
					} `xml:"residence" json:"residence,omitempty"`
				} `xml:"us-applicant" json:"us-applicant,omitempty"`
			} `xml:"us-applicants" json:"us-applicants,omitempty"`
			Inventors struct {
				Text     string `xml:",chardata" json:"text,omitempty"`
				Inventor []struct {
					Text        string `xml:",chardata" json:"text,omitempty"`
					Sequence    string `xml:"sequence,attr" json:"sequence,omitempty"`
					Designation string `xml:"designation,attr" json:"designation,omitempty"`
					Addressbook struct {
						Text     string `xml:",chardata" json:"text,omitempty"`
						LastName struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"last-name" json:"last-name,omitempty"`
						FirstName struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"first-name" json:"first-name,omitempty"`
						Address struct {
							Text string `xml:",chardata" json:"text,omitempty"`
							City struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"city" json:"city,omitempty"`
							State struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"state" json:"state,omitempty"`
							Country struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"country" json:"country,omitempty"`
						} `xml:"address" json:"address,omitempty"`
					} `xml:"addressbook" json:"addressbook,omitempty"`
				} `xml:"inventor" json:"inventor,omitempty"`
			} `xml:"inventors" json:"inventors,omitempty"`
			Agents struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				Agent []struct {
					Text        string `xml:",chardata" json:"text,omitempty"`
					Sequence    string `xml:"sequence,attr" json:"sequence,omitempty"`
					RepType     string `xml:"rep-type,attr" json:"rep-type,omitempty"`
					Addressbook struct {
						Text     string `xml:",chardata" json:"text,omitempty"`
						LastName struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"last-name" json:"last-name,omitempty"`
						FirstName struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"first-name" json:"first-name,omitempty"`
						Address struct {
							Text    string `xml:",chardata" json:"text,omitempty"`
							Country struct {
								Text string `xml:",chardata" json:"text,omitempty"`
							} `xml:"country" json:"country,omitempty"`
						} `xml:"address" json:"address,omitempty"`
					} `xml:"addressbook" json:"addressbook,omitempty"`
				} `xml:"agent" json:"agent,omitempty"`
			} `xml:"agents" json:"agents,omitempty"`
		} `xml:"us-parties" json:"us-parties,omitempty"`
		Assignees struct {
			Text     string `xml:",chardata" json:"text,omitempty"`
			Assignee struct {
				Text        string `xml:",chardata" json:"text,omitempty"`
				Addressbook struct {
					Text    string `xml:",chardata" json:"text,omitempty"`
					Orgname struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"orgname" json:"orgname,omitempty"`
					Role struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"role" json:"role,omitempty"`
					Address struct {
						Text string `xml:",chardata" json:"text,omitempty"`
						City struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"city" json:"city,omitempty"`
						State struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"state" json:"state,omitempty"`
						Country struct {
							Text string `xml:",chardata" json:"text,omitempty"`
						} `xml:"country" json:"country,omitempty"`
					} `xml:"address" json:"address,omitempty"`
				} `xml:"addressbook" json:"addressbook,omitempty"`
			} `xml:"assignee" json:"assignee,omitempty"`
		} `xml:"assignees" json:"assignees,omitempty"`
		Examiners struct {
			Text            string `xml:",chardata" json:"text,omitempty"`
			PrimaryExaminer struct {
				Text     string `xml:",chardata" json:"text,omitempty"`
				LastName struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"last-name" json:"last-name,omitempty"`
				FirstName struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"first-name" json:"first-name,omitempty"`
				Department struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"department" json:"department,omitempty"`
			} `xml:"primary-examiner" json:"primary-examiner,omitempty"`
			AssistantExaminer struct {
				Text     string `xml:",chardata" json:"text,omitempty"`
				LastName struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"last-name" json:"last-name,omitempty"`
				FirstName struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"first-name" json:"first-name,omitempty"`
			} `xml:"assistant-examiner" json:"assistant-examiner,omitempty"`
		} `xml:"examiners" json:"examiners,omitempty"`
	} `xml:"us-bibliographic-data-grant" json:"us-bibliographic-data-grant,omitempty"`
	Abstract struct {
		Text string `xml:",chardata" json:"text,omitempty"`
		ID   string `xml:"id,attr" json:"id,omitempty"`
		P    struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			ID   string `xml:"id,attr" json:"id,omitempty"`
			Num  string `xml:"num,attr" json:"num,omitempty"`
		} `xml:"p" json:"p,omitempty"`
	} `xml:"abstract" json:"abstract,omitempty"`
	Drawings struct {
		Text   string `xml:",chardata" json:"text,omitempty"`
		ID     string `xml:"id,attr" json:"id,omitempty"`
		Figure []struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			ID   string `xml:"id,attr" json:"id,omitempty"`
			Num  string `xml:"num,attr" json:"num,omitempty"`
			Img  struct {
				Text        string `xml:",chardata" json:"text,omitempty"`
				ID          string `xml:"id,attr" json:"id,omitempty"`
				He          string `xml:"he,attr" json:"he,omitempty"`
				Wi          string `xml:"wi,attr" json:"wi,omitempty"`
				File        string `xml:"file,attr" json:"file,omitempty"`
				Alt         string `xml:"alt,attr" json:"alt,omitempty"`
				ImgContent  string `xml:"img-content,attr" json:"img-content,omitempty"`
				ImgFormat   string `xml:"img-format,attr" json:"img-format,omitempty"`
				Orientation string `xml:"orientation,attr" json:"orientation,omitempty"`
			} `xml:"img" json:"img,omitempty"`
		} `xml:"figure" json:"figure,omitempty"`
	} `xml:"drawings" json:"drawings,omitempty"`
	Description struct {
		Text    string `xml:",chardata" json:"text,omitempty"`
		ID      string `xml:"id,attr" json:"id,omitempty"`
		Heading []struct {
			Text  string `xml:",chardata" json:"text,omitempty"`
			ID    string `xml:"id,attr" json:"id,omitempty"`
			Level string `xml:"level,attr" json:"level,omitempty"`
		} `xml:"heading" json:"heading,omitempty"`
		P []struct {
			Text   string `xml:",chardata" json:"text,omitempty"`
			ID     string `xml:"id,attr" json:"id,omitempty"`
			Num    string `xml:"num,attr" json:"num,omitempty"`
			Figref []struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				Idref string `xml:"idref,attr" json:"idref,omitempty"`
				B     []struct {
					Text string `xml:",chardata" json:"text,omitempty"`
				} `xml:"b" json:"b,omitempty"`
			} `xml:"figref" json:"figref,omitempty"`
			B []struct {
				Text string `xml:",chardata" json:"text,omitempty"`
			} `xml:"b" json:"b,omitempty"`
		} `xml:"p" json:"p,omitempty"`
		DescriptionOfDrawings struct {
			Text    string `xml:",chardata" json:"text,omitempty"`
			Heading struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				ID    string `xml:"id,attr" json:"id,omitempty"`
				Level string `xml:"level,attr" json:"level,omitempty"`
			} `xml:"heading" json:"heading,omitempty"`
			P []struct {
				Text   string `xml:",chardata" json:"text,omitempty"`
				ID     string `xml:"id,attr" json:"id,omitempty"`
				Num    string `xml:"num,attr" json:"num,omitempty"`
				Figref []struct {
					Text  string `xml:",chardata" json:"text,omitempty"`
					Idref string `xml:"idref,attr" json:"idref,omitempty"`
					B     struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"b" json:"b,omitempty"`
				} `xml:"figref" json:"figref,omitempty"`
			} `xml:"p" json:"p,omitempty"`
		} `xml:"description-of-drawings" json:"description-of-drawings,omitempty"`
	} `xml:"description" json:"description,omitempty"`
	UsClaimStatement struct {
		Text string `xml:",chardata" json:"text,omitempty"`
	} `xml:"us-claim-statement" json:"us-claim-statement,omitempty"`
	Claims struct {
		Text  string `xml:",chardata" json:"text,omitempty"`
		ID    string `xml:"id,attr" json:"id,omitempty"`
		Claim []struct {
			Text      string `xml:",chardata" json:"text,omitempty"`
			ID        string `xml:"id,attr" json:"id,omitempty"`
			Num       string `xml:"num,attr" json:"num,omitempty"`
			ClaimText struct {
				Text      string `xml:",chardata" json:"text,omitempty"`
				ClaimText []struct {
					Text      string `xml:",chardata" json:"text,omitempty"`
					ClaimText []struct {
						Text string `xml:",chardata" json:"text,omitempty"`
					} `xml:"claim-text" json:"claim-text,omitempty"`
				} `xml:"claim-text" json:"claim-text,omitempty"`
				ClaimRef struct {
					Text  string `xml:",chardata" json:"text,omitempty"`
					Idref string `xml:"idref,attr" json:"idref,omitempty"`
				} `xml:"claim-ref" json:"claim-ref,omitempty"`
			} `xml:"claim-text" json:"claim-text,omitempty"`
		} `xml:"claim" json:"claim,omitempty"`
	} `xml:"claims" json:"claims,omitempty"`
}
