module mainfiles

require (
	configs v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

replace configs => ../configs
replace utils => ../utils

go 1.17
