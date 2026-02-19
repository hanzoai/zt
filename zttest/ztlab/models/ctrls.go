package models

import (
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/zt/v2/ztrest"
	"github.com/hanzozt/zt/zttest/ztlab/chaos"
)

type CtrlClients struct {
	ctrls       []*ztrest.Clients
	ctrlMap     map[string]*ztrest.Clients
	initialized atomic.Bool
}

func (self *CtrlClients) Init(run model.Run, selector string) error {
	if !self.initialized.CompareAndSwap(false, true) {
		return nil
	}

	self.ctrlMap = map[string]*ztrest.Clients{}
	ctrls := run.GetModel().SelectComponents(selector)
	resultC := make(chan struct {
		err     error
		id      string
		clients *ztrest.Clients
	}, len(ctrls))

	for _, ctrl := range ctrls {
		go func() {
			clients, err := chaos.EnsureLoggedIntoCtrl(run, ctrl, time.Minute)
			resultC <- struct {
				err     error
				id      string
				clients *ztrest.Clients
			}{
				err:     err,
				id:      ctrl.Id,
				clients: clients,
			}
		}()
	}

	for i := 0; i < len(ctrls); i++ {
		result := <-resultC
		if result.err != nil {
			return result.err
		}
		self.ctrls = append(self.ctrls, result.clients)
		self.ctrlMap[result.id] = result.clients
	}
	return nil
}

func (self *CtrlClients) GetRandomCtrl() *ztrest.Clients {
	return self.ctrls[rand.IntN(len(self.ctrls))]
}

func (self *CtrlClients) GetCtrl(id string) *ztrest.Clients {
	return self.ctrlMap[id]
}
