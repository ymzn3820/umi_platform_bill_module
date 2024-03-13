package product_service

import "github.com/EDDYCJY/go-gin-example/models"

type Product struct {
	ProdId            int     `json:"prod_id"`
	ProdName          string  `json:"prod_name"`
	ProdPrice         float32 `json:"modified_by"`
	ProdCateId        int     `json:"prod_cate_id"`
	ValidPeriodDays   int     `json:"valid_period_days"`
	Hashrate          int     `json:"hashrate"`
	UniversalHashrate int     `json:"universal_hashrate"`
	DirectedHashrate  int     `json:"directed_hashrate"`
	IsShow            int     `json:"is_show"`
	IsDelete          int     `json:"is_delete"`
}

func (p *Product) GetProductById() (Product, error) {
	product, err := models.GetProductInfoById(p.ProdId, p.ProdCateId)
	if err != nil {
		return Product{}, err
	}

	return Product{
		ProdId:            product.ProdId,
		ProdName:          product.ProdName,
		ProdPrice:         product.ProdPrice,
		ProdCateId:        product.ProdCateId,
		ValidPeriodDays:   product.ValidPeriodDays,
		Hashrate:          product.Hashrate,
		DirectedHashrate:  product.DirectedHashrate,
		UniversalHashrate: product.UniversalHashrate,
		IsShow:            product.IsShow,
		IsDelete:          product.IsDelete,
	}, nil
}
