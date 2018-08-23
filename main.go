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
	flag "github.com/spf13/pflag"
	"github.com/go-kit/kit/log"
	"github.com/sw-samuraj/oci-sd/oci"
	"github.com/sirupsen/logrus"
	"github.com/prometheus/documentation/examples/custom-sd/adapter"
	"context"
	"os"
)

var (
	configFile string
	outputFile string
)

func init() {
	flag.StringVarP(&configFile, "config-file", "c", "oci-sd.toml", "external config file")
	flag.StringVarP(&outputFile, "output-file", "o", "oci-sd.json", "output file for file_sd compatible file")
}

func main() {
	flag.Parse()

	logger := logrus.New()
	cfg := parseConfig()
	ctx := context.Background()

	disc, err := oci.NewDiscovery(&cfg.SDConfig, logger)
	if err != nil {
		logger.WithFields(logrus.Fields{"func": "main", "err": err}).Fatal("can't create a discovery")
	}

	sdAdapter := adapter.NewAdapter(ctx, outputFile, "ociSD", disc, log.NewLogfmtLogger(os.Stdout))
	sdAdapter.Run()

	<-ctx.Done()
}