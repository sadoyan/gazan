module gazan

go 1.20

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace configs => ./configs

replace mainfiles => ./mainfiles

replace utils => ./utils

require (
	configs v0.0.0-00010101000000-000000000000
	mainfiles v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)
