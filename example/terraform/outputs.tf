output "private-ips" {
  value = ["${oci_core_instance.oci-sd-instance.*.private_ip}"]
}

output "public-ips" {
  value = ["${oci_core_instance.oci-sd-instance.*.public_ip}"]
}
