// Copyright The OpenTelemetry Authors
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

package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type DeleteKeyArguments[K any] struct {
	Target ottl.PMapGetter[K] `ottlarg:"0"`
	Key    string             `ottlarg:"1"`
}

func NewDeleteKeyFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("delete_key", &DeleteKeyArguments[K]{}, createDeleteKeyFunction[K])
}

func createDeleteKeyFunction[K any](_ ottl.FunctionContext, oArgs ottl.Arguments) (ottl.ExprFunc[K], error) {
	args, ok := oArgs.(*DeleteKeyArguments[K])

	if !ok {
		return nil, fmt.Errorf("DeleteKeysFactory args must be of type *DeleteKeyArguments[K]")
	}

	return deleteKey(args.Target, args.Key), nil
}

func deleteKey[K any](target ottl.PMapGetter[K], key string) ottl.ExprFunc[K] {
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		val, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}
		val.Remove(key)
		return nil, nil
	}
}
