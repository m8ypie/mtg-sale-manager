module github.com/m8ypie/mtg-sale-manager

go 1.25.0

require (
	github.com/m8ypie/mtg-sale-manager/config v0.0.0
	github.com/m8ypie/mtg-sale-manager/db v0.0.0
	github.com/m8ypie/mtg-sale-manager/models v0.0.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/xo/dburl v0.23.8 // indirect
)

replace github.com/m8ypie/mtg-sale-manager/models => ./models

replace github.com/m8ypie/mtg-sale-manager/db => ./db

replace github.com/m8ypie/mtg-sale-manager/config => ./config
