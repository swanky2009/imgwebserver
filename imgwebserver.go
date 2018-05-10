package imgwebserver

import (
	"github.com/swanky2009/imgwebserver/utils"
	"sync"
	"sync/atomic"
)

type IMGWEBSERVER struct {
	sync.RWMutex

	opts atomic.Value

	waitGroup utils.WaitGroupWrapper
}

func New(opts *Options) *IMGWEBSERVER {

	utils.InitLog(opts.LogLevel)

	n := &IMGWEBSERVER{}
	n.swapOpts(opts)
	return n
}

func (n *IMGWEBSERVER) getOpts() *Options {
	return n.opts.Load().(*Options)
}

func (n *IMGWEBSERVER) swapOpts(opts *Options) {
	n.opts.Store(opts)
}

func (n *IMGWEBSERVER) Main() {

	ctx := &imgContext{n}

	imageServer := &imageServer{ctx: ctx}

	n.waitGroup.Wrap(func() {
		imageServer.Init()
		imageServer.Start()
	})
}

func (n *IMGWEBSERVER) Exit() {
	n.waitGroup.Wait()
}
