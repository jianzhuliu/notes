package db

var KindSizeMap = map[string]int{
	"int8":      1,
	"uint8":     1,
	"int16":     2,
	"uint16":    2,
	"int32":     4,
	"uint32":    4,
	"float32":   4,
	"float64":   8,
	"int":       8,
	"uint":      8,
	"string":    16,
	"time.Time": 24,
	"[]byte":    24,
	"[]uint8":   24,
}
