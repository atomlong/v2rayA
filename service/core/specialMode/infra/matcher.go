package infra

import (
	"github.com/xtls/xray-core/common/strmatcher"
	"sync/atomic"
)

type DomainMatcherGroup struct {
	id uint32
	g  strmatcher.DomainMatcherGroup
	strmatcher.Matcher
}

func (g *DomainMatcherGroup) Match(dm string) bool {
	return g.g.Match(dm) != nil
}

func (g *DomainMatcherGroup) Add(dm string) {
	atomic.AddUint32(&g.id, 1)
	g.g.Add(dm, g.id)
}

type FullMatcherGroup struct {
	id uint32
	g  strmatcher.FullMatcherGroup
	strmatcher.Matcher
}

func (g *FullMatcherGroup) Match(dm string) bool {
	return g.g.Match(dm) != nil
}

func (g *FullMatcherGroup) Add(dm string) {
	atomic.AddUint32(&g.id, 1)
	g.g.Add(dm, g.id)
}
