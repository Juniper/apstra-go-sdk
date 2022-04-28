package apstraTelemetry

import "strconv"

func intSliceContains(in []int, t int) bool {
	for _, i := range in {
		if i == t {
			return true
		}
	}
	return false
}

func intSliceToStringSlice(in []int) []string {
	var result []string
	for _, i := range in {
		result = append(result, strconv.Itoa(i))
	}
	return result
}
