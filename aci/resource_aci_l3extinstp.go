package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciExternalNetworkInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciExternalNetworkInstanceProfileCreate,
		Update: resourceAciExternalNetworkInstanceProfileUpdate,
		Read:   resourceAciExternalNetworkInstanceProfileRead,
		Delete: resourceAciExternalNetworkInstanceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciExternalNetworkInstanceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
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

			"exception_tag": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"flood_on_encap": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"match_t": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "match criteria",
			},

			"name_alias": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"pref_gr_memb": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"prio": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "qos priority class id",
			},

			"target_dscp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "target dscp",
			},

			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to fvEPg",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prov": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Set:         schema.HashString,
			},
			"relation_l3ext_rs_l3_inst_p_to_dom_p": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to extnwDomP",
			},
			"relation_l3ext_rs_inst_p_to_nat_mapping_e_pg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvAEPg",
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzCPIf",
				Set:         schema.HashString,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to qosCustomPol",
			},
			"relation_l3ext_rs_inst_p_to_profile": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to rtctrlProfile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzTaboo",
				Set:         schema.HashString,
			},
			"relation_fv_rs_intra_epg": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Set:         schema.HashString,
			},
		}),
	}
}
func getRemoteExternalNetworkInstanceProfile(client *client.Client, dn string) (*models.ExternalNetworkInstanceProfile, error) {
	l3extInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extInstP := models.ExternalNetworkInstanceProfileFromContainer(l3extInstPCont)

	if l3extInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("Bridge Domain %s not found", l3extInstP.DistinguishedName)
	}

	return l3extInstP, nil
}

func setExternalNetworkInstanceProfileAttributes(l3extInstP *models.ExternalNetworkInstanceProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extInstP.DistinguishedName)
	d.Set("description", l3extInstP.Description)
	d.Set("name", GetMOName(l3extInstP.DistinguishedName))
	d.Set("l3_outside_dn", GetParentDn(l3extInstP.DistinguishedName))
	l3extInstPMap, _ := l3extInstP.ToMap()

	d.Set("annotation", l3extInstPMap["annotation"])
	d.Set("exception_tag", l3extInstPMap["exceptionTag"])
	d.Set("flood_on_encap", l3extInstPMap["floodOnEncap"])
	d.Set("match_t", l3extInstPMap["matchT"])
	d.Set("name_alias", l3extInstPMap["nameAlias"])
	d.Set("pref_gr_memb", l3extInstPMap["prefGrMemb"])
	d.Set("prio", l3extInstPMap["prio"])
	d.Set("target_dscp", l3extInstPMap["targetDscp"])
	return d
}

func resourceAciExternalNetworkInstanceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setExternalNetworkInstanceProfileAttributes(l3extInstP, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciExternalNetworkInstanceProfileCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		l3extInstPAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		l3extInstPAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		l3extInstPAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extInstPAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		l3extInstPAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extInstPAttr.Prio = Prio.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extInstPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extInstP := models.NewExternalNetworkInstanceProfile(fmt.Sprintf("instP-%s", name), L3OutsideDn, desc, l3extInstPAttr)

	err := aciClient.Save(l3extInstP)
	if err != nil {
		return err
	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTol3extRsL3InstPToDomP, ok := d.GetOk("relation_l3ext_rs_l3_inst_p_to_dom_p"); ok {
		relationParam := relationTol3extRsL3InstPToDomP.(string)
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTol3extRsInstPToNatMappingEPg, ok := d.GetOk("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg"); ok {
		relationParam := relationTol3extRsInstPToNatMappingEPg.(string)
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTol3extRsInstPToProfile, ok := d.GetOk("relation_l3ext_rs_inst_p_to_profile"); ok {

		relationParamList := relationTol3extRsInstPToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}

	d.SetId(l3extInstP.DistinguishedName)
	return resourceAciExternalNetworkInstanceProfileRead(d, m)
}

func resourceAciExternalNetworkInstanceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		l3extInstPAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		l3extInstPAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		l3extInstPAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extInstPAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		l3extInstPAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extInstPAttr.Prio = Prio.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extInstPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extInstP := models.NewExternalNetworkInstanceProfile(fmt.Sprintf("instP-%s", name), L3OutsideDn, desc, l3extInstPAttr)

	l3extInstP.Status = "modified"

	err := aciClient.Save(l3extInstP)

	if err != nil {
		return err
	}

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_l3ext_rs_l3_inst_p_to_dom_p") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_inst_p_to_dom_p")
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg")
		err = aciClient.DeleteRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_profile") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_inst_p_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
		}

	}
	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}

	d.SetId(l3extInstP.DistinguishedName)
	return resourceAciExternalNetworkInstanceProfileRead(d, m)

}

func resourceAciExternalNetworkInstanceProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setExternalNetworkInstanceProfileAttributes(l3extInstP, d)
	return nil
}

func resourceAciExternalNetworkInstanceProfileDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extInstP")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
