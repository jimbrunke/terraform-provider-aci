package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciVRF() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVRFCreate,
		Update: resourceAciVRFUpdate,
		Read:   resourceAciVRFRead,
		Delete: resourceAciVRFDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVRFImport,
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

			"bd_enforced_enable": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"ip_data_plane_learning": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"knw_mcast_act": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "specifies if known multicast traffic is forwarded",
			},

			"name_alias": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"pc_enf_dir": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"pc_enf_pref": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "preferred policy control",
			},

			"relation_fv_rs_ospf_ctx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to ospfCtxPol",
			},
			"relation_fv_rs_vrf_validation_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to l3extVrfValidationPol",
			},
			"relation_fv_rs_ctx_mcast_to": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzFilter",
				Set:         schema.HashString,
			},
			"relation_fv_rs_ctx_to_eigrp_ctx_af_pol": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to eigrpCtxAfPol",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_eigrp_ctx_af_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_ctx_to_ospf_ctx_pol": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to ospfCtxPol",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_ospf_ctx_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_ctx_to_ep_ret": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvEpRetPol",
			},
			"relation_fv_rs_bgp_ctx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to bgpCtxPol",
			},
			"relation_fv_rs_ctx_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to monEPGPol",
			},
			"relation_fv_rs_ctx_to_ext_route_tag_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to l3extRouteTagPol",
			},
			"relation_fv_rs_ctx_to_bgp_ctx_af_pol": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to bgpCtxAfPol",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_bgp_ctx_af_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		}),
	}
}
func getRemoteVRF(client *client.Client, dn string) (*models.VRF, error) {
	fvCtxCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvCtx := models.VRFFromContainer(fvCtxCont)

	if fvCtx.DistinguishedName == "" {
		return nil, fmt.Errorf("VRF %s not found", fvCtx.DistinguishedName)
	}

	return fvCtx, nil
}

func setVRFAttributes(fvCtx *models.VRF, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvCtx.DistinguishedName)
	d.Set("description", fvCtx.Description)
	d.Set("name", GetMOName(fvCtx.DistinguishedName))
	d.Set("tenant_dn", GetParentDn(fvCtx.DistinguishedName))
	fvCtxMap, _ := fvCtx.ToMap()

	d.Set("annotation", fvCtxMap["annotation"])
	d.Set("bd_enforced_enable", fvCtxMap["bdEnforcedEnable"])
	d.Set("ip_data_plane_learning", fvCtxMap["ipDataPlaneLearning"])
	d.Set("knw_mcast_act", fvCtxMap["knwMcastAct"])
	d.Set("name_alias", fvCtxMap["nameAlias"])
	d.Set("pc_enf_dir", fvCtxMap["pcEnfDir"])
	d.Set("pc_enf_pref", fvCtxMap["pcEnfPref"])
	return d
}

func resourceAciVRFImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	fvCtx, err := getRemoteVRF(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVRFAttributes(fvCtx, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVRFCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvCtxAttr := models.VRFAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvCtxAttr.Annotation = Annotation.(string)
	}
	if BdEnforcedEnable, ok := d.GetOk("bd_enforced_enable"); ok {
		fvCtxAttr.BdEnforcedEnable = BdEnforcedEnable.(string)
	}
	if IpDataPlaneLearning, ok := d.GetOk("ip_data_plane_learning"); ok {
		fvCtxAttr.IpDataPlaneLearning = IpDataPlaneLearning.(string)
	}
	if KnwMcastAct, ok := d.GetOk("knw_mcast_act"); ok {
		fvCtxAttr.KnwMcastAct = KnwMcastAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvCtxAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfDir, ok := d.GetOk("pc_enf_dir"); ok {
		fvCtxAttr.PcEnfDir = PcEnfDir.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvCtxAttr.PcEnfPref = PcEnfPref.(string)
	}
	fvCtx := models.NewVRF(fmt.Sprintf("ctx-%s", name), TenantDn, desc, fvCtxAttr)

	err := aciClient.Save(fvCtx)
	if err != nil {
		return err
	}

	if relationTofvRsOspfCtxPol, ok := d.GetOk("relation_fv_rs_ospf_ctx_pol"); ok {
		relationParam := relationTofvRsOspfCtxPol.(string)
		err = aciClient.CreateRelationfvRsOspfCtxPolFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsVrfValidationPol, ok := d.GetOk("relation_fv_rs_vrf_validation_pol"); ok {
		relationParam := relationTofvRsVrfValidationPol.(string)
		err = aciClient.CreateRelationfvRsVrfValidationPolFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsCtxMcastTo, ok := d.GetOk("relation_fv_rs_ctx_mcast_to"); ok {
		relationParamList := toStringList(relationTofvRsCtxMcastTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsCtxMcastToFromVRF(fvCtx.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTofvRsCtxToEigrpCtxAfPol, ok := d.GetOk("relation_fv_rs_ctx_to_eigrp_ctx_af_pol"); ok {

		relationParamList := relationTofvRsCtxToEigrpCtxAfPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}
	if relationTofvRsCtxToOspfCtxPol, ok := d.GetOk("relation_fv_rs_ctx_to_ospf_ctx_pol"); ok {

		relationParamList := relationTofvRsCtxToOspfCtxPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}
	if relationTofvRsCtxToEpRet, ok := d.GetOk("relation_fv_rs_ctx_to_ep_ret"); ok {
		relationParam := relationTofvRsCtxToEpRet.(string)
		err = aciClient.CreateRelationfvRsCtxToEpRetFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsBgpCtxPol, ok := d.GetOk("relation_fv_rs_bgp_ctx_pol"); ok {
		relationParam := relationTofvRsBgpCtxPol.(string)
		err = aciClient.CreateRelationfvRsBgpCtxPolFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsCtxMonPol, ok := d.GetOk("relation_fv_rs_ctx_mon_pol"); ok {
		relationParam := relationTofvRsCtxMonPol.(string)
		err = aciClient.CreateRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsCtxToExtRouteTagPol, ok := d.GetOk("relation_fv_rs_ctx_to_ext_route_tag_pol"); ok {
		relationParam := relationTofvRsCtxToExtRouteTagPol.(string)
		err = aciClient.CreateRelationfvRsCtxToExtRouteTagPolFromVRF(fvCtx.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsCtxToBgpCtxAfPol, ok := d.GetOk("relation_fv_rs_ctx_to_bgp_ctx_af_pol"); ok {

		relationParamList := relationTofvRsCtxToBgpCtxAfPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}

	d.SetId(fvCtx.DistinguishedName)
	return resourceAciVRFRead(d, m)
}

func resourceAciVRFUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvCtxAttr := models.VRFAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvCtxAttr.Annotation = Annotation.(string)
	}
	if BdEnforcedEnable, ok := d.GetOk("bd_enforced_enable"); ok {
		fvCtxAttr.BdEnforcedEnable = BdEnforcedEnable.(string)
	}
	if IpDataPlaneLearning, ok := d.GetOk("ip_data_plane_learning"); ok {
		fvCtxAttr.IpDataPlaneLearning = IpDataPlaneLearning.(string)
	}
	if KnwMcastAct, ok := d.GetOk("knw_mcast_act"); ok {
		fvCtxAttr.KnwMcastAct = KnwMcastAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvCtxAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfDir, ok := d.GetOk("pc_enf_dir"); ok {
		fvCtxAttr.PcEnfDir = PcEnfDir.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvCtxAttr.PcEnfPref = PcEnfPref.(string)
	}
	fvCtx := models.NewVRF(fmt.Sprintf("ctx-%s", name), TenantDn, desc, fvCtxAttr)

	fvCtx.Status = "modified"

	err := aciClient.Save(fvCtx)

	if err != nil {
		return err
	}

	if d.HasChange("relation_fv_rs_ospf_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ospf_ctx_pol")
		err = aciClient.CreateRelationfvRsOspfCtxPolFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_vrf_validation_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_vrf_validation_pol")
		err = aciClient.CreateRelationfvRsVrfValidationPolFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_ctx_mcast_to") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_mcast_to")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsCtxMcastToFromVRF(fvCtx.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_eigrp_ctx_af_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_eigrp_ctx_af_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ospf_ctx_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_ospf_ctx_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ep_ret")
		err = aciClient.CreateRelationfvRsCtxToEpRetFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_bgp_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_bgp_ctx_pol")
		err = aciClient.CreateRelationfvRsBgpCtxPolFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_ctx_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_mon_pol")
		err = aciClient.DeleteRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ext_route_tag_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ext_route_tag_pol")
		err = aciClient.CreateRelationfvRsCtxToExtRouteTagPolFromVRF(fvCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_bgp_ctx_af_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_bgp_ctx_af_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
		}

	}

	d.SetId(fvCtx.DistinguishedName)
	return resourceAciVRFRead(d, m)

}

func resourceAciVRFRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	fvCtx, err := getRemoteVRF(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVRFAttributes(fvCtx, d)
	return nil
}

func resourceAciVRFDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvCtx")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
