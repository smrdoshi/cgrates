/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package config

import (
	"github.com/cgrates/cgrates/utils"
)

type RadiusAgentCfg struct {
	Enabled            bool
	ListenNet          string // udp or tcp
	ListenAuth         string
	ListenAcct         string
	ClientSecrets      map[string]string
	ClientDictionaries map[string]string
	SMGenericConns     []*HaPoolConfig
	CreateCDR          bool
	CDRRequiresSession bool
	Timezone           string
	RequestProcessors  []*RARequestProcessor
}

func (self *RadiusAgentCfg) loadFromJsonCfg(jsnCfg *RadiusAgentJsonCfg) error {
	if jsnCfg == nil {
		return nil
	}
	if jsnCfg.Enabled != nil {
		self.Enabled = *jsnCfg.Enabled
	}
	if jsnCfg.Listen_net != nil {
		self.ListenNet = *jsnCfg.Listen_net
	}
	if jsnCfg.Listen_auth != nil {
		self.ListenAuth = *jsnCfg.Listen_auth
	}
	if jsnCfg.Listen_acct != nil {
		self.ListenAcct = *jsnCfg.Listen_acct
	}
	if jsnCfg.Client_secrets != nil {
		if self.ClientSecrets == nil {
			self.ClientSecrets = make(map[string]string)
		}
		for k, v := range *jsnCfg.Client_secrets {
			self.ClientSecrets[k] = v
		}
	}
	if jsnCfg.Client_dictionaries != nil {
		if self.ClientDictionaries == nil {
			self.ClientDictionaries = make(map[string]string)
		}
		for k, v := range *jsnCfg.Client_dictionaries {
			self.ClientDictionaries[k] = v
		}
	}
	if jsnCfg.Sm_generic_conns != nil {
		self.SMGenericConns = make([]*HaPoolConfig, len(*jsnCfg.Sm_generic_conns))
		for idx, jsnHaCfg := range *jsnCfg.Sm_generic_conns {
			self.SMGenericConns[idx] = NewDfltHaPoolConfig()
			self.SMGenericConns[idx].loadFromJsonCfg(jsnHaCfg)
		}
	}
	if jsnCfg.Create_cdr != nil {
		self.CreateCDR = *jsnCfg.Create_cdr
	}
	if jsnCfg.Cdr_requires_session != nil {
		self.CDRRequiresSession = *jsnCfg.Cdr_requires_session
	}
	if jsnCfg.Timezone != nil {
		self.Timezone = *jsnCfg.Timezone
	}
	if jsnCfg.Request_processors != nil {
		for _, reqProcJsn := range *jsnCfg.Request_processors {
			rp := new(RARequestProcessor)
			var haveID bool
			for _, rpSet := range self.RequestProcessors {
				if reqProcJsn.Id != nil && rpSet.Id == *reqProcJsn.Id {
					rp = rpSet // Will load data into the one set
					haveID = true
					break
				}
			}
			if err := rp.loadFromJsonCfg(reqProcJsn); err != nil {
				return nil
			}
			if !haveID {
				self.RequestProcessors = append(self.RequestProcessors, rp)
			}
		}
	}
	return nil
}

// One Diameter request processor configuration
type RARequestProcessor struct {
	Id                string
	DryRun            bool
	RequestFilter     utils.RSRFields
	Flags             utils.StringMap // Various flags to influence behavior
	ContinueOnSuccess bool
	AppendReply       bool
	RequestFields     []*CfgCdrField
	ReplyFields       []*CfgCdrField
}

func (self *RARequestProcessor) loadFromJsonCfg(jsnCfg *RAReqProcessorJsnCfg) error {
	if jsnCfg == nil {
		return nil
	}
	if jsnCfg.Id != nil {
		self.Id = *jsnCfg.Id
	}
	if jsnCfg.Dry_run != nil {
		self.DryRun = *jsnCfg.Dry_run
	}
	var err error
	if jsnCfg.Request_filter != nil {
		if self.RequestFilter, err = utils.ParseRSRFields(*jsnCfg.Request_filter, utils.INFIELD_SEP); err != nil {
			return err
		}
	}
	if jsnCfg.Flags != nil {
		self.Flags = utils.StringMapFromSlice(*jsnCfg.Flags)
	}
	if jsnCfg.Continue_on_success != nil {
		self.ContinueOnSuccess = *jsnCfg.Continue_on_success
	}
	if jsnCfg.Append_reply != nil {
		self.AppendReply = *jsnCfg.Append_reply
	}
	if jsnCfg.Request_fields != nil {
		if self.RequestFields, err = CfgCdrFieldsFromCdrFieldsJsonCfg(*jsnCfg.Request_fields); err != nil {
			return err
		}
	}
	if jsnCfg.Reply_fields != nil {
		if self.ReplyFields, err = CfgCdrFieldsFromCdrFieldsJsonCfg(*jsnCfg.Reply_fields); err != nil {
			return err
		}
	}
	return nil
}
