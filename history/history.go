package history

type History struct {
	inputs map[int]string
	cursor int
}

func NewHistory() *History {
	return &History{
		map[int]string{},
		0,
	}
}

func (h *History) Reset() *History {
	h.inputs = map[int]string{}
	h.cursor = 0

	return h
}

func (h *History) All() map[int]string {
	return h.inputs
}

func (h *History) Add(input string) *History {
	h.cursor = len(h.inputs)
	h.inputs[h.cursor] = input

	return h
}

func (h *History) Previous() *string {
	if input, ok := h.inputs[h.cursor-1]; ok {
		h.cursor--
		return &input
	}

	return nil
}

func (h *History) Next() *string {
	if input, ok := h.inputs[h.cursor+1]; ok {
		h.cursor++
		return &input
	}

	return nil
}
