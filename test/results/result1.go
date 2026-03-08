package results

type Result1 struct {
	X int
}

func (r Result1) LastOne() bool {
	return true
}
