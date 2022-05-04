module gazan

go 1.17

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	gopkg.in/ini.v1 v1.66.4
)

replace configs => ./configs

replace mainfiles => ./mainfiles

replace utils => ./utils

require (
	configs v0.0.0-00010101000000-000000000000 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	mainfiles v0.0.0-00010101000000-000000000000 // indirect
	utils v0.0.0-00010101000000-000000000000 // indirect
)
