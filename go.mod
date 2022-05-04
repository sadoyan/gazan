module gazan

go 1.17

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
)

replace configs => ./configs

replace mainfiles => ./mainfiles

replace utils => ./utils

require (
	configs v0.0.0-00010101000000-000000000000
	mainfiles v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)
