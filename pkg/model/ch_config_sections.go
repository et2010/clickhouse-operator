// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/config"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

type configSections struct {
	// commonConfigSections maps section name to section XML config string
	commonConfigSections map[string]string
	// commonUsersConfigSections maps section name to section XML config string
	commonUsersConfigSections map[string]string

	// ClickHouse config generator
	chConfigGenerator *ClickHouseConfigGenerator
	// clickhouse-operator configuration
	chopConfig *config.Config
}

func NewConfigSections(chConfigGenerator *ClickHouseConfigGenerator, chopConfig *config.Config) *configSections {
	return &configSections{
		commonConfigSections:      make(map[string]string),
		commonUsersConfigSections: make(map[string]string),
		chConfigGenerator:         chConfigGenerator,
		chopConfig:                chopConfig,
	}
}

func (c *configSections) CreateConfigsCommon() {
	// commonConfigSections maps section name to section XML chopConfig of the following sections:
	// 1. remote servers
	// 2. zookeeper
	// 3. settings
	// 4. files
	util.IncludeNonEmpty(c.commonConfigSections, filenameRemoteServersXML, c.chConfigGenerator.GetRemoteServers())
	util.IncludeNonEmpty(c.commonConfigSections, filenameZookeeperXML, c.chConfigGenerator.GetZookeeper())
	util.IncludeNonEmpty(c.commonConfigSections, filenameSettingsXML, c.chConfigGenerator.GetSettings())
	util.MergeStringMaps(c.commonConfigSections, c.chConfigGenerator.GetFiles())
	// Extra user-specified config files
	util.MergeStringMaps(c.commonConfigSections, c.chopConfig.ChCommonConfigs)
}

func (c *configSections) CreateConfigsUsers() {
	// commonConfigSections maps section name to section XML chopConfig of the following sections:
	// 1. users
	// 2. quotas
	// 3. profiles
	util.IncludeNonEmpty(c.commonUsersConfigSections, filenameUsersXML, c.chConfigGenerator.GetUsers())
	util.IncludeNonEmpty(c.commonUsersConfigSections, filenameQuotasXML, c.chConfigGenerator.GetQuotas())
	util.IncludeNonEmpty(c.commonUsersConfigSections, filenameProfilesXML, c.chConfigGenerator.GetProfiles())
	// Extra user-specified config files
	util.MergeStringMaps(c.commonUsersConfigSections, c.chopConfig.ChUsersConfigs)
}

func (c *configSections) CreateConfigsPod(host *v1.ChiHost) map[string]string {
	// Prepare for this replica deployment chopConfig files map as filename->content
	podConfigSections := make(map[string]string)
	util.IncludeNonEmpty(podConfigSections, filenameMacrosXML, c.chConfigGenerator.GetHostMacros(host))
	// Extra user-specified config files
	util.MergeStringMaps(podConfigSections, c.chopConfig.ChHostConfigs)

	return podConfigSections
}
