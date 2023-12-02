package fipa

type FipaContent interface {
	GetCart() *[]Item
}

type CFPContentType struct {
	Cart  []Item  `json:"cart"`
	Total float32 `json:"total"`
}

func (cfp CFPContentType) GetCart() *[]Item {
	return &cfp.Cart
}
