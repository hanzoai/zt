package edge

import (
	"github.com/hanzozt/fablab/kernel/lib/actions/component"
	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/ziti/zititest/zitilab"
)

func ReEnrollEdgeRouters(componentSpec string, concurrency int) model.Action {
	return &reEnrollEdgeRoutersAction{
		componentSpec: componentSpec,
		concurrency:   concurrency,
	}
}

func (action *reEnrollEdgeRoutersAction) Execute(run model.Run) error {
	return component.ExecInParallel(action.componentSpec, action.concurrency, zitilab.RouterActionsReEnroll).Execute(run)
}

type reEnrollEdgeRoutersAction struct {
	componentSpec string
	concurrency   int
}
