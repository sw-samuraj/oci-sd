package main

import "testing"

func TestReadConfigFile(t *testing.T) {
	configFile = "testdata/oci-sd.toml"
	cf, _ := readConfigFile()

	if len(cf) == 0 {
		t.Error("length of config file content == 0, expected > 0")
	}
}
