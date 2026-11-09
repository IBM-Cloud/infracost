package apiclient

import (
	"testing"

	"github.com/infracost/infracost/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestPricingAPIClient_BuildQuery_CleansIBMRegion(t *testing.T) {
	client := &PricingAPIClient{
		Currency: "USD",
	}

	vendor := "ibm"
	region := "eu-de-1"
	service := "is.instance"
	planName := "cxf-2x4"

	productFilter := &schema.ProductFilter{
		VendorName: &vendor,
		Region:     &region,
		Service:    &service,
		AttributeFilters: []*schema.AttributeFilter{
			{Key: "planName", Value: &planName},
		},
	}

	query := client.buildQuery(productFilter, nil)

	// Original productFilter region should not be mutated
	assert.Equal(t, "eu-de-1", *productFilter.Region)

	// Query productFilter should have cleaned region
	pfMap, ok := query.Variables["productFilter"].(*schema.ProductFilter)
	assert.True(t, ok)
	assert.Equal(t, "eu-de", *pfMap.Region)
}
