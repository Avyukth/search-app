package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	LinkHash  string             `bson:"linkHash"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Patent struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	PatentTitle     string             `bson:"patentTitle"`
	PatentNumber    string             `bson:"patentNumber"`
	InventorNames   []string           `bson:"inventorNames"`
	AssigneeName    string             `bson:"assigneeName"`
	ApplicationDate string             `bson:"applicationDate"`
	IssueDate       string             `bson:"issueDate"`
	DesignClass     string             `bson:"designClass,omitempty"`
}

type Index struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PatentObj Patent             `bson:"patentObj"`
}
