package ibm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

var imageMap map[string]struct {
	Vendor  string
	Version string
}

type Image struct {
	Href            string          `json:"href"`
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	OperatingSystem OperatingSystem `json:"operating_system"`
}

type OperatingSystem struct {
	AllowUserImageCreation bool   `json:"allow_user_image_creation"`
	Architecture           string `json:"architecture"`
	DedicatedHostOnly      bool   `json:"dedicated_host_only"`
	DisplayName            string `json:"display_name"`
	Family                 string `json:"family"`
	Href                   string `json:"href"`
	Name                   string `json:"name"`
	UserDataFormat         string `json:"user_data_format"`
	Vendor                 string `json:"vendor"`
	Version                string `json:"version"`
}

func getIsInstanceRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_is_instance",
		RFunc: newIsInstance,
	}
}

// valid profile values https://cloud.ibm.com/docs/vpc?topic=vpc-profiles&interface=ui
// profile names in Global Catalog contain dots instead of dashes
func newIsInstance(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	imagePath := "./images.json" // path for using infracost excecutable

	byteValue, err := os.ReadFile(imagePath)
	if err != nil {
		byteValue, err = os.ReadFile("../../../../images.json") //path for individual testing directly from this file (go test)
		if err != nil {
			fmt.Printf("Error reading file: %v", err)
		}
	}

	var images []Image
	err = json.Unmarshal(byteValue, &images)
	if err != nil {
		fmt.Printf("Error unmarshaling json: %v ", err)
	}

	imageMap = make(map[string]struct {
		Vendor  string
		Version string
	})

	for _, image := range images {
		imageMap[image.ID] = struct {
			Vendor  string
			Version string
		}{
			Vendor:  image.OperatingSystem.Vendor,
			Version: image.OperatingSystem.Version,
		}
	}

	imageId := d.Get("image").String()
	region := d.Get("region").String()
	profile := d.Get("profile").String()
	vendor := imageMap[imageId].Vendor
	version := imageMap[imageId].Version
	zone := d.Get("zone").String()
	dedicatedHost := strings.TrimSpace(d.Get("dedicated_host").String())
	dedicatedHostGroup := strings.TrimSpace(d.Get("dedicated_host_group").String())
	isDedicated := !((dedicatedHost == "") && (dedicatedHostGroup == ""))
	name := d.Get("name").String()

	// Defaults
	bootVolumeName := "Unnamed boot volume"
	var bootVolumeSize int64 = 100

	bv := d.Get("boot_volume").Array()
	if len(bv) > 0 {
		if bv[0].Get("name").String() != "" {
			bootVolumeName = bv[0].Get("name").String()
		}
		if bv[0].Get("size").Int() != 0 {
			bootVolumeSize = bv[0].Get("size").Int()
		}
	}

	r := &ibm.IsInstance{
		Address:     d.Address,
		Region:      region,
		Profile:     profile,
		Vendor:      vendor,
		Version:     version,
		Zone:        zone,
		IsDedicated: isDedicated,
		BootVolume: struct {
			Name string
			Size int64
		}{Name: bootVolumeName, Size: bootVolumeSize},
	}

	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["name"] = name
	configuration["on_dedicated_host"] = isDedicated
	configuration["profile"] = profile
	configuration["region"] = region

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}
