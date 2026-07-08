package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"cooperative-service/internal/config"
	"cooperative-service/internal/database"
	targetmodule "cooperative-service/internal/modules/target"
)

const source = "ເປົ້າໝາຍສ້າງສະຫະກອນ 2026 - 20230.xlsx"

//go:embed targets.json
var targetData []byte

type sourceGroup struct {
	ProvinceCode     string `json:"provinceCode"`
	Province         string `json:"province"`
	Name             string `json:"name"`
	District         string `json:"district"`
	MemberCount      int    `json:"memberCount"`
	CooperativeCount int    `json:"cooperativeCount"`
	Activity         string `json:"activity"`
	Capital          string `json:"capital"`
	Implementer      string `json:"implementer"`
}

func aggregate(records []sourceGroup) []targetmodule.ProvinceTarget {
	byProvince := make(map[string]*targetmodule.ProvinceTarget)
	for _, record := range records {
		item := byProvince[record.ProvinceCode]
		if item == nil {
			item = &targetmodule.ProvinceTarget{ProvinceCode: record.ProvinceCode, Province: record.Province, Districts: make([]string, 0), Groups: make([]targetmodule.TargetGroup, 0)}
			byProvince[record.ProvinceCode] = item
		}
		item.TargetCount += record.CooperativeCount
		item.PotentialMembers += record.MemberCount
		if !contains(item.Districts, record.District) {
			item.Districts = append(item.Districts, record.District)
		}
		item.Groups = append(item.Groups, targetmodule.TargetGroup{Name: record.Name, District: record.District, MemberCount: record.MemberCount, CooperativeCount: record.CooperativeCount, Activity: record.Activity, Capital: record.Capital, Implementer: record.Implementer})
	}
	items := make([]targetmodule.ProvinceTarget, 0, len(byProvince))
	for _, item := range byProvince {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ProvinceCode < items[j].ProvinceCode })
	return items
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func main() {
	var records []sourceGroup
	if err := json.Unmarshal(targetData, &records); err != nil {
		log.Fatal(err)
	}
	targets := aggregate(records)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	repository := targetmodule.NewRepository(client.Database(cfg.Database))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = repository.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	now := time.Now().UTC()
	for _, item := range targets {
		item.StartYear, item.EndYear, item.Source, item.UpdatedAt = 2026, 2030, source, now
		if err = repository.Upsert(ctx, item); err != nil {
			log.Fatalf("seed %s: %v", item.ProvinceCode, err)
		}
	}
	fmt.Printf("Seeded %d target groups across %d provinces successfully.\n", len(records), len(targets))
}
