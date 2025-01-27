package expression

import "github.com/pkg/errors"

type Variable struct {
	name string
}

func (v Variable) Eval(vars map[string]int64) (int64, error) {
	if val, exists := vars[v.name]; exists {
		return val, nil
	}
	return 0, errors.New("variable not found")
}

func (v Variable) EvalKnown(vars map[string]int64) (int64, error) {
	if val, exists := vars[v.name]; exists {
		return val, nil
	}
	return 0, errors.Errorf("unknown var %s", v.name)
}

func (v Variable) String() string {
	return v.name
}
