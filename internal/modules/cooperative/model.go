package cooperative

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	TypeAgriculture = "agriculture"
	TypeSCU         = "scu"
)

type Cooperative struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Type             string             `bson:"type" json:"type"`
	ProvinceCode     string             `bson:"provinceCode" json:"provinceCode"`
	Province         string             `bson:"province" json:"province"`
	District         string             `bson:"district" json:"district"`
	Village          string             `bson:"village" json:"village"`
	Chairman         string             `bson:"chairman" json:"chairman"`
	Phone            string             `bson:"phone" json:"phone"`
	MemberCount      int                `bson:"memberCount" json:"memberCount"`
	Description      string             `bson:"description" json:"description"`
	Status           string             `bson:"status" json:"status"`
	Source           string             `bson:"source,omitempty" json:"source,omitempty"`
	SourceNo         int                `bson:"sourceNo,omitempty" json:"sourceNo,omitempty"`
	ProductionAreaHa *float64           `bson:"productionAreaHa,omitempty" json:"productionAreaHa,omitempty"`
	EstablishedYear  int                `bson:"establishedYear,omitempty" json:"establishedYear,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
