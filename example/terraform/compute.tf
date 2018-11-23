resource "oci_core_instance" "oci-sd-instance" {
  count = "${var.NumInstances}"
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[2],"name")}"
  compartment_id = "${var.compartment_ocid}"
  display_name = "oci-sd-instance-${count.index}"
  shape = "${var.instance_shape}"
  subnet_id = "${oci_core_subnet.oci-sd-subnet.id}"

  source_details {
    source_type = "image"
    source_id = "${var.instance_image_ocid[var.region]}"
  }

  metadata {
    # ssh_authorized_keys = "${file(var.ssh_public_key)}"
    user_data = "${base64encode(file(var.BootStrap))}"
  }

  freeform_tags = "${
    map(
      "prometheus_exporter", "node_exporter"
    )
  }"
}