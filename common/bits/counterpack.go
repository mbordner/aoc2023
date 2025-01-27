package bits

import "fmt"

type CounterPack uint64

const (
	MASK            uint64 = 0xFF
	MaxCounterValue        = uint8(0xFF)
	CounterBits            = 8
	MaxCounters            = 8
)

func NewCounterPack(values []int) (CounterPack, error) {
	if len(values) >= MaxCounters {
		return 0, fmt.Errorf("invalid number of values: %d", len(values))
	}

	var pack uint64
	for i := 0; i < len(values); i++ {
		if values[i] > int(MaxCounterValue) {
			return 0, fmt.Errorf("value %d exceeds max value", values[i])
		}
		pack |= uint64(values[i]) << (i * CounterBits)
	}
	return CounterPack(pack), nil
}

func (cp CounterPack) GetCount(index int) (int, error) {
	if index < 0 || index >= MaxCounters {
		return 0, fmt.Errorf("invalid index: %d", index)
	}

	shift := uint(index * CounterBits)
	shifted := uint64(cp) >> shift

	return int(shifted & MASK), nil
}

func (cp CounterPack) SetCounter(index int, value int) (CounterPack, error) {
	if index < 0 || index >= MaxCounters {
		return cp, fmt.Errorf("invalid index: %d", index)
	}

	if value > int(MaxCounterValue) {
		return cp, fmt.Errorf("value %d exceeds max value", value)
	}

	shift := uint(index * CounterBits)
	clearMask := ^(MASK << shift)

	newCP := uint64(cp) & clearMask
	newCounterBlock := uint64(value) << shift

	return CounterPack(newCP | newCounterBlock), nil
}

func (cp CounterPack) IncrementCount(index int) (CounterPack, error) {
	if index < 0 || index >= MaxCounters {
		return cp, fmt.Errorf("invalid index: %d", index)
	}
	currentValue, _ := cp.GetCount(index)

	if currentValue < int(MaxCounterValue) {
		return cp.SetCounter(index, currentValue+1)
	}
	return cp, fmt.Errorf("count at index: %d is at max", index)
}

func (cp CounterPack) DecrementCount(index int) (CounterPack, error) {
	if index < 0 || index >= MaxCounters {
		return cp, fmt.Errorf("invalid index: %d", index)
	}
	currentValue, _ := cp.GetCount(index)

	if currentValue > 0 {
		return cp.SetCounter(index, currentValue-1)
	}
	return cp, fmt.Errorf("count at index: %d is at 0", index)
}

func (cp CounterPack) Values() []int {
	values := make([]int, 8)
	for i := 0; i < len(values); i++ {
		values[i], _ = cp.GetCount(i)
	}
	return values
}
