package target

import "context"

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context) ([]ProvinceTarget, Summary, error) {
	items, err := s.repository.List(ctx)
	if err != nil {
		return nil, Summary{}, err
	}
	summary := Summary{
		StartYear: 2026, EndYear: 2030, ProvinceCount: len(items),
		StrategicPlanTotal: 34, ProgramPlanTotal: 92,
		AnnualTargets: []AnnualTarget{{Year: 2026, Count: 20}, {Year: 2027, Count: 5}, {Year: 2028, Count: 5}, {Year: 2029, Count: 3}, {Year: 2030, Count: 1}},
		Source:        "ເປົ້າໝາຍສ້າງສະຫະກອນ 2026 - 20230.xlsx",
	}
	for _, item := range items {
		summary.ProvincePipeline += item.TargetCount
		summary.PotentialMembers += item.PotentialMembers
	}
	return items, summary, nil
}
