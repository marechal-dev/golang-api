package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var stringBuilder strings.Builder

	alphabetLenght := len(alphabet)

	for i := 0; i < n; i++ {
		carret := alphabet[rand.Intn(alphabetLenght)]
		stringBuilder.WriteByte(carret)
	}

	return stringBuilder.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD", "BRL"}
	currenciesLength := len(currencies)
	return currencies[rand.Intn(currenciesLength)]
}
