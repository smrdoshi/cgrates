// +build integration

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
package v1

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

var (
	tSv1CfgPath string
	tSv1Cfg     *config.CGRConfig
	tSv1Rpc     *rpc.Client
	tPrfl       *engine.ThresholdProfile
	tSv1ConfDIR string //run tests for specific configuration
	thdsDelay   int
)

var tEvs = []*engine.ThresholdEvent{
	&engine.ThresholdEvent{
		Tenant: "cgrates.org",
		ID:     "event1",
		Fields: map[string]interface{}{
			utils.EventSource: utils.StatService,
			utils.StatID:      "Stats1",
			utils.MetaASR:     35.0,
			utils.MetaACD:     time.Duration(2*time.Minute + 45*time.Second),
			utils.MetaTCC:     12.7,
			utils.MetaTCD:     time.Duration(12*time.Minute + 15*time.Second),
			utils.MetaACC:     0.75,
			utils.MetaPDD:     time.Duration(2 * time.Second),
		}},
	&engine.ThresholdEvent{
		Tenant: "cgrates.org",
		ID:     "event2",
		Fields: map[string]interface{}{
			utils.EventSource:  utils.AccountService,
			utils.AccountID:    "1002",
			utils.BalanceType:  utils.MONETARY,
			utils.BalanceID:    utils.META_DEFAULT,
			utils.BalanceValue: 12.3}},
	&engine.ThresholdEvent{
		Tenant: "cgrates.org",
		ID:     "event3",
		Fields: map[string]interface{}{
			utils.EventSource: utils.ResourceS,
			utils.ResourceID:  "ResGroup1",
			utils.USAGE:       10.0}},
}

var sTestsThresholdSV1 = []func(t *testing.T){
	testV1TSLoadConfig,
	testV1TSInitDataDb,
	testV1TSStartEngine,
	testV1TSRpcConn,
	testV1TSFromFolder,
	testV1TSGetThresholds,
	//testV1STSProcessEvent,
	//testV1TSGetThresholdsAfterRestart,
	//testV1STSSetThresholdProfile,
	//testV1STSUpdateThresholdProfile,
	//testV1STSRemoveThresholdProfile,
	testV1TSStopEngine,
}

// Test start here
func TestTSV1ITMySQL(t *testing.T) {
	tSv1ConfDIR = "tutmysql"
	for _, stest := range sTestsThresholdSV1 {
		t.Run(tSv1ConfDIR, stest)
	}
}

func TestTSV1ITMongo(t *testing.T) {
	tSv1ConfDIR = "tutmongo"
	time.Sleep(time.Duration(5 * time.Second)) // give time for engine to start
	for _, stest := range sTestsThresholdSV1 {
		t.Run(tSv1ConfDIR, stest)
	}
}

func testV1TSLoadConfig(t *testing.T) {
	var err error
	tSv1CfgPath = path.Join(*dataDir, "conf", "samples", tSv1ConfDIR)
	if tSv1Cfg, err = config.NewCGRConfigFromFolder(tSv1CfgPath); err != nil {
		t.Error(err)
	}
	switch tSv1ConfDIR {
	case "tutmongo": // Mongo needs more time to reset db, need to investigate
		thdsDelay = 4000
	default:
		thdsDelay = 1000
	}
}

func testV1TSInitDataDb(t *testing.T) {
	if err := engine.InitDataDb(tSv1Cfg); err != nil {
		t.Fatal(err)
	}
}

func testV1TSStartEngine(t *testing.T) {
	if _, err := engine.StopStartEngine(tSv1CfgPath, thdsDelay); err != nil {
		t.Fatal(err)
	}
}

func testV1TSRpcConn(t *testing.T) {
	var err error
	tSv1Rpc, err = jsonrpc.Dial("tcp", tSv1Cfg.RPCJSONListen) // We connect over JSON so we can also troubleshoot if needed
	if err != nil {
		t.Fatal("Could not connect to rater: ", err.Error())
	}
}

func testV1TSFromFolder(t *testing.T) {
	var reply string
	attrs := &utils.AttrLoadTpFromFolder{FolderPath: path.Join(*dataDir, "tariffplans", "tutorial")}
	if err := tSv1Rpc.Call("ApierV1.LoadTariffPlanFromFolder", attrs, &reply); err != nil {
		t.Error(err)
	}
}

func testV1TSGetThresholds(t *testing.T) {
	var tIDs []string
	expectedIDs := []string{"Threshold1"}
	if err := tSv1Rpc.Call("ThresholdSV1.GetThresholdIDs", "cgrates.org", &tIDs); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expectedIDs, tIDs) {
		t.Errorf("expecting: %+v, received reply: %s", expectedIDs, tIDs)
	}
	var td engine.Threshold
	eTd := engine.Threshold{Tenant: "cgrates.org", ID: expectedIDs[0]}
	if err := tSv1Rpc.Call("ThresholdSV1.GetThreshold",
		&utils.TenantID{Tenant: "cgrates.org", ID: expectedIDs[0]}, &td); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eTd, td) {
		t.Errorf("expecting: %+v, received: %+v", eTd, td)
	}
}

/*

func testV1STSProcessEvent(t *testing.T) {
	var reply string
	ev1 := engine.StatEvent{
		Tenant: "cgrates.org",
		ID:     "event1",
		Fields: map[string]interface{}{
			utils.ACCOUNT:     "1001",
			utils.ANSWER_TIME: time.Date(2014, 7, 14, 14, 25, 0, 0, time.UTC),
			utils.USAGE:       time.Duration(135 * time.Second),
			utils.COST:        123.0,
			utils.PDD:         time.Duration(12 * time.Second)}}
	if err := tSv1Rpc.Call("StatSV1.ProcessEvent", &ev1, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("received reply: %s", reply)
	}
	//process with one event (should be N/A becaus MinItems is 2)
	expectedMetrics := map[string]string{
		utils.MetaASR: utils.NOT_AVAILABLE,
		utils.MetaACD: utils.NOT_AVAILABLE,
		utils.MetaTCC: utils.NOT_AVAILABLE,
		utils.MetaTCD: utils.NOT_AVAILABLE,
		utils.MetaACC: utils.NOT_AVAILABLE,
		utils.MetaPDD: utils.NOT_AVAILABLE,
	}
	var metrics map[string]string
	if err := tSv1Rpc.Call("StatSV1.GetQueueStringMetrics", &utils.TenantID{Tenant: "cgrates.org", ID: "Stats1"}, &metrics); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expectedMetrics, metrics) {
		t.Errorf("expecting: %+v, received reply: %s", expectedMetrics, metrics)
	}
	ev2 := engine.StatEvent{
		Tenant: "cgrates.org",
		ID:     "event2",
		Fields: map[string]interface{}{
			utils.ACCOUNT:     "1002",
			utils.ANSWER_TIME: time.Date(2014, 7, 14, 14, 25, 0, 0, time.UTC),
			utils.USAGE:       time.Duration(45 * time.Second)}}
	if err := tSv1Rpc.Call("StatSV1.ProcessEvent", &ev2, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("received reply: %s", reply)
	}
	ev3 := &engine.StatEvent{
		Tenant: "cgrates.org",
		ID:     "event3",
		Fields: map[string]interface{}{
			utils.ACCOUNT:    "1002",
			utils.SETUP_TIME: time.Date(2014, 7, 14, 14, 25, 0, 0, time.UTC),
			utils.USAGE:      0}}
	if err := tSv1Rpc.Call("StatSV1.ProcessEvent", &ev3, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("received reply: %s", reply)
	}
	expectedMetrics2 := map[string]string{
		utils.MetaASR: "66.66667%",
		utils.MetaACD: "1m30s",
		utils.MetaACC: "61.5",
		utils.MetaTCD: "3m0s",
		utils.MetaTCC: "123",
		utils.MetaPDD: "4s",
	}
	var metrics2 map[string]string
	if err := tSv1Rpc.Call("StatSV1.GetQueueStringMetrics", &utils.TenantID{Tenant: "cgrates.org", ID: "Stats1"}, &metrics2); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expectedMetrics2, metrics2) {
		t.Errorf("expecting: %+v, received reply: %s", expectedMetrics2, metrics2)
	}
}

func testV1TSGetThresholdsAfterRestart(t *testing.T) {
	time.Sleep(time.Second)
	if _, err := engine.StopStartEngine(tSv1CfgPath, thdsDelay); err != nil {
		t.Fatal(err)
	}
	var err error
	tSv1Rpc, err = jsonrpc.Dial("tcp", tSv1Cfg.RPCJSONListen) // We connect over JSON so we can also troubleshoot if needed
	if err != nil {
		t.Fatal("Could not connect to rater: ", err.Error())
	}
	//get stats metrics after restart
	expectedMetrics2 := map[string]string{
		utils.MetaASR: "66.66667%",
		utils.MetaACD: "1m30s",
		utils.MetaACC: "61.5",
		utils.MetaTCD: "3m0s",
		utils.MetaTCC: "123",
		utils.MetaPDD: "4s",
	}
	var metrics2 map[string]string
	if err := tSv1Rpc.Call("StatSV1.GetQueueStringMetrics", &utils.TenantID{Tenant: "cgrates.org", ID: "Stats1"}, &metrics2); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expectedMetrics2, metrics2) {
		t.Errorf("After restat expecting: %+v, received reply: %s", expectedMetrics2, metrics2)
	}
	time.Sleep(time.Duration(1 * time.Second))
}

func testV1STSSetThresholdProfile(t *testing.T) {
	var reply *engine.ThresholdProfile
	if err := tSv1Rpc.Call("ApierV1.GetThresholdProfile",
		&utils.TenantID{Tenant: "cgrates.org", ID: "TEST_PROFILE1"}, &reply); err == nil ||
		err.Error() != utils.ErrNotFound.Error() {
		t.Error(err)
	}
	tPrfl = &engine.ThresholdProfile{
		Tenant: "cgrates.org",
		ID:     "TEST_PROFILE1",
		Filters: []*engine.RequestFilter{
			&engine.RequestFilter{
				Type:      "type",
				FieldName: "Name",
				Values:    []string{"FilterValue1", "FilterValue2"},
			},
		},
		ActivationInterval: &utils.ActivationInterval{
			ActivationTime: time.Date(2014, 7, 14, 14, 25, 0, 0, time.UTC).Local(),
			ExpiryTime:     time.Date(2014, 7, 14, 14, 25, 0, 0, time.UTC).Local(),
		},
		QueueLength: 10,
		TTL:         time.Duration(10) * time.Second,
		Metrics:     []string{"MetricValue", "MetricValueTwo"},
		Thresholds:  []string{"Val1", "Val2"},
		Blocker:     true,
		Stored:      true,
		Weight:      20,
		MinItems:    1,
	}
	var result string
	if err := tSv1Rpc.Call("ApierV1.SetThresholdProfile", tPrfl, &result); err != nil {
		t.Error(err)
	} else if result != utils.OK {
		t.Error("Unexpected reply returned", result)
	}
	if err := tSv1Rpc.Call("ApierV1.GetThresholdProfile",
		&utils.TenantID{Tenant: "cgrates.org", ID: "TEST_PROFILE1"}, &reply); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(tPrfl, reply) {
		t.Errorf("Expecting: %+v, received: %+v", tPrfl, reply)
	}
}

func testV1STSUpdateThresholdProfile(t *testing.T) {
	var result string
	tPrfl.Filters = []*engine.RequestFilter{
		&engine.RequestFilter{
			Type:      "type",
			FieldName: "Name",
			Values:    []string{"FilterValue1", "FilterValue2"},
		},
		&engine.RequestFilter{
			Type:      "*string",
			FieldName: "Accout",
			Values:    []string{"1001", "1002"},
		},
		&engine.RequestFilter{
			Type:      "*string_prefix",
			FieldName: "Destination",
			Values:    []string{"10", "20"},
		},
	}
	if err := tSv1Rpc.Call("ApierV1.SetThresholdProfile", tPrfl, &result); err != nil {
		t.Error(err)
	} else if result != utils.OK {
		t.Error("Unexpected reply returned", result)
	}
	time.Sleep(time.Duration(1 * time.Second))
	var reply *engine.ThresholdProfile
	if err := tSv1Rpc.Call("ApierV1.GetThresholdProfile",
		&utils.TenantID{Tenant: "cgrates.org", ID: "TEST_PROFILE1"}, &reply); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(tPrfl, reply) {
		t.Errorf("Expecting: %+v, received: %+v", tPrfl, reply)
	}
}

func testV1STSRemoveThresholdProfile(t *testing.T) {
	var resp string
	if err := tSv1Rpc.Call("ApierV1.RemThresholdProfile",
		&utils.TenantID{Tenant: "cgrates.org", ID: "TEST_PROFILE1"}, &resp); err != nil {
		t.Error(err)
	} else if resp != utils.OK {
		t.Error("Unexpected reply returned", resp)
	}
	var sqp *engine.ThresholdProfile
	if err := tSv1Rpc.Call("ApierV1.GetThresholdProfile",
		&utils.TenantID{Tenant: "cgrates.org", ID: "TEST_PROFILE1"}, &sqp); err == nil || err.Error() != utils.ErrNotFound.Error() {
		t.Error(err)
	}
}
*/
func testV1TSStopEngine(t *testing.T) {
	if err := engine.KillEngine(100); err != nil {
		t.Error(err)
	}
}
