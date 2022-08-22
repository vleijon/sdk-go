// Copyright 2020 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestBasicValidation(t *testing.T) {
	rootPath := "./testdata/workflows"
	files, err := ioutil.ReadDir(rootPath)
	assert.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			workflow, err := FromFile(filepath.Join(rootPath, file.Name()))
			assert.NoError(t, err, "error parsing workflow %s", file.Name())
			assert.NotEmpty(t, workflow.Name)
			assert.NotEmpty(t, workflow.ID)
			assert.NotEmpty(t, workflow.States)

			_, err = json.Marshal(workflow)
			assert.NoError(t, err)
		}
	}
}

func TestBasicValidationv08(t *testing.T) {
	rootPath := "./testdata/workflows/0.8"
	files, err := ioutil.ReadDir(rootPath)
	assert.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			filename := filepath.Join(rootPath, file.Name())
			workflow, err := FromFile(filename)
			assert.NoErrorf(t, err, "Error parsing workflow: %s (%s)", file.Name(), err)
			assert.NotNilf(t, workflow, "Workflow is nil: %s", file.Name())
			assert.NotEmpty(t, workflow.Name)
			assert.NotEmpty(t, workflow.ID)
			assert.NotEmpty(t, workflow.States)

			_, err = json.Marshal(workflow)
			assert.NoError(t, err)
		}
	}
}

func TestCustomValidators(t *testing.T) {
	rootPath := "./testdata/workflows/witherrors"
	files, err := ioutil.ReadDir(rootPath)
	assert.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			_, err := FromFile(filepath.Join(rootPath, file.Name()))
			assert.Error(t, err)
		}
	}
}

// TestFromFilev08 tests all examples for specVersion 0.8.
// The examples are taken from https://github.com/serverlessworkflow/specification/blob/main/examples/README.md
// The intent is not to run complete tests but to ensure that the examples are valid and test some key elements of them.
func TestFromFilev08(t *testing.T) {
	files := map[string]func(*testing.T, *model.Workflow){
		"./testdata/workflows/0.8/helloworld.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "helloworld", w.ID)
			assert.Equal(t, "Hello State", w.Start.StateName)
			assert.IsType(t, w.States[0], &model.InjectState{})
			assert.Equal(t, "Hello World!", w.States[0].(*model.InjectState).Data["result"])
		},
		"./testdata/workflows/0.8/greetings.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "greeting", w.ID)
			assert.IsType(t, &model.OperationState{}, w.States[0])
			assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
		},
		"./testdata/workflows/0.8/greetings.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			assert.Equal(t, "greeting", w.ID)
			assert.NotEmpty(t, w.States[0].(*model.OperationState).Actions)
			assert.NotNil(t, w.States[0].(*model.OperationState).Actions[0].FunctionRef)
			assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
		},
		"./testdata/workflows/0.8/eventbasedgreeting.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			assert.Equal(t, true, eventState.Exclusive)
		},
		"./testdata/workflows/0.8/eventbasedgreeting.sw.p.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
		},
		"./testdata/workflows/0.8/eventbasedgreeting.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			assert.Equal(t, true, eventState.Exclusive)
		},
		"./testdata/workflows/0.8/solvemathproblems.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "solvemathproblems", w.ID)
			assert.IsType(t, w.States[0], &model.ForEachState{})
			state := w.States[0].(*model.ForEachState)
			assert.Equal(t, "Solve", state.Name)
			assert.Equal(t, state.IterationParam, "singleexpression")
		},
		"./testdata/workflows/0.8/parallelexec.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "parallelexec", w.ID)
			assert.IsType(t, w.States[0], &model.ParallelState{})
		},
		"./testdata/workflows/0.8/asyncfunction.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "sendcustomeremail", w.ID)
			assert.IsType(t, w.States[0], &model.OperationState{})
			state := w.States[0].(*model.OperationState)
			assert.Equal(t, "Send Email", state.Name)
		},
		"./testdata/workflows/0.8/asyncsubflow.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "onboardcustomer", w.ID)
			assert.IsType(t, w.States[0], &model.OperationState{})
			state := w.States[0].(*model.OperationState)
			assert.Equal(t, "Onboard", state.Name)
			act := state.Actions[0]
			assert.NotNil(t, act)
			sfref := act.SubFlowRef
			assert.NotNil(t, sfref)
			assert.Equal(t, "customeronboardingworkflow", sfref.WorkflowID)
			assert.Equal(t, "1.0", sfref.Version)
			assert.EqualValues(t, "async", sfref.Invoke)
		},
		"./testdata/workflows/0.8/eventbasedtransitions.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "eventbasedswitchstate", w.ID)
			events := w.Events
			assert.Len(t, events, 2)
			assert.Equal(t, "visaApprovedEvent", events[0].Name)
			assert.Equal(t, "visaRejectedEvent", events[1].Name)
			assert.Equal(t, "VisaApproved", events[0].Type)
			assert.Equal(t, "VisaRejected", events[1].Type)
			assert.Equal(t, "visaCheckSource", events[0].Source)
			assert.Equal(t, "visaCheckSource", events[1].Source)

			switchState := w.States[0].(*model.EventBasedSwitchState)
			assert.Equal(t, "PT1H", switchState.Timeouts.EventTimeout)
			eventConds := switchState.EventConditions
			assert.Len(t, eventConds, 2)
			assert.IsType(t, &model.TransitionEventCondition{}, eventConds[0])
			assert.IsType(t, &model.TransitionEventCondition{}, eventConds[1])
			c0 := eventConds[0].(*model.TransitionEventCondition)
			c1 := eventConds[1].(*model.TransitionEventCondition)
			assert.Equal(t, c0.GetEventRef(), "visaApprovedEvent")
			assert.Equal(t, c1.GetEventRef(), "visaRejectedEvent")
			assert.Equal(t, c0.Transition.NextState, "HandleApprovedVisa")
			assert.Equal(t, c1.Transition.NextState, "HandleRejectedVisa")

			assert.Equal(t, switchState.DefaultCondition.Transition.NextState, "HandleNoVisaDecision")
		},
		"./testdata/workflows/0.8/applicantrequest.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "applicantrequest", w.ID)
			states := w.States
			assert.Len(t, states, 3)
			s3 := states[2].(*model.OperationState)
			act := s3.Actions[0]
			assert.NotNil(t, act)
			assert.Equal(t, act.FunctionRef.RefName, "sendRejectionEmailFunction")
			args := act.FunctionRef.Arguments
			assert.Contains(t, args, "applicant")
			assert.Equal(t, args["applicant"], "${ .applicant }")
		},
		"./testdata/workflows/0.8/provisionorders.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "provisionorders", w.ID)
			errors := w.Errors
			assert.Len(t, errors, 3)
			assert.Equal(t, "Missing order id", errors[0].Name)
			assert.Equal(t, "Missing order item", errors[1].Name)
			assert.Equal(t, "Missing order quantity", errors[2].Name)
			s0 := w.States[0].(*model.OperationState)
			onerr := s0.OnErrors
			assert.Len(t, onerr, 3)
			assert.Equal(t, "Missing order id", onerr[0].ErrorRef)
			assert.Equal(t, "Missing order item", onerr[1].ErrorRef)
			assert.Equal(t, "Missing order quantity", onerr[2].ErrorRef)
			assert.Equal(t, "MissingId", onerr[0].Transition.NextState)
			assert.Equal(t, "MissingItem", onerr[1].Transition.NextState)
			assert.Equal(t, "MissingQuantity", onerr[2].Transition.NextState)
			s1 := w.States[1].(*model.OperationState)
			assert.False(t, s1.End.Terminate)
		},
		"./testdata/workflows/0.8/jobmonitoring.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "jobmonitoring", w.ID)
			states := w.States
			assert.Len(t, states, 6)
			s1 := states[1]
			assert.IsType(t, &model.SleepState{}, s1)
			sleepstate := s1.(*model.SleepState)
			assert.Equal(t, "PT5S", sleepstate.Duration)
			assert.Equal(t, "GetJobStatus", sleepstate.Transition.NextState)
			s2 := states[2]
			assert.IsType(t, &model.OperationState{}, s2)
			opstate := s2.(*model.OperationState)
			assert.EqualValues(t, "sequential", opstate.ActionMode)
			assert.Equal(t, "DetermineCompletion", opstate.Transition.NextState)
		},
		"./testdata/workflows/0.8/sendcloudeventonprovision.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "sendcloudeventonprovision", w.ID)
			events := w.Events
			assert.Len(t, events, 1)
			assert.Equal(t, "provisioningCompleteEvent", events[0].Name)
			assert.Equal(t, "provisionCompleteType", events[0].Type)
			assert.EqualValues(t, "produced", events[0].Kind)
			s0 := w.States[0].(*model.ForEachState)
			assert.Equal(t, "${ .orders }", s0.InputCollection)
			assert.Equal(t, "${ .provisionedOrders }", s0.OutputCollection)
			endEvents := s0.End.ProduceEvents
			assert.Len(t, endEvents, 1)
			assert.Equal(t, "provisioningCompleteEvent", endEvents[0].EventRef)
			assert.Equal(t, "${ .provisionedOrders }", endEvents[0].Data)
		},
		"./testdata/workflows/0.8/patientVitalsWorkflow.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "patientVitalsWorkflow", w.ID)
			events := w.Events
			assert.Len(t, events, 3)
			assert.Equal(t, "HighBodyTemperature", events[0].Name)
			assert.Equal(t, "org.monitor.highBodyTemp", events[0].Type)
			assert.EqualValues(t, "patientId", events[0].Correlation[0].ContextAttributeName)
			eventState := w.States[0].(*model.EventState)
			assert.Equal(t, eventState.Exclusive, true)
			onEvents := eventState.OnEvents
			assert.Len(t, onEvents, 3)
			assert.Equal(t, "HighBodyTemperature", onEvents[0].EventRefs[0])
			assert.Equal(t, "sendTylenolOrder", onEvents[0].Actions[0].FunctionRef.RefName)
			assert.Equal(t, "${ .patientId }", onEvents[0].Actions[0].FunctionRef.Arguments["patientid"])
			assert.Equal(t, true, eventState.End.Terminate)
		},
		"./testdata/workflows/0.8/finalizeCollegeApplication.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "finalizeCollegeApplication", w.ID)
			eventState := w.States[0].(*model.EventState)
			assert.Equal(t, false, eventState.Exclusive)
			assert.Len(t, eventState.OnEvents[0].EventRefs, 3)
		},
		"./testdata/workflows/0.8/customercreditcheck.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "customercreditcheck", w.ID)
			states := w.States
			assert.Len(t, states, 4)
			callBackState := states[0].(*model.CallbackState)
			assert.Equal(t, "PT15M", callBackState.Timeouts.StateExecTimeout.Total)
			assert.Equal(t, "CreditCheckCompletedEvent", callBackState.EventRef)
			assert.Equal(t, "EvaluateDecision", callBackState.Transition.NextState)
			switchState := states[1].(*model.DataBasedSwitchState)
			conditions := switchState.DataConditions
			assert.Len(t, conditions, 2)
			c0 := conditions[0].(*model.TransitionDataCondition)
			assert.Equal(t, "${ .creditCheck | .decision == \"Approved\" }", c0.GetCondition())
			assert.Equal(t, "StartApplication", c0.Transition.NextState)
			assert.Equal(t, "RejectApplication", switchState.DefaultCondition.Transition.NextState)
		},
		"./testdata/workflows/0.8/handleCarAuctionBid.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "handleCarAuctionBid", w.ID)
			assert.Equal(t, "StoreCarAuctionBid", w.Start.StateName)
			assert.Equal(t, "R/PT2H", w.Start.Schedule.Interval)
			s0 := w.States[0].(*model.EventState)
			assert.Equal(t, true, s0.Exclusive)
			assert.Equal(t, "CarBidEvent", s0.OnEvents[0].EventRefs[0])
		},
		"./testdata/workflows/0.8/checkInbox.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "checkInbox", w.ID)
			assert.Equal(t, "0 0/15 * * * ?", w.Start.Schedule.Cron.Expression)
			assert.Equal(t, "checkInboxFunction", w.Functions[0].Name)
			assert.Equal(t, "http://myapis.org/inboxapi.json#checkNewMessages", w.Functions[0].Operation)
			foreachState := w.States[1].(*model.ForEachState)
			assert.Equal(t, "${ .messages }", foreachState.InputCollection)

		},
		"./testdata/workflows/0.8/VetAppointmentWorkflow.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "VetAppointmentWorkflow", w.ID)
			events := w.Events
			assert.Len(t, events, 2)
			assert.Equal(t, "MakeVetAppointment", events[0].Name)
			assert.Equal(t, "VetAppointmentInfo", events[1].Name)
			assert.EqualValues(t, "produced", events[0].Kind)
			assert.EqualValues(t, "consumed", events[1].Kind)
			s0 := w.States[0].(*model.OperationState)
			assert.Equal(t, "PT15M", s0.Timeouts.ActionExecTimeout)
			act := s0.Actions[0]
			assert.Equal(t, "MakeAppointmentAction", act.Name)
			assert.Equal(t, "MakeVetAppointment", act.EventRef.ProduceEventRef)
			assert.Equal(t, "VetAppointmentInfo", act.EventRef.ConsumeEventRef)
			assert.Equal(t, "${ .patientInfo }", act.EventRef.Data)
		},
		"./testdata/workflows/0.8/paymentconfirmation.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "paymentconfirmation", w.ID)
			assert.Len(t, w.Functions, 3)
			assert.Len(t, w.Events, 2)

			states := w.States
			assert.Len(t, states, 4)
			s2 := states[2].(*model.OperationState)
			assert.Equal(t, "ConfirmationCompletedEvent", s2.End.ProduceEvents[0].EventRef)
			assert.Equal(t, "${ .payment }", s2.End.ProduceEvents[0].Data)
			s3 := states[3].(*model.OperationState)
			assert.Equal(t, "ConfirmationCompletedEvent", s3.End.ProduceEvents[0].EventRef)
			assert.Equal(t, "${ .payment }", s3.End.ProduceEvents[0].Data)
		},
		"./testdata/workflows/0.8/patientonboarding.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "patientonboarding", w.ID)
			s0 := w.States[0].(*model.EventState)
			assert.Equal(t, "NewPatientEvent", s0.OnEvents[0].EventRefs[0])
			acts := s0.OnEvents[0].Actions
			assert.Len(t, acts, 3)
			assert.Equal(t, "StorePatient", acts[0].FunctionRef.RefName)
			assert.Equal(t, "ServicesNotAvailableRetryStrategy", acts[0].RetryRef)
			assert.Len(t, acts[0].RetryableErrors, 1)
			assert.Equal(t, "ServiceNotAvailable", acts[0].RetryableErrors[0])
			assert.Equal(t, "ServiceNotAvailable", s0.OnErrors[0].ErrorRef)
			errors := w.Errors
			assert.Len(t, errors, 1)
			assert.Equal(t, "ServiceNotAvailable", errors[0].Name)
			assert.Equal(t, "503", errors[0].Code)

			retries := w.Retries
			assert.Len(t, retries, 1)
			assert.Equal(t, "ServicesNotAvailableRetryStrategy", retries[0].Name)
			assert.Equal(t, "PT3S", retries[0].Delay)
			assert.EqualValues(t, 10, retries[0].MaxAttempts.IntVal)
		},
		"./testdata/workflows/0.8/order.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "order", w.ID)
			assert.Equal(t, "CancelOrder", w.Timeouts.WorkflowExecTimeout.RunBefore)
			assert.Equal(t, "PT30D", w.Timeouts.WorkflowExecTimeout.Duration)
			assert.Len(t, w.States, 4)
			assert.Len(t, w.Events, 5)
		},
		"./testdata/workflows/0.8/roomreadings.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "roomreadings", w.ID)
			assert.Equal(t, true, w.KeepActive)
			assert.Equal(t, "GenerateReport", w.Timeouts.WorkflowExecTimeout.RunBefore)
			assert.Equal(t, "PT1H", w.Timeouts.WorkflowExecTimeout.Duration)
			s0 := w.States[0].(*model.EventState)
			assert.Len(t, s0.OnEvents[0].EventRefs, 2)
		},
		"./testdata/workflows/0.8/checkcarvitals.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "checkcarvitals", w.ID)
			s1 := w.States[1].(*model.OperationState)
			actions := s1.Actions
			assert.Len(t, actions, 1)
			assert.Equal(t, "vitalscheck", actions[0].SubFlowRef.WorkflowID)
			assert.Equal(t, "PT1S", actions[0].Sleep.After)
			s2 := w.States[2].(*model.EventBasedSwitchState)
			assert.Equal(t, "DoCarVitalChecks", s2.DefaultCondition.Transition.NextState)
			eventCond := s2.EventConditions[0].(*model.EndEventCondition)
			assert.Equal(t, "CarTurnedOffEvent", eventCond.GetEventRef())
			assert.Equal(t, false, eventCond.End.Terminate)
		},
		"./testdata/workflows/0.8/booklending.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "booklending", w.ID)
		},
		"./testdata/workflows/0.8/customfunction.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "customerbankingtransactions", w.ID)
			//			fmt.Printf("%+v\n", spew.Sdump(w))
		},
	}
	for file, f := range files {
		t.Run(file, func(t *testing.T) {
			workflow, err := FromFile(file)
			assert.NoErrorf(t, err, "Test File", file)
			assert.NotNilf(t, workflow, "Test File", file)
			f(t, workflow)
		})
	}
}

func TestFromFile(t *testing.T) {
	files := map[string]func(*testing.T, *model.Workflow){
		"./testdata/workflows/eventbasedgreetingexclusive.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.Equal(t, "GreetingEvent2", w.Events[1].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			assert.Equal(t, "GreetingEvent2", eventState.OnEvents[1].EventRefs[0])
			assert.Equal(t, true, eventState.Exclusive)
		},
		"./testdata/workflows/eventbasedgreetingnonexclusive.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.Equal(t, "GreetingEvent2", w.Events[1].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			assert.Equal(t, "GreetingEvent2", eventState.OnEvents[0].EventRefs[1])
			assert.Equal(t, false, eventState.Exclusive)
		},

		"./testdata/workflows/eventbasedswitch.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.EventBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.EventBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.EventConditions)
			assert.NotEmpty(t, eventState.Name)
			assert.IsType(t, &model.TransitionEventCondition{}, eventState.EventConditions[0])
		},
		"./testdata/workflows/applicationrequest.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			assert.Equal(t, "CheckApplication", w.Start.StateName)
			assert.IsType(t, &model.OperationState{}, w.States[1])
			operationState := w.States[1].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Equal(t, "startApplicationWorkflowId", operationState.Actions[0].SubFlowRef.WorkflowID)
			assert.NotNil(t, w.Auth)
			assert.NotNil(t, w.Auth.Defs)
			assert.Equal(t, len(w.Auth.Defs), 1)
			assert.Equal(t, "testAuth", w.Auth.Defs[0].Name)
			assert.Equal(t, model.AuthTypeBearer, w.Auth.Defs[0].Scheme)
			bearerProperties := w.Auth.Defs[0].Properties.(*model.BearerAuthProperties).Token
			assert.Equal(t, "test_token", bearerProperties)
		},
		"./testdata/workflows/applicationrequest.multiauth.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			assert.Equal(t, "CheckApplication", w.Start.StateName)
			assert.IsType(t, &model.OperationState{}, w.States[1])
			operationState := w.States[1].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Equal(t, "startApplicationWorkflowId", operationState.Actions[0].SubFlowRef.WorkflowID)
			assert.NotNil(t, w.Auth)
			assert.NotNil(t, w.Auth.Defs)
			assert.Equal(t, len(w.Auth.Defs), 2)
			assert.Equal(t, "testAuth", w.Auth.Defs[0].Name)
			assert.Equal(t, model.AuthTypeBearer, w.Auth.Defs[0].Scheme)
			bearerProperties := w.Auth.Defs[0].Properties.(*model.BearerAuthProperties).Token
			assert.Equal(t, "test_token", bearerProperties)
			assert.Equal(t, "testAuth2", w.Auth.Defs[1].Name)
			assert.Equal(t, model.AuthTypeBasic, w.Auth.Defs[1].Scheme)
			basicProperties := w.Auth.Defs[1].Properties.(*model.BasicAuthProperties)
			assert.Equal(t, "test_user", basicProperties.Username)
			assert.Equal(t, "test_pwd", basicProperties.Password)

		},
		"./testdata/workflows/applicationrequest.rp.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/workflows/applicationrequest.url.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/workflows/checkinbox.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			operationState := w.States[0].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, w.States, 2)
		},
		// validates: https://github.com/serverlessworkflow/specification/pull/175/
		"./testdata/workflows/provisionorders.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			operationState := w.States[0].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, operationState.OnErrors, 3)
			assert.Equal(t, "Missing order id", operationState.OnErrors[0].ErrorRef)
			assert.Equal(t, "MissingId", operationState.OnErrors[0].Transition.NextState)
			assert.Equal(t, "Missing order item", operationState.OnErrors[1].ErrorRef)
			assert.Equal(t, "MissingItem", operationState.OnErrors[1].Transition.NextState)
			assert.Equal(t, "Missing order quantity", operationState.OnErrors[2].ErrorRef)
			assert.Equal(t, "MissingQuantity", operationState.OnErrors[2].Transition.NextState)
		}, "./testdata/workflows/checkinbox.cron-test.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "0 0/15 * * * ?", w.Start.Schedule.Cron.Expression)
			assert.Equal(t, "checkInboxFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
			assert.Equal(t, "SendTextForHighPriority", w.States[0].GetTransition().NextState)
			assert.False(t, w.States[1].GetEnd().Terminate)
		}, "./testdata/workflows/applicationrequest-issue16.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			dataBaseSwitchState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, dataBaseSwitchState)
			assert.NotEmpty(t, dataBaseSwitchState.DataConditions)
			assert.Equal(t, "CheckApplication", w.States[0].GetName())
		},
		// validates: https://github.com/serverlessworkflow/sdk-go/issues/36
		"./testdata/workflows/patientonboarding.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, w.Retries)
			assert.Len(t, w.Retries, 1)
			assert.Equal(t, float32(0.0), w.Retries[0].Jitter.FloatVal)
			assert.Equal(t, float32(1.1), w.Retries[0].Multiplier.FloatVal)
		},
		"./testdata/workflows/greetings-secret.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Len(t, w.Secrets, 1)
		},
		"./testdata/workflows/greetings-secret-file.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Len(t, w.Secrets, 3)
		},
		"./testdata/workflows/greetings-constants-file.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.NotEmpty(t, w.Constants)
			assert.NotEmpty(t, w.Constants.Data["Translations"])
		},
		"./testdata/workflows/roomreadings.timeouts.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.NotNil(t, w.Timeouts)
			assert.Equal(t, "PT1H", w.Timeouts.WorkflowExecTimeout.Duration)
			assert.Equal(t, "GenerateReport", w.Timeouts.WorkflowExecTimeout.RunBefore)
		},
		"./testdata/workflows/roomreadings.timeouts.file.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.NotNil(t, w.Timeouts)
			assert.Equal(t, "PT1H", w.Timeouts.WorkflowExecTimeout.Duration)
			assert.Equal(t, "GenerateReport", w.Timeouts.WorkflowExecTimeout.RunBefore)
		},
		"./testdata/workflows/purchaseorderworkflow.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.NotNil(t, w.Timeouts)
			assert.Equal(t, "PT30D", w.Timeouts.WorkflowExecTimeout.Duration)
			assert.Equal(t, "CancelOrder", w.Timeouts.WorkflowExecTimeout.RunBefore)
		},
	}
	for file, f := range files {
		t.Run(file, func(t *testing.T) {
			workflow, err := FromFile(file)
			assert.NoError(t, err, "Test File", file)
			assert.NotNil(t, workflow, "Test File", file)
			f(t, workflow)
		})
	}
}
