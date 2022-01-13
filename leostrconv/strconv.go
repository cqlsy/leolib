package leostrconv

import "strconv"

// AtoIDefault0
func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func AtoIDefault(s string, defValue int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defValue
	}
	return i
}

// FormatInt
func FormatInt(i int) string {
	return strconv.Itoa(i)
}

// FormatFloat
func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// FormatFloatPrec0
func FormatFloatPrec0(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}

// FormatFloatPrec0
func FormatFloatPrec2(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// FormatFloatPrec0
func FormatFloatPrec4(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}

// ParseFloat64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseFloat64Default0
func ParseFloat64Default0(s string) float64 {
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return out
}

// IntArrayToStringArr
func IntArrayToStringArr(i []int) []string {
	strArr := make([]string, len(i))
	for k, v := range i {
		strArr[k] = FormatInt(v)
	}
	return strArr
}
