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

package migrator

type V1DataDB interface {
	getKeysForPrefix(prefix string) ([]string, error)
	getv1Account() (v1Acnt *v1Account, err error)
	setV1Account(x *v1Account) (err error)
	getV1ActionPlans() (v1aps *v1ActionPlans, err error)
	setV1ActionPlans(x *v1ActionPlans) (err error)
	getV1Actions() (v1acs *v1Actions, err error)
	setV1Actions(x *v1Actions) (err error)
	getV1ActionTriggers() (v1acts *v1ActionTriggers, err error)
	setV1ActionTriggers(x *v1ActionTriggers) (err error)
	getV1SharedGroup() (v1acts *v1SharedGroup, err error)
	setV1SharedGroup(x *v1SharedGroup) (err error)
}
