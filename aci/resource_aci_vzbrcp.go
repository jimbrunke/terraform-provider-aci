package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractCreate,
		Update: resourceAciContractUpdate,
		Read:   resourceAciContractRead,
		Delete: resourceAciContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"name_alias": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"prio": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "priority level of the service contract",
			},

			"scope": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "scope of contract",
			},

			"target_dscp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "target dscp",
			},

			"relation_vz_rs_graph_att": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to vnsAbsGraph",
			},
		}),
	}
}
func getRemoteContract(client *client.Client, dn string) (*models.Contract, error) {
	vzBrCPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzBrCP := models.ContractFromContainer(vzBrCPCont)

	if vzBrCP.DistinguishedName == "" {
		return nil, fmt.Errorf("Contract %s not found", vzBrCP.DistinguishedName)
	}

	return vzBrCP, nil
}

func setContractAttributes(vzBrCP *models.Contract, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vzBrCP.DistinguishedName)
	d.Set("description", vzBrCP.Description)
	d.Set("name", GetMOName(vzBrCP.DistinguishedName))
	d.Set("tenant_dn", GetParentDn(vzBrCP.DistinguishedName))
	vzBrCPMap, _ := vzBrCP.ToMap()

	d.Set("annotation", vzBrCPMap["annotation"])
	d.Set("name_alias", vzBrCPMap["nameAlias"])
	d.Set("prio", vzBrCPMap["prio"])
	d.Set("scope", vzBrCPMap["scope"])
	d.Set("target_dscp", vzBrCPMap["targetDscp"])
	return d
}

func resourceAciContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setContractAttributes(vzBrCP, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	err := aciClient.Save(vzBrCP)
	if err != nil {
		return err
	}

	if relationTovzRsGraphAtt, ok := d.GetOk("relation_vz_rs_graph_att"); ok {
		relationParam := relationTovzRsGraphAtt.(string)
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	d.SetId(vzBrCP.DistinguishedName)
	return resourceAciContractRead(d, m)
}

func resourceAciContractUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	vzBrCP.Status = "modified"

	err := aciClient.Save(vzBrCP)

	if err != nil {
		return err
	}

	if d.HasChange("relation_vz_rs_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_graph_att")
		err = aciClient.DeleteRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}

	d.SetId(vzBrCP.DistinguishedName)
	return resourceAciContractRead(d, m)

}

func resourceAciContractRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setContractAttributes(vzBrCP, d)
	return nil
}

func resourceAciContractDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzBrCP")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
