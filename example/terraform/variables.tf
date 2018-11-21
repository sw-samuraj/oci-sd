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
  default = "VM.Standard2.1"
}

variable "instance_image_ocid" {
  description = "Oracle-provided image 'Oracle-Linux-7.4-2018.02.21-1'"
  type = "map"
  default = {
    // See https://docs.us-phoenix-1.oraclecloud.com/images/
    us-phoenix-1 = "ocid1.image.oc1.phx.aaaaaaaaupbfz5f5hdvejulmalhyb6goieolullgkpumorbvxlwkaowglslq"
    us-ashburn-1   = "ocid1.image.oc1.iad.aaaaaaaajlw3xfie2t5t52uegyhiq2npx7bqyu4uvi2zyu3w3mqayc2bxmaa"
    eu-frankfurt-1 = "ocid1.image.oc1.eu-frankfurt-1.aaaaaaaa7d3fsb6272srnftyi4dphdgfjf6gurxqhmv6ileds7ba3m2gltxq"
    uk-london-1    = "ocid1.image.oc1.uk-london-1.aaaaaaaaa6h6gj6v4n56mqrbgnosskq63blyv2752g36zerymy63cfkojiiq"
  }
}

variable "BootStrap" {
  default = "./userdata/bootstrap"
}
