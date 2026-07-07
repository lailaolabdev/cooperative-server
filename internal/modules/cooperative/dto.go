package cooperative

type UpsertRequest struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"`
	ProvinceCode string `json:"provinceCode" binding:"required"`
	Province     string `json:"province" binding:"required"`
	District     string `json:"district" binding:"required"`
	Village      string `json:"village" binding:"required"`
	Chairman     string `json:"chairman"`
	Phone        string `json:"phone"`
	MemberCount  int    `json:"memberCount" binding:"min=0"`
	Description  string `json:"description"`
	Status       string `json:"status"`
}

func (r UpsertRequest) ValidType() bool {
	return r.Type == TypeAgriculture || r.Type == TypeSCU
}
