resource "oci_core_vcn" "oci-sd-vcn" {
  cidr_block = "${var.vcn-cidr}"
  compartment_id = "${var.compartment_ocid}"
  display_name = "oci-sd-vcn"
  dns_label = "ocisdvcn"
}

resource "oci_core_internet_gateway" "oci-sd-gw" {
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_vcn.oci-sd-vcn.id}"
  display_name = "oci-sd-gw"
}

resource "oci_core_route_table" "oci-sd-rt" {
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_vcn.oci-sd-vcn.id}"
  display_name = "oci-sd-rt"

  route_rules {
    destination = "0.0.0.0/0"
    network_entity_id = "${oci_core_internet_gateway.oci-sd-gw.id}"
  }
}

resource "oci_core_security_list" "oci-sd-sl" {
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_vcn.oci-sd-vcn.id}"
  display_name = "oci-sd-sl"
  egress_security_rules = [
    {
      protocol = "6"
      destination = "0.0.0.0/0"
    }]
  ingress_security_rules = [
    {
      tcp_options {
        "max" = 22
        "min" = 22
      }

      protocol = "6"
      source = "0.0.0.0/0"
    },{
      tcp_options {
        "max" = 9100
        "min" = 9100
      }

      protocol = "6"
      source = "0.0.0.0/0"
    }
  ]
}

resource "oci_core_subnet" "oci-sd-subnet" {
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[2],"name")}"
  cidr_block = "${var.subnet-cidr}"
  display_name = "oci-sd-subnet-ad3"
  dns_label = "ocisdsubad3"
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_vcn.oci-sd-vcn.id}"
  route_table_id = "${oci_core_route_table.oci-sd-rt.id}"
  security_list_ids = [
    "${oci_core_security_list.oci-sd-sl.id}"
  ]
  dhcp_options_id = "${oci_core_vcn.oci-sd-vcn.default_dhcp_options_id}"
}
