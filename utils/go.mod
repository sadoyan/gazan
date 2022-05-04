module utils

go 1.17

require (
	configs v0.0.0-00010101000000-000000000000
)

replace configs => ../configs
require github.com/golang-jwt/jwt v3.2.2+incompatible
