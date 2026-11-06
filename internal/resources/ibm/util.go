package ibm

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

func CleanRegion(region string) string {
	parts := strings.Split(region, "-")
	if len(parts) > 2 {
		return strings.Join(parts[:2], "-")
	}
	return region
}

func strPtr(s string) *string {
	return &s
}

func regexPtr(regex string) *string {
	return strPtr(fmt.Sprintf("/%s/i", regex))
}

func decimalPtr(d decimal.Decimal) *decimal.Decimal {
	return &d
}
