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
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/sw-samuraj/oci-sd/adapter"
	"github.com/sw-samuraj/oci-sd/oci"
)

var (
	configFile  string
	outputFile  string
	authvar     bool
	cfg         config
	compartment string
	sanitise    bool
)

func init() {
	flag.StringVarP(&configFile, "config-file", "c", "oci-sd.toml", "external config file")
	flag.StringVarP(&outputFile, "output-file", "o", "oci-sd.json", "output file for file_sd compatible file")
	flag.BoolVarP(&authvar, "instance-principal", "i", false, "initialise with instance principal authentication")
	flag.StringVarP(&compartment, "compartment", "t", "", "compartment for discovering targets")
	flag.BoolVarP(&sanitise, "sanitise", "s", false, "sanitise instance tags to fit Prometheus requirements by removing special characters (:, -)")
}

const LOG_PATH = "/var/log/oci-sd/"

func main() {
	flag.Parse()
	logger := log.New()
	logSetup(LOG_PATH, logger)

	if authvar {
		logger.Info("initialising with instance principal authentication")
		if compartment == "" {
			logger.Fatal("flag --compartment (or shorthand -t) cannot be empty if instance principal is used")
		}
		cfg.SDConfig.InstancePrincipal = true
		cfg.SDConfig.Compartment = compartment
		if err := cfg.SDConfig.ApplyDefault(); err != nil {
			logger.WithField("error", err).Warn("error to apply default config values")
		}
	} else {
		logger.Info("initialising with user authentication")
		cfg = parseConfig()
	}
	cfg.SDConfig.Sanitise = sanitise
	ctx := context.Background()

	disc, err := oci.NewDiscovery(&cfg.SDConfig, logger)
	if err != nil {
		logger.WithFields(log.Fields{"func": "main", "err": err}).Fatal("can't create a discovery")
	}

	sdAdapter := adapter.NewAdapter(ctx, outputFile, "ociSD", disc, *logger)
	sdAdapter.Run()

	<-ctx.Done()
}

func logSetup(logPath string, logger *log.Logger) {
	logger.SetReportCaller(true)
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		logger.Fatal("Can't create a log path. Err: ", err)
	}
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logger.SetFormatter(Formatter)
	logFile, err := os.OpenFile(logPath+"oci-sd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal("Can't create a log path. Err: ", err)
	}
	logger.SetOutput(logFile)
}
