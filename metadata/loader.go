package metadata

import (
	"context"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"
	"sync"
)

type loader struct {
}

var Loader = new(loader)

func (loader *loader) Init(cxt context.Context, param *xxl.RunReq) (msg string) {
	zap.L().Info("metadata init xxl task")
	//loader.InitializeMetadata()
	return
}

func (loader *loader) InitializeMetadata(funcs []func(wg *sync.WaitGroup)) {
	zap.L().Info("metadata initializeMetadata start")
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		go f(&wg)
	}
	wg.Wait()
	zap.L().Info("metadata initializeMetadata success")
	return
}
