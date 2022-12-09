variable "tenancy_ocid" {
  description = "The global identifier for your account, always shown on the bottom of the web console."
}

variable "user_ocid" {
  description = " The identifier of the user account you will be using for Terraform."
}

variable "private_key_path" {
  description = "The path to the private key stored on your computer."
}

variable "fingerprint" {
  description = "The fingerprint of the public key."
}

variable "pass_phrase" {
  description = "The pass phrase to the private key."
  default = ""
}

variable "region" {
  description = "The region to target with this provider configuration."
  default = "us-phoenix-1"
}

variable "compartment_ocid" {
  description = "The identifier of the compartment you will be using for Terraform."
}

variable "NumInstances" {
  description = "Defines the number of instances to deploy."
  default = "3"
}

variable "vcn-cidr" {
  default = "10.1.0.0/24"
}

variable "subnet-cidr" {
  default = "10.1.0.0/28"
}

variable "instance_shape" {
  default = "VM.Standard.E4.Flex"
}

variable "instance_ocpus" {
  default = 1
}

variable "instance_shape_config_memory_in_gbs" {
  default = 8
}

variable "instance_image_ocid" {
  description = "Oracle-provided image 'Oracle-Linux-8.6-2022.10.04-0'"
  type = map(string)
  default = {
    // See https://docs.us-phoenix-1.oraclecloud.com/images/
    us-phoenix-1 = "ocid1.image.oc1.phx.aaaaaaaaqdlspgo5d5tdew5m3ntswbkxfoclc35nvcv3r3a7a5wjwxphuoeq"
    us-ashburn-1   = "ocid1.image.oc1.iad.aaaaaaaaorro6lk6mljfs3dafptdskbupyjjbindwgqc6nf4ohbe3ucklrqq"
    eu-frankfurt-1 = "ocid1.image.oc1.eu-frankfurt-1.aaaaaaaa47555lp4mjbiuf64doxtnbimrwk57m4sfgu3gonaf5i2cteil5iq"
    uk-london-1    = "ocid1.image.oc1.uk-london-1.aaaaaaaa5e2iiw5k4gclwn2akqjxl7xwtohw5ivx3ly7q3cjn6ibggx5ywla"
  }
}

variable "BootStrap" {
  default = "./userdata/bootstrap"
}

/*
variable "ssh_public_key" {
  description = "The path to the SSL public key stored on your computer."
}
*/


