package ibm_test

import (
	"testing"

	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInstance_BuildResource(t *testing.T) {
	t.Run("regular profile", func(t *testing.T) {
		instance := &ibm.IsInstance{
			Address: "ibm_is_instance.test",
			Region:  "eu-de-1",
			Profile: "bx2-2x8",
		}

		res := instance.BuildResource()
		require.NotNil(t, res)
		require.NotEmpty(t, res.CostComponents)

		instanceHoursComponent := res.CostComponents[0]
		assert.Equal(t, "Instance Hours (bx2-2x8)", instanceHoursComponent.Name)
		assert.Equal(t, "ibm", *instanceHoursComponent.ProductFilter.VendorName)
		assert.Equal(t, "eu-de", *instanceHoursComponent.ProductFilter.Region)
		assert.Equal(t, "is.reservation", *instanceHoursComponent.ProductFilter.Service)
		assert.Equal(t, "bx2-2x8", *instanceHoursComponent.ProductFilter.AttributeFilters[0].Value)
		assert.Equal(t, "RESERVATION_HOURS_HOURLY", *instanceHoursComponent.PriceFilter.Unit)
	})

	t.Run("confidential profile", func(t *testing.T) {
		instance := &ibm.IsInstance{
			Address: "ibm_is_instance.test_confidential",
			Region:  "eu-de-1",
			Profile: "cx3dc-128x320",
		}

		res := instance.BuildResource()
		require.NotNil(t, res)
		require.NotEmpty(t, res.CostComponents)

		instanceHoursComponent := res.CostComponents[0]
		assert.Equal(t, "Instance Hours (cx3dc-128x320)", instanceHoursComponent.Name)
		assert.Equal(t, "ibm", *instanceHoursComponent.ProductFilter.VendorName)
		assert.Equal(t, "eu-de", *instanceHoursComponent.ProductFilter.Region)
		assert.Equal(t, "is.instance", *instanceHoursComponent.ProductFilter.Service)
		assert.Equal(t, "cx3dc-128x320", *instanceHoursComponent.ProductFilter.AttributeFilters[0].Value)
		assert.Equal(t, "INSTANCE_HOURS_MULTI_TENANT", *instanceHoursComponent.PriceFilter.Unit)
	})
}
