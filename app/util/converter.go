package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

func ConvertToSpaced(input string) string {
	// Regular expression to match uppercase letters (except the first one)
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	// Replace matches with a space between the letters
	spaced := re.ReplaceAllString(input, `${1} ${2}`)
	// Convert the first letter to uppercase and the rest to lowercase
	return strings.ToUpper(spaced[:1]) + strings.ToLower(spaced[1:])
}

func ExtractIDFromReference(input string) (int64, error) {

	lastDashIndex := strings.LastIndex(input, "-")
	if lastDashIndex == -1 {
		return 0, fmt.Errorf("invalid reference ID format: %s", input)
	}

	// Ambil substring setelah tanda '-' terakhir
	userIDStr := input[lastDashIndex+1:]

	// Konversi substring ke int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func CurrencyFormat(amount float64) string {
	humanizeValue := humanize.CommafWithDigits(amount, 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	return "Rp " + stringValue
}
