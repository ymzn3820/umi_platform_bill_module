package v1

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/product_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

// HashRateInByProdId @Summary	写入用户算力值
// @Produce	json
// @Param		user_id	   query		string	false	"用户id"
// @Param		prod_id	   query		int64	false	"产品id"
// @Param		quantity	query		int64	false	"数量"
// @Param		prod_cate_id	query		int64	false	"种类"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/hashrate [post]
func HashRateInByProdId(c *gin.Context) {
	appG := app.Gin{C: c}

	userId := c.PostForm("user_id")
	prodId := c.PostForm("prod_id")
	quantity := c.PostForm("quantity")
	prodCate := c.PostForm("prod_cate_id")

	if userId == "" || prodId == "" || prodCate == "" {
		logging.Error("UserId or ProdId is missing")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	prodIdInt, err := strconv.Atoi(prodId)
	if err != nil {
		logging.Error("Invalid ProdId format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	quantityInt, err := strconv.Atoi(quantity) // Convert quantity to int
	if err != nil {
		logging.Error("Invalid quantity format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	prodCateInt, err := strconv.Atoi(prodCate) // Convert quantity to int
	if err != nil {
		logging.Error("Invalid quantity format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	productService := product_service.Product{ProdId: prodIdInt, ProdCateId: prodCateInt}
	product, err := productService.GetProductById()
	logging.Info("Debug: JSON Body = %+v\n", product)
	if err != nil {
		logging.Error("Error retrieving product by ID: %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_GET_PRODUCT_BY_ID, nil)
		return
	}

	// Calculate the hashrate and valid period by multiplying with quantity
	totalHashrateDirected := int64(product.DirectedHashrate) * int64(quantityInt)
	totalUniversalHashrate := int64(product.UniversalHashrate) * int64(quantityInt)
	validPeriod := time.Duration(product.ValidPeriodDays*quantityInt) * 24 * time.Hour

	// Add directed hashrate if any,which one will be called by prod cate id
	if prodCateInt == 3 {
		if totalHashrateDirected > 0 {
			err := gredis.AddHashRateDirected(userId, totalHashrateDirected, validPeriod)
			if err != nil {
				logging.Error("Error adding directed hashrate: %v", err)
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_HASHRATE, nil)
				return
			}
		}

		// Add member universal hashrate if any
		if totalUniversalHashrate > 0 {
			err := gredis.AddHashRateUniversal(userId, totalUniversalHashrate, validPeriod)
			if err != nil {
				logging.Error("Error adding universal hashrate: %v", err)
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_HASHRATE, nil)
				return
			}
		}
	} else {
		// Add package hashrate package if any

		if totalUniversalHashrate > 0 {
			err := gredis.AddHashRatePackage(userId, totalUniversalHashrate, validPeriod)
			if err != nil {
				logging.Error("Error adding universal hashrate: %v", err)
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_HASHRATE, nil)
				return
			}
		}
	}

	logging.Info("Hash rate added successfully for user: %s, product: %d", userId, prodIdInt)
	returnData := map[string]string{"user_id": userId}
	appG.Response(http.StatusOK, e.SUCCESS, returnData)
}

// ConsumeHashRate @Summary	消耗用户算力值
// @Produce	json
// @Param		user_id	   query		string	false	"用户id"
// @Param		hashrate	   query		float64 	false	"消耗的算力值"
// @Param		scene	query		int64	false	"场景，对话为1，其余随便传"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/hashrate [put]
func ConsumeHashRate(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := c.PostForm("user_id")
	hashrate := c.PostForm("hashrate")
	scene := c.PostForm("scene")

	if userId == "" || scene == "" || hashrate == "" {
		logging.Error("UserId or ProdId is missing")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	hashrateInt, err := strconv.ParseInt(hashrate, 10, 64)

	if err != nil {
		logging.Error("Error Consume Universal Hashrate: %v", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	sceneInt, err := strconv.ParseInt(scene, 10, 64)

	if err != nil {
		logging.Error("Error Consume Universal Hashrate: %v", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	var successConsume = 0
	var hasError = true
	// scene 1代表对话场景，扣除定向算力值，如果非对话场景，则先扣流量包的算力，如果扣完了，则扣除
	// 会员里的通用算力值
	if sceneInt != 1 {
		successConsume, hasError = gredis.ConsumeHashRateUniversal(userId, hashrateInt)
	} else {
		successConsume, hasError = gredis.ConsumeHashRateDirected(userId, hashrateInt)
	}

	if !hasError {

		if successConsume == 20000 {

			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		} else if successConsume == 30013 {
			successConsumeUniversal, hasErrorUniversal := gredis.ConsumeHashRatePackage(userId, hashrateInt)

			if !hasErrorUniversal {
				if successConsumeUniversal == 20000 {
					appG.Response(http.StatusOK, e.SUCCESS, nil)
					return
				} else if successConsumeUniversal == 30013 {
					appG.Response(http.StatusOK, e.INSUFFICIENT_HASHRATE, nil)
					return
				}
			} else {
				appG.Response(http.StatusInternalServerError, e.ERROR_CONSUME_HASHRATE, nil)
				return
			}
			appG.Response(http.StatusInternalServerError, successConsume, nil)
			return
		}
	} else {
		appG.Response(http.StatusInternalServerError, successConsume, nil)
		return
	}
}

type HashRateStatistic struct {
	Universal int64 `json:"universal"`
	Directed  int64 `json:"directed"`
}

type HashRateResponse struct {
	HashRateTotals map[string]int64             `json:"hash_rates"`
	Pricing        []gredis.HashRateRulesFields `json:"rules"`
}

// GetHashRate @Summary	汇总用户算力值
// @Param		user_id	   path		string	false	"用户id"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/hashrate/{user_id} [get]
func GetHashRate(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := c.Param("user_id")
	logging.Info("GetHashRate Received Params: %+v\n", c.PostForm)

	hashrates, err := gredis.SummarizeUserHashRates(userId)

	if err != nil {
		logging.Error("Error Get SummarizeUserHashRates: %v", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_HASHRATES, nil)
		return
	}
	pricing, err := gredis.HashRateRules()

	if err != nil {
		logging.Error("Error Get HashRateRules: %v", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_HASHRATES, nil)
		return
	}
	returnData := HashRateResponse{
		HashRateTotals: hashrates,
		Pricing:        pricing,
	}
	appG.Response(http.StatusOK, e.SUCCESS, returnData)
	return
}

// RenewHashRate @Summary	续费
// @Param		user_id	   query		string	false	"用户id"
// @Param		prodId	   query		string	false	"产品id"
// @Param		quantity	   query		string	false	"数量"
// @Param		prodCate	   query		string	false	"种类"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/renew [post]
func RenewHashRate(c *gin.Context) {
	appG := app.Gin{C: c}

	userId := c.PostForm("user_id")
	prodId := c.PostForm("prod_id")
	quantity := c.PostForm("quantity")
	prodCate := c.PostForm("prod_cate_id")

	logging.Info("RenewHashRate Received Params: %+v\n", c.PostForm)
	if userId == "" || prodId == "" || quantity == "" {
		logging.Error("UserId or ProdId or quantity is missing")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	prodIdInt, err := strconv.Atoi(prodId)
	if err != nil {
		logging.Error("Invalid ProdId format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	quantityInt, err := strconv.Atoi(quantity) // Convert quantity to int
	if err != nil {
		logging.Error("Invalid quantity format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}

	prodCateInt, err := strconv.Atoi(prodCate) // Convert prodCateInt to int

	if err != nil {
		logging.Error("Invalid prod cate format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}
	productService := product_service.Product{ProdId: prodIdInt, ProdCateId: prodCateInt}
	product, err := productService.GetProductById()

	logging.Info("Debug: JSON Body = %+v\n", product)
	if err != nil {
		logging.Error("Error retrieving product by ID: %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_GET_PRODUCT_BY_ID, nil)
		return
	}
	// Calculate the hashrate and valid period by multiplying with quantity
	totalHashrateDirected := int64(product.DirectedHashrate) * int64(quantityInt)
	totalUniversalHashrate := int64(product.UniversalHashrate) * int64(quantityInt)
	validPeriod := time.Duration(product.ValidPeriodDays*quantityInt) * 24 * time.Hour

	if prodCateInt == 3 {
		if totalHashrateDirected > 0 && totalUniversalHashrate > 0 {
			err := gredis.RenewHashRate(userId, totalHashrateDirected, totalUniversalHashrate, validPeriod)
			if err != nil {
				logging.Error("Error RenewHashRate: %v", err)
				appG.Response(http.StatusInternalServerError, e.ERROR_RENEW_HASHRATES, nil)
				return
			}
		} else {
			logging.Error("Error RenewHashRate: %v", err)
			appG.Response(http.StatusInternalServerError, e.ERROR_RENEW_HASHRATES, nil)
			return
		}
	} else {
		if totalUniversalHashrate > 0 {

			err := gredis.RenewHashRate(userId, 0, totalUniversalHashrate, validPeriod)

			if err != nil {
				logging.Error("Error RenewHashRate: %v", err)
				appG.Response(http.StatusInternalServerError, e.ERROR_RENEW_HASHRATES, nil)
				return
			}
		} else {
			logging.Error("Error RenewHashRate: %v", err)
			appG.Response(http.StatusInternalServerError, e.ERROR_RENEW_HASHRATES, nil)
			return
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// CheckUserIsActive @Summary	检查用户是否活跃
// @Param		user_id	   query		string	false	"用户id"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/active [post]
func CheckUserIsActive(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := c.PostForm("user_id")

	isActive, err := gredis.CheckIsActive(userId)
	if err != nil {
		appG.Response(http.StatusOK, e.SUCCESS, isActive)
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, isActive)
		return
	}
}

// ComplimentaryHashRste @Summary	赠送算力值
// @Accept multipart/form-data
// @Param		user_id	   formData		string	false	"用户id"
// @Param		hashrate	   formData		int64	false	"赠送的算力值"
// @Success	20000		{object}	app.Response
// @Failure	500		{object}	app.Response
// @Router		/api/v1/complimentary [post]
func ComplimentaryHashRste(c *gin.Context) {
	appG := app.Gin{C: c}

	userId := c.PostForm("user_id")
	hashrate := c.PostForm("hashrate")
	reason := c.PostForm("reason")

	fmt.Println(reason)
	fmt.Println("reasonreasonreason")
	if userId == "" || hashrate == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	hashrateInt, err := strconv.Atoi(hashrate)
	if err != nil {
		logging.Error("Invalid hashrate format %v", err)
		appG.Response(http.StatusBadRequest, e.ERROR_CONVERT_PARAM, nil)
		return
	}
	validPeriod := time.Duration(365) * 24 * time.Hour

	errAddPackage := gredis.AddHashRatePackage(userId, int64(hashrateInt), validPeriod)
	if errAddPackage != nil {
		logging.Error("Error Add ComplimentaryHashRste", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_HASHRATE, nil)
		return
	}

	_, errAddComplimentary := gredis.ComplimentaryHashRate(userId, int64(hashrateInt), reason)

	if errAddComplimentary != nil {
		if errAddPackage != nil {
			logging.Error("Error Add errAddComplimentary", err)
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_HASHRATE, nil)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
