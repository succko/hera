package metadata

import (
	"context"
	"github.com/succko/hera/global"
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

func (loader *loader) InitializeMetadata() {
	zap.L().Info("metadata initializeMetadata start")
	funcs := global.App.RunConfig.MetaData
	if len(funcs) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(funcs))
		for _, f := range funcs {
			f := f
			go func() {
				defer wg.Done()
				f()
			}()
		}
		wg.Wait()
	}
	zap.L().Info("metadata initializeMetadata success")
	return
}
