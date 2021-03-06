package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAciAny() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAnyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciAnyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("any")
	VRFDn := d.Get("vrf_dn").(string)

	dn := fmt.Sprintf("%s/%s", VRFDn, rn)

	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		return err
	}
	setAnyAttributes(vzAny, d)
	return nil
}
