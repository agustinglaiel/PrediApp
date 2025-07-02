module prediapp.local

go 1.23.7

require (
	github.com/joho/godotenv v1.5.1
	prediapp.local/db v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	gorm.io/gorm v1.30.0 // indirect
)

replace prediapp.local/db => ./db
