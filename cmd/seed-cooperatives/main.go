package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cooperative-service/internal/config"
	"cooperative-service/internal/database"
	coopmodule "cooperative-service/internal/modules/cooperative"
)

type sourceRecord struct {
	No               int    `json:"no"`
	Name             string `json:"cooperative_name"`
	Village          string `json:"village"`
	District         string `json:"district"`
	Province         string `json:"province"`
	ProductionAreaHa any    `json:"production_area_ha"`
	MembersCount     int    `json:"members_count"`
	EstablishedYear  int    `json:"established_year"`
	ChairmanName     string `json:"chairman_name"`
	Phone            string `json:"phone"`
}

var provinceCodes = map[string]string{
	"ຊຽງຂວາງ": "XK", "ຫຼວງນ້ຳທາ": "LN", "ຈຳປາສັກ": "CH", "ຈໍາປາສັກ": "CH",
	"ໄຊຍະບູລີ": "XY", "ຫົວພັນ": "HP", "ນະຄອນຫຼວງ": "VC", "ບໍ່ແກ້ວ": "BK",
	"ວຽງຈັນ": "VT", "ຜັ້ງສາລີ": "PH", "ເຊກອງ": "SK", "ສາລະວັນ": "SL",
	"ຄຳມ່ວນ": "KM", "ອັດຕະປື": "AT", "ໄຊສົມບູນ": "XS", "ສະຫວັນນະເຂດ": "SV",
	"ສະຫັວນນະເຂດ": "SV",
}

func optionalFloat(value any) *float64 {
	var number float64
	switch typed := value.(type) {
	case float64:
		number = typed
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return nil
		}
		number = parsed
	default:
		return nil
	}
	return &number
}

func main() {
	path := flag.String("file", "../cooperatives_58.json", "seed JSON file")
	flag.Parse()
	data, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	var records []sourceRecord
	if err = json.Unmarshal(data, &records); err != nil {
		log.Fatal(err)
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	repository := coopmodule.NewRepository(client.Database(cfg.Database))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	now := time.Now().UTC()
	seeded := 0
	for _, record := range records {
		province := strings.TrimSpace(record.Province)
		code, ok := provinceCodes[province]
		if !ok {
			log.Printf("skip no=%d: unknown province %q", record.No, province)
			continue
		}
		item := coopmodule.Cooperative{Name: record.Name, Type: coopmodule.TypeAgriculture, ProvinceCode: code, Province: province, District: record.District, Village: record.Village, Chairman: record.ChairmanName, Phone: record.Phone, MemberCount: record.MembersCount, Status: "active", Source: "cooperatives_58", SourceNo: record.No, ProductionAreaHa: optionalFloat(record.ProductionAreaHa), EstablishedYear: record.EstablishedYear, CreatedAt: now, UpdatedAt: now}
		if err = repository.UpsertSeed(ctx, item); err != nil {
			log.Fatalf("seed no=%d: %v", record.No, err)
		}
		seeded++
	}
	fmt.Printf("Seeded %d cooperative records successfully.\n", seeded)
}
