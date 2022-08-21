// Copyright 2021 The Serverless Workflow Specification Authors
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

package model

import (
	"reflect"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
	"gopkg.in/go-playground/validator.v8"
)

const (
	// EventKindConsumed ...
	EventKindConsumed EventKind = "consumed"
	// EventKindProduced ...
	EventKindProduced EventKind = "produced"
)

func init() {
	val.GetValidator().RegisterStructValidation(EventStructLevelValidation, Event{})
}

// EventStructLevelValidation custom validator for event kind consumed
func EventStructLevelValidation(v *validator.Validate, structLevel *validator.StructLevel) {
	event := structLevel.CurrentStruct.Interface().(Event)

	if event.Kind == EventKindConsumed && len(event.Type) == 0 {
		structLevel.ReportError(reflect.ValueOf(event.Type), "Type", "type", "reqtypeconsumed")
	}
}

// EventKind ...
type EventKind string

// Event ...
type Event struct {
	Common
	// Unique event name
	Name string `json:"name" validate:"required"`
	// CloudEvent source
	Source string `json:"source,omitempty"`
	// CloudEvent type
	Type string `json:"type" validate:"required"`
	// Defines the CloudEvent as either 'consumed' or 'produced' by the workflow. Default is 'consumed'
	Kind EventKind `json:"kind,omitempty"`
	// If `true`, only the Event payload is accessible to consuming Workflow states. If `false`, both event payload and context attributes should be accessible"
	DataOnly bool `json:"dataOnly,omitempty"`
	// CloudEvent correlation definitions
	Correlation []Correlation `json:"correlation,omitempty" validate:"omitempty,dive"`
}

// Correlation ...
type Correlation struct {
	// CloudEvent Extension Context Attribute name
	ContextAttributeName string `json:"contextAttributeName" validate:"required"`
	// CloudEvent Extension Context Attribute value
	ContextAttributeValue string `json:"contextAttributeValue,omitempty"`
}

// EventRef ...
// TODO: Determine if we should remove old fields or at least add validation depending on spec version.
type EventRef struct {
	// Reference to the unique name of a 'produced' event definition
	TriggerEventRef string `json:"triggerEventRef" validate:"required"` // Older spec version
	ProduceEventRef string `json:"produceEventref" validate:"required"` // 0.8
	// Reference to the unique name of a 'consumed' event definition
	ResultEventRef  string `json:"resultEventRef" validate:"required"` // Older spec version
	ConsumeEventRef string `json:"consumeEventRef" validate:"required"`
	// Maximum amount of time (ISO 8601 format) to wait for the consume event. If not defined it be set to the actionExecutionTimeout
	ConsumeEventTimeout string `json:"consumeEventTimeout,omitempty"`
	// TODO: create StringOrMap structure
	// If string type, an expression which selects parts of the states data output to become the data (payload) of the event referenced by 'triggerEventRef'. If object type, a custom object to become the data (payload) of the event referenced by 'triggerEventRef'.
	Data interface{} `json:"data,omitempty"`
	// Add additional extension context attributes to the produced event
	ContextAttributes map[string]interface{} `json:"contextAttributes,omitempty"`
}
