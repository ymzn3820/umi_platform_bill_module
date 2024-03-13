package models

type Product struct {
	ProdId            int     `json:"prod_id"`
	ProdName          string  `json:"prod_name"`
	ProdPrice         float32 `json:"modified_by"`
	ProdCateId        int     `json:"prod_cate_id"`
	ValidPeriodDays   int     `json:"valid_period_days"`
	Hashrate          int     `json:"hashrate"`
	DirectedHashrate  int     `json:"directed_hashrate"`
	UniversalHashrate int     `json:"universal_hashrate"`
	IsShow            int     `json:"is_show"`
	IsDelete          int     `json:"is_delete"`
}

func (Product) TableName() string {
	return "pp_products"
}

func GetProductInfoById(prodId int, prodCateId int) (*Product, error) {
	var product Product
	err := db.Debug().Where("prod_id = ? AND prod_cate_id = ? AND is_delete = 0 ", prodId, prodCateId).First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}
