package util

import (
	"fmt"
	"strconv"
)

// Round 截取小数点
func Round(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	if inst, err := strconv.ParseFloat(floatStr, 64); err == nil {
		return inst
	}
	return 0.0
}
