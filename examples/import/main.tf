provider "aci" {
  username = "admin"
  password = "cisco123"
  url      = "https://192.168.10.102"
  insecure = true
}

resource "aci_tenant" "test" {
  #name = "1234"
}

# resource "aci_app_profile" "testap" {}

# resource "aci_bridge_domain" "testbd" {}

# resource "aci_contract" "testcontract" {}

# resource "aci_subnet" "testsubnet" {
#   name             = "10.0.3.30/30"
#   bridge_domain_dn = "uni/tn-Girish/BD-testingbd"
#   scope            = ["private"]
#   ip_address       = "10.0.3.30/30"
# }

# resource "aci_subject" "testsubject" {}


# resource "aci_filter" "testfilter" {}


# resource "aci_entry" "testentry" {}


# resource "aci_epg" "testepg" {}

