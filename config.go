// Copyright (c) 2018, Vít Kotačka
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"errors"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/sw-samuraj/oci-sd/oci"
	"io/ioutil"
)

type config struct {
	SDConfig oci.SDConfig
}

func parseConfig() config {
	logger := log.WithField("func", "parseConfig")

	file, err := readConfigFile()
	if err != nil {
		logger.WithFields(log.Fields{"file": configFile, "error": err}).Fatal("can't read a config file")
	}
	config := config{}
	toml.Unmarshal(file, &config)

	if err := config.SDConfig.Validate(); err != nil {
		logger.WithField("error", err).Fatal("invalid config file")
	}

	if err := config.SDConfig.ApplyDefault(); err != nil {
		logger.WithField("error", err).Warn("error to apply default config values")
	}

	return config
}

func readConfigFile() ([]byte, error) {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("config file has not been found")
	}

	return configFile, nil
}
