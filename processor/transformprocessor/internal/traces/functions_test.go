// Copyright  The OpenTelemetry Authors
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

package traces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/common"
)

func Test_newFunctionCall(t *testing.T) {
	input := ptrace.NewSpan()
	input.SetName("bear")
	attrs := pcommon.NewMap()
	attrs.InsertString("test", "1")
	attrs.InsertInt("test2", 3)
	attrs.InsertBool("test3", true)
	attrs.CopyTo(input.Attributes())

	tests := []struct {
		name string
		inv  common.Invocation
		want func(ptrace.Span)
	}{
		{
			name: "set name",
			inv: common.Invocation{
				Function: "set",
				Arguments: []common.Value{
					{
						Path: &common.Path{
							Fields: []common.Field{
								{
									Name: "name",
								},
							},
						},
					},
					{
						String: strp("cat"),
					},
				},
			},
			want: func(span ptrace.Span) {
				input.CopyTo(span)
				span.SetName("cat")
			},
		},
		{
			name: "set status.code",
			inv: common.Invocation{
				Function: "set",
				Arguments: []common.Value{
					{
						Path: &common.Path{
							Fields: []common.Field{
								{
									Name: "status",
								},
								{
									Name: "code",
								},
							},
						},
					},
					{
						Int: intp(1),
					},
				},
			},
			want: func(span ptrace.Span) {
				input.CopyTo(span)
				span.Status().SetCode(ptrace.StatusCodeOk)
			},
		},
		{
			name: "keep_keys one",
			inv: common.Invocation{
				Function: "keep_keys",
				Arguments: []common.Value{
					{
						Path: &common.Path{
							Fields: []common.Field{
								{
									Name: "attributes",
								},
							},
						},
					},
					{
						String: strp("test"),
					},
				},
			},
			want: func(span ptrace.Span) {
				input.CopyTo(span)
				span.Attributes().Clear()
				attrs := pcommon.NewMap()
				attrs.InsertString("test", "1")
				attrs.CopyTo(span.Attributes())
			},
		},
		{
			name: "keep_keys two",
			inv: common.Invocation{
				Function: "keep_keys",
				Arguments: []common.Value{
					{
						Path: &common.Path{
							Fields: []common.Field{
								{
									Name: "attributes",
								},
							},
						},
					},
					{
						String: strp("test"),
					},
					{
						String: strp("test2"),
					},
				},
			},
			want: func(span ptrace.Span) {
				input.CopyTo(span)
				span.Attributes().Clear()
				attrs := pcommon.NewMap()
				attrs.InsertString("test", "1")
				attrs.InsertInt("test2", 3)
				attrs.CopyTo(span.Attributes())
			},
		},
		{
			name: "keep_keys none",
			inv: common.Invocation{
				Function: "keep_keys",
				Arguments: []common.Value{
					{
						Path: &common.Path{
							Fields: []common.Field{
								{
									Name: "attributes",
								},
							},
						},
					},
				},
			},
			want: func(span ptrace.Span) {
				input.CopyTo(span)
				span.Attributes().Clear()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span := ptrace.NewSpan()
			input.CopyTo(span)

			evaluate, err := common.NewFunctionCall(tt.inv, DefaultFunctions(), ParsePath)
			assert.NoError(t, err)
			evaluate(spanTransformContext{
				span:     span,
				il:       pcommon.NewInstrumentationScope(),
				resource: pcommon.NewResource(),
			})

			expected := ptrace.NewSpan()
			tt.want(expected)
			assert.Equal(t, expected, span)
		})
	}
}
