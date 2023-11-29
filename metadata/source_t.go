package metadata

type sourceInterface interface {
	List() []interface{}
	One(id int64) (interface{}, error)
}
