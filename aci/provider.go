package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_USERNAME", nil),
				Description: "Username for the APIC Account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PASSWORD", nil),
				Description: "Password for the APIC Account",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_URL", nil),
				Description: "URL of the Cisco ACI web interface",
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure HTTPS client",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PRIVATE_KEY", nil),
				Description: "Private key path for signature calculation",
			},
			"cert_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_CERT_NAME", nil),
				Description: "Certificate name for the User in Cisco ACI.",
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PROXY_URL", nil),
				Description: "Proxy Server URL with port number",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aci_tenant":                            resourceAciTenant(),
			"aci_application_profile":               resourceAciApplicationProfile(),
			"aci_bridge_domain":                     resourceAciBridgeDomain(),
			"aci_contract":                          resourceAciContract(),
			"aci_application_epg":                   resourceAciApplicationEPG(),
			"aci_contract_subject":                  resourceAciContractSubject(),
			"aci_subnet":                            resourceAciSubnet(),
			"aci_filter":                            resourceAciFilter(),
			"aci_filter_entry":                      resourceAciFilterEntry(),
			"aci_vmm_domain":                        resourceAciVMMDomain(),
			"aci_vrf":                               resourceAciVRF(),
			"aci_rest":                              resourceAciRest(),
			"aci_external_network_instance_profile": resourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                        resourceAciL3Outside(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aci_tenant":                            dataSourceAciTenant(),
			"aci_application_profile":               dataSourceAciApplicationProfile(),
			"aci_bridge_domain":                     dataSourceAciBridgeDomain(),
			"aci_contract":                          dataSourceAciContract(),
			"aci_application_epg":                   dataSourceAciApplicationEPG(),
			"aci_contract_subject":                  dataSourceAciContractSubject(),
			"aci_subnet":                            dataSourceAciSubnet(),
			"aci_filter":                            dataSourceAciFilter(),
			"aci_filter_entry":                      dataSourceAciFilterEntry(),
			"aci_vmm_domain":                        dataSourceAciVMMDomain(),
			"aci_vrf":                               dataSourceAciVRF(),
			"aci_external_network_instance_profile": dataSourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                        dataSourceAciL3Outside(),
		},

		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		URL:        d.Get("url").(string),
		IsInsecure: d.Get("insecure").(bool),
		PrivateKey: d.Get("private_key").(string),
		Certname:   d.Get("cert_name").(string),
		ProxyUrl:   d.Get("proxy_url").(string),
	}

	if err := config.Valid(); err != nil {
		return nil, err
	}

	return config.getClient(), nil
}

func (c Config) Valid() error {

	if c.Username == "" {
		return fmt.Errorf("Username must be provided for the ACI provider")
	}

	if c.Password == "" {
		if c.PrivateKey == "" && c.Certname == "" {

			return fmt.Errorf("Either of private_key/cert_name or password is required")
		} else if c.PrivateKey == "" || c.Certname == "" {
			return fmt.Errorf("private_key and cert_name both must be provided")
		}
	}

	if c.URL == "" {
		return fmt.Errorf("The URL must be provided for the ACI provider")
	}

	return nil
}

func (c Config) getClient() interface{} {
	if c.Password != "" {

		return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl))

	} else {

		return client.GetClient(c.URL, c.Username, client.PrivateKey(c.PrivateKey), client.AdminCert(c.Certname), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl))
	}
}

// Config
type Config struct {
	Username   string
	Password   string
	URL        string
	IsInsecure bool
	PrivateKey string
	Certname   string
	ProxyUrl   string
}
