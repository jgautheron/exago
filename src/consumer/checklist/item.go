package checklist

type CheckItemParams func(sp, sgp string) bool

type CheckItem struct {
	Name string `json:"name"`
	Desc string `json:"-"`
	fn   func() CheckItemParams
}

func (ci CheckItem) run(sp, sgp string) (success bool) {
	return ci.fn()(sp, sgp)
}
