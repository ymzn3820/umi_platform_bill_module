module github.com/EDDYCJY/go-gin-example

go 1.20

require (
	github.com/360EntSecGroup-Skylar/excelize v1.3.1-0.20180527032555-9e463b461434
	github.com/astaxie/beego v1.9.3-0.20171218111859-f16688817aa4
	github.com/boombuler/barcode v1.0.1-0.20180315051053-3c06908149f7
	github.com/dgrijalva/jwt-go v3.1.0+incompatible
	github.com/gin-gonic/gin v1.9.0
	github.com/go-ini/ini v1.32.1-0.20180214101753-32e4be5f41bb
	github.com/gomodule/redigo v2.0.1-0.20180401191855-9352ab68be13+incompatible
	github.com/jinzhu/gorm v0.0.0-20180213101209-6e1387b44c64
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
	github.com/tealeg/xlsx v1.0.4-0.20180419195153-f36fa3be8893
	github.com/unknwon/com v1.0.1
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190920000552-128d9f4ae1cd // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/go-sql-driver/mysql v1.4.1-0.20190510102335-877a9775f068 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/jinzhu/inflection v0.0.0-20170102125226-1c35d901db3d // indirect
	github.com/jinzhu/now v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/appengine v1.6.3 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.47.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)


replace (
	github.com/EDDYCJY/go-gin-example/pkg/e => ./pkg/e
	github.com/EDDYCJY/go-gin-example/models => ./models
	github.com/EDDYCJY/go-gin-example/pkg/logging => ./pkg/logging
	github.com/EDDYCJY/go-gin-example/pkg/setting => ./pkg/setting
	github.com/EDDYCJY/go-gin-example/pkg/app => ./pkg/app
	github.com/EDDYCJY/go-gin-example/service/product_service => ./service/product_service
	github.com/EDDYCJY/go-gin-example/pkg/gredis => ./pkg/gredis
	github.com/EDDYCJY/go-gin-example/pkg/file => ./pkg/file
	github.com/EDDYCJY/go-gin-example/pkg/util => ./pkg/util
	github.com/EDDYCJY/go-gin-example/service/tag_service => ./service/tag_service
	github.com/EDDYCJY/go-gin-example/service/auth_service => ./service/auth_service
	github.com/EDDYCJY/go-gin-example/pkg/upload => ./pkg/upload
	github.com/EDDYCJY/go-gin-example/middleware/cors => ./middleware/cors
	github.com/EDDYCJY/go-gin-example/docs => ./docs
	github.com/EDDYCJY/go-gin-example/pkg/export => ./pkg/export
	github.com/EDDYCJY/go-gin-example/routers/api => ./routers/api
)