package metadata

import "sync"

type loaderInterface interface {
	Load(wg *sync.WaitGroup)
	Update(id int64)
}
