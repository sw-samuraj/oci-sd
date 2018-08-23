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

package oci

import (
	"context"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"time"
)

const (
	ociLabel                   = model.MetaLabelPrefix + "oci_"
	ociLabelAvailabilityDomain = ociLabel + "availability_domain"
	ociLabelCompartmentID      = ociLabel + "compartment_id"
	ociLabelInstanceID         = ociLabel + "instance_id"
	ociLabelInstanceName       = ociLabel + "instance_name"
	ociLabelInstanceState      = ociLabel + "instance_state"
	ociLabelPrivateIP          = ociLabel + "private_ip"
	ociLabelPublicIP           = ociLabel + "public_ip"
	ociLabelDefinedTag         = ociLabel + "defined_tag_"
	ociLabelFreeformTag        = ociLabel + "freeform_tag_"
)

// defaultSDConfig is the default OCI SD configuration
var defaultSDConfig = SDConfig{
	Port:            9100,
	RefreshInterval: model.Duration(60 * time.Second),
}

// SDConfig is the configuration for OCI based service discovery.
type SDConfig struct {
	User            string
	FingerPrint     string
	KeyFile         string
	PassPhrase      string `toml:",omitempty"`
	Tenancy         string
	Region          string
	Compartment     string
	Port            int            `toml:",omitempty"`
	RefreshInterval model.Duration `toml:",omitempty"`
}

func (c *SDConfig) Validate() error {
	if c.User == "" {
		return fmt.Errorf("oci sd configuration requires a user")
	}
	if c.FingerPrint == "" {
		return fmt.Errorf("oci sd configuration requires a fingerprint")
	}
	if c.KeyFile == "" {
		return fmt.Errorf("oci sd configuration requires a key file")
	}
	if c.Tenancy == "" {
		return fmt.Errorf("oci sd configuration requires a tenancy")
	}
	if c.Region == "" {
		return fmt.Errorf("oci sd configuration requires a region")
	}
	if c.Compartment == "" {
		return fmt.Errorf("oci sd configuration requires a compartment")
	}

	return nil
}

func (c *SDConfig) ApplyDefault() error {
	if err := mergo.Merge(c, defaultSDConfig); err != nil {
		return err
	}
	return nil
}

// discovery periodically performs OCI-SD requests. It implements
// the Discoverer interface.
type discovery struct {
	sdConfig  *SDConfig
	ociConfig common.ConfigurationProvider
	interval  time.Duration
	logger    log.Logger
}

// NewDiscovery returns a new OCI discovery which periodically refreshes its targets.
func NewDiscovery(conf *SDConfig, logger *log.Logger) (*discovery, error) {
	if logger == nil {
		logger = log.New()
	}
	privateKey, err := loadKey(conf.KeyFile)
	if err != nil {
		return nil, err
	}
	ociConfig := common.NewRawConfigurationProvider(
		conf.Tenancy,
		conf.User,
		conf.Region,
		conf.FingerPrint,
		privateKey,
		&conf.PassPhrase,
	)
	return &discovery{
		sdConfig:  conf,
		ociConfig: ociConfig,
		interval:  time.Duration(conf.RefreshInterval),
		logger:    *logger,
	}, nil
}

// Run implements the Discoverer interface.
func (d *discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	// Get an initial set right away.
	tg, err := d.refresh()
	if err != nil {
		d.logger.WithField("err", err).Error("refreshing targets failed")
	} else {
		select {
		case ch <- []*targetgroup.Group{tg}:
		case <-ctx.Done():
			return
		}
	}

	for {
		select {
		case <-ticker.C:
			tg, err := d.refresh()
			if err != nil {
				d.logger.WithField("err", err).Error("refreshing targets failed")
				continue
			}

			select {
			case ch <- []*targetgroup.Group{tg}:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (d *discovery) refresh() (tg *targetgroup.Group, err error) {
	tg = &targetgroup.Group{
		Source: d.sdConfig.Region,
	}

	computeClient, err := core.NewComputeClientWithConfigurationProvider(d.ociConfig)
	if err != nil {
		return nil, err
	}
	vnicClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(d.ociConfig)
	if err != nil {
		return nil, err
	}
	res, err := computeClient.ListInstances(
		context.Background(),
		core.ListInstancesRequest{CompartmentId: &d.sdConfig.Compartment},
	)
	if err != nil {
		return nil, fmt.Errorf("could not obtain list of instances: %s", err)
	}

	for _, instance := range res.Items {
		res, err := computeClient.ListVnicAttachments(
			context.Background(),
			core.ListVnicAttachmentsRequest{
				CompartmentId: &d.sdConfig.Compartment,
				InstanceId:    instance.Id,
			},
		)
		if err != nil {
			d.logger.WithField("ocid", instance.Id).Error("could not obtain attached vnic")
			continue
		}
		for _, vnic := range res.Items {
			res, err := vnicClient.GetVnic(
				context.Background(),
				core.GetVnicRequest{VnicId: vnic.VnicId},
			)
			if err != nil {
				d.logger.WithField("ocid", vnic.VnicId).Error("could not obtain vnic")
				continue
			}
			if *res.IsPrimary {
				labels := model.LabelSet{
					ociLabelInstanceID:         model.LabelValue(*instance.Id),
					ociLabelInstanceName:       model.LabelValue(*instance.DisplayName),
					ociLabelInstanceState:      model.LabelValue(instance.LifecycleState),
					ociLabelCompartmentID:      model.LabelValue(*instance.CompartmentId),
					ociLabelAvailabilityDomain: model.LabelValue(*instance.AvailabilityDomain),
					ociLabelPrivateIP:          model.LabelValue(*res.PrivateIp),
				}
				if *res.PublicIp != "" {
					labels[ociLabelPublicIP] = model.LabelValue(*res.PublicIp)
				}
				addr := net.JoinHostPort(*res.PrivateIp, fmt.Sprintf("%d", d.sdConfig.Port))
				labels[model.AddressLabel] = model.LabelValue(addr)
				for key, value := range instance.FreeformTags {
					labels[ociLabelFreeformTag+model.LabelName(key)] = model.LabelValue(value)
				}
				for ns, tags := range instance.DefinedTags {
					for key, value := range tags {
						labelName := model.LabelName(ociLabelDefinedTag + ns + "_" + key)
						labels[labelName] = model.LabelValue(value.(string))
					}
				}
				tg.Targets = append(tg.Targets, labels)
			}
		}
	}

	return tg, nil
}

func loadKey(keyFile string) (string, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
