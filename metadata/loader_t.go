package metadata

type loaderInterface interface {
	Load()
	Update(id int64)
}
