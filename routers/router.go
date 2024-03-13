package routers

import (
	"github.com/EDDYCJY/go-gin-example/middleware/cors"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/EDDYCJY/go-gin-example/pkg/export"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/EDDYCJY/go-gin-example/routers/api"
	"github.com/EDDYCJY/go-gin-example/routers/api/v1"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(cors.CORSMiddleware())
	//apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		apiv1.POST("/tags/export", v1.ExportTag)
		//导入标签
		apiv1.GET("/product/:prod_id", v1.GetProductById)
		apiv1.POST("/hashrate", v1.HashRateInByProdId)
		apiv1.PUT("/hashrate", v1.ConsumeHashRate)
		apiv1.GET("/hashrate/:user_id", v1.GetHashRate)
		apiv1.POST("/renew", v1.RenewHashRate)
		apiv1.POST("/active", v1.CheckUserIsActive)
		apiv1.POST("/complimentary", v1.ComplimentaryHashRste)
	}

	return r
}
