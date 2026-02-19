package edge

import (
	"github.com/hanzozt/fablab/kernel/lib/actions/component"
	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/zt/zttest/ztlab"
)

func ReEnrollIdentities(componentSpec string, concurrency int) model.Action {
	return &reEnrollIdentitiesAction{
		componentSpec: componentSpec,
		concurrency:   concurrency,
	}
}

func (action *reEnrollIdentitiesAction) Execute(run model.Run) error {
	return component.ExecInParallel(action.componentSpec, action.concurrency, ztlab.ZitiTunnelActionsReEnroll).Execute(run)
}

type reEnrollIdentitiesAction struct {
	componentSpec string
	concurrency   int
}
