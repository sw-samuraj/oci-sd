resource "oci_core_instance" "oci-sd-instance" {
  count = "${var.NumInstances}"
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[2],"name")}"
  compartment_id = "${var.compartment_ocid}"
  display_name = "oci-sd-instance-${count.index}"
  shape = "${var.instance_shape}"

  shape_config {
      ocpus = "${var.instance_ocpus}"
      memory_in_gbs = "${var.instance_shape_config_memory_in_gbs}"
    }

  source_details {
    source_type = "image"
    source_id = "${var.instance_image_ocid[var.region]}"
  }

  create_vnic_details {
      subnet_id                 = oci_core_subnet.oci-sd-subnet.id
      display_name              = "Primaryvnic"
      assign_public_ip          = true
      assign_private_dns_record = true
    }

  metadata = {
    # ssh_authorized_keys = "${file(var.ssh_public_key)}"
    user_data = "${base64encode(file(var.BootStrap))}"
  }

  freeform_tags = {
    "prometheus_exporter" = "node_exporter"
  }
}