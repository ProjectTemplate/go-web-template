package utils

func FillZero(str string, totalLength int) string {
	fillStr := ""
	for i := len(str) + 1; i <= totalLength; i++ {
		fillStr = fillStr + "0"
	}
	return fillStr + str
}
