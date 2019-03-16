module daemonw

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20181213202711-891ebc4b82d6
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20181213200352-4d1cda033e06
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20181108054448-85acf8d2951c
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/lib/pq v1.0.0
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.11.0
	github.com/stretchr/testify v1.2.2 // indirect
	github.com/ugorji/go/codec v0.0.0-20181209151446-772ced7fd4c2 // indirect
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
