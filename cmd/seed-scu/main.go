package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"cooperative-service/internal/config"
	"cooperative-service/internal/database"
	coopmodule "cooperative-service/internal/modules/cooperative"
	"github.com/xuri/excelize/v2"
)

type provinceInfo struct{ Code, LaoName string }

var provinces = map[string]provinceInfo{
	"savannakhet": {"SV", "ສະຫວັນນະເຂດ"}, "lungphabang": {"LP", "ຫຼວງພະບາງ"},
	"champasak": {"CH", "ຈຳປາສັກ"}, "vientiane capital": {"VC", "ນະຄອນຫຼວງວຽງຈັນ"},
	"saravan": {"SL", "ສາລະວັນ"}, "lungnamtha": {"LN", "ຫຼວງນ້ຳທາ"},
	"bolikhamxay": {"BL", "ບໍລິຄຳໄຊ"}, "houaphan": {"HP", "ຫົວພັນ"},
	"attapeu": {"AT", "ອັດຕະປື"}, "xayabuly": {"XY", "ໄຊຍະບູລີ"},
	"xeingkhoung": {"XK", "ຊຽງຂວາງ"}, "vientiane": {"VT", "ວຽງຈັນ"},
	"khammuan": {"KM", "ຄຳມ່ວນ"}, "bokeo": {"BK", "ບໍ່ແກ້ວ"},
}

func cell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func main() {
	path := flag.String("file", "../LASCU members and SCU in Laos 2025.xlsx", "SCU Excel file")
	flag.Parse()
	book, err := excelize.OpenFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer book.Close()
	rows, err := book.GetRows("List of SCU 2023")
	if err != nil {
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
	deleted, err := repository.DeleteByTypes(ctx, []string{"saving_credit", "service"})
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().UTC()
	seeded := 0
	for rowIndex, row := range rows {
		if rowIndex < 5 {
			continue
		}
		var sourceNo int
		if _, err = fmt.Sscanf(cell(row, 0), "%d", &sourceNo); err != nil || sourceNo < 1 || sourceNo > 32 {
			continue
		}
		name := cell(row, 1)
		if name == "" {
			continue
		}
		province, ok := provinces[strings.ToLower(cell(row, 4))]
		if !ok {
			log.Printf("skip row %d: unknown province %q", rowIndex+1, cell(row, 4))
			continue
		}
		item := coopmodule.Cooperative{Name: name, Type: coopmodule.TypeSCU, ProvinceCode: province.Code, Province: province.LaoName, District: cell(row, 3), Village: cell(row, 2), Phone: cell(row, 5), MemberCount: 0, Status: "active", Source: "lascu_scu_2025", SourceNo: sourceNo, CreatedAt: now, UpdatedAt: now}
		if err = repository.UpsertSeed(ctx, item); err != nil {
			log.Fatalf("seed row %d: %v", rowIndex+1, err)
		}
		seeded++
	}
	fmt.Printf("Seeded %d SCU records; removed %d deprecated records.\n", seeded, deleted)
}
