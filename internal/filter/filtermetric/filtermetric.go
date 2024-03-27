// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package filtermetric // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/filtermetric"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/expr"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
)

// NewSkipExpr creates a BoolExpr that on evaluation returns true if a metric should NOT be processed or kept.
// The logic determining if a metric should be processed is based on include and exclude settings.
// Include properties are checked before exclude settings are checked.
func NewSkipExpr(include *MatchProperties, exclude *MatchProperties) (expr.BoolExpr[ottlmetric.TransformContext], error) {
	var matchers []expr.BoolExpr[ottlmetric.TransformContext]
	inclExpr, err := newExpr(include)
	if err != nil {
		return nil, err
	}
	if inclExpr != nil {
		matchers = append(matchers, expr.Not(inclExpr))
	}
	exclExpr, err := newExpr(exclude)
	if err != nil {
		return nil, err
	}
	if exclExpr != nil {
		matchers = append(matchers, exclExpr)
	}
	return expr.Or(matchers...), nil
}

// NewMatcher constructs a metric Matcher. If an 'expr' match type is specified,
// returns an expr matcher, otherwise a name matcher.
func newExpr(mp *MatchProperties) (expr.BoolExpr[ottlmetric.TransformContext], error) {
	if mp == nil {
		return nil, nil
	}

	if mp.MatchType == Expr {
		if len(mp.Expressions) == 0 {
			return nil, nil
		}
		return newExprMatcher(mp.Expressions)
	}
	if len(mp.MetricNames) == 0 {
		return nil, nil
	}
	return newNameMatcher(mp)
}
