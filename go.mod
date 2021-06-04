module github.com/tossp/tsgo

go 1.16

require (
	github.com/allegro/bigcache v1.2.1
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/casbin/casbin/v2 v2.16.0
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getsentry/sentry-go v0.11.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-xorm/xorm v0.7.9
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200209183636-89e6cbcd0b6d // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf
	github.com/jackc/pgconn v1.8.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgtype v1.6.1
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/inflection v1.0.0
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo/v4 v4.3.0
	github.com/lib/pq v1.8.0
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/minio/minio v0.0.0-20201103204752-b9277c803098
	github.com/minio/minio-go/v6 v6.0.58-0.20200612001654-a57fec8037ec
	github.com/minio/minio-go/v7 v7.0.10 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/montanaflynn/stats v0.6.3 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/secure-io/sio-go v0.3.1 // indirect
	github.com/shirou/gopsutil v3.20.10+incompatible // indirect
	github.com/shopspring/decimal v0.0.0-20200227202807-02e2044944cc // indirect
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/tjfoc/gmsm v1.3.2
	github.com/ugorji/go v1.1.13 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/gorm v1.21.10 // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
	xorm.io/builder v0.3.7 // indirect
	xorm.io/core v0.7.3 // indirect
)

replace (
	github.com/jackc/pgtype v1.6.1 => github.com/tossp/pgtype v1.6.2-0.20201126104256-ff11ce768d3d
	github.com/pdfcpu/pdfcpu v0.3.11 => github.com/tossp/pdfcpu v0.3.12-0.20210428151629-ffe29d5d1606
)
