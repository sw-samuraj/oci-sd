#!/bin/sh

envsubst < oci-sd.tmpl > oci-sd.toml
rm oci-sd.tmpl start.sh
exec bin/oci-sd -c oci-sd.toml -o "${OCI_DISCOVERY_OUTPUT}"

