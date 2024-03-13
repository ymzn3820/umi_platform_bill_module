package v1

import (
	"github.com/EDDYCJY/go-gin-example/service/product_service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

// @Summary	Get Product Information By Product Id
// @Produce	json
// @Param		prod_id	  path		int	false	"prod_id"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/product/{prod_id} [get]
func GetProductById(c *gin.Context) {
	appG := app.Gin{C: c}
	prodID := c.Param("prod_id")             // Get the prod_id from the URL query parameters
	prodIDInt := com.StrTo(prodID).MustInt() // Convert prod_id to int
	if prodIDInt == 0 && prodID != "0" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	productService := product_service.Product{ProdId: prodIDInt}
	product, err := productService.GetProductById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_PRODUCT_BY_ID, nil)
		return
	}

	if product.ProdId == 0 {
		appG.Response(http.StatusNotFound, e.ERROR_GET_PRODUCT_BY_ID, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}
