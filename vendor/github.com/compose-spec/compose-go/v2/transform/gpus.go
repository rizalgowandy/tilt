/*
   Copyright 2020 The Compose Specification Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package transform

import (
	"fmt"

	"github.com/compose-spec/compose-go/v2/tree"
)

func transformGpus(data any, p tree.Path, ignoreParseError bool) (any, error) {
	switch v := data.(type) {
	case []any:
		return transformSequence(v, p, ignoreParseError)
	case string:
		return []any{
			map[string]any{
				"count": "all",
			},
		}, nil
	default:
		return data, fmt.Errorf("%s: invalid type %T for gpus", p, v)
	}
}
