package target

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProvinceTarget struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProvinceCode     string             `bson:"provinceCode" json:"provinceCode"`
	Province         string             `bson:"province" json:"province"`
	TargetCount      int                `bson:"targetCount" json:"targetCount"`
	PotentialMembers int                `bson:"potentialMembers" json:"potentialMembers"`
	Districts        []string           `bson:"districts" json:"districts"`
	Groups           []TargetGroup      `bson:"groups" json:"groups"`
	StartYear        int                `bson:"startYear" json:"startYear"`
	EndYear          int                `bson:"endYear" json:"endYear"`
	Source           string             `bson:"source" json:"source"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type TargetGroup struct {
	Name             string `bson:"name" json:"name"`
	District         string `bson:"district" json:"district"`
	MemberCount      int    `bson:"memberCount" json:"memberCount"`
	CooperativeCount int    `bson:"cooperativeCount" json:"cooperativeCount"`
	Activity         string `bson:"activity" json:"activity"`
	Capital          string `bson:"capital" json:"capital"`
	Implementer      string `bson:"implementer" json:"implementer"`
}

type AnnualTarget struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

type Summary struct {
	StartYear          int            `json:"startYear"`
	EndYear            int            `json:"endYear"`
	ProvincePipeline   int            `json:"provincePipeline"`
	PotentialMembers   int            `json:"potentialMembers"`
	ProvinceCount      int            `json:"provinceCount"`
	StrategicPlanTotal int            `json:"strategicPlanTotal"`
	ProgramPlanTotal   int            `json:"programPlanTotal"`
	AnnualTargets      []AnnualTarget `json:"annualTargets"`
	Source             string         `json:"source"`
}
