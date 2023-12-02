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

type ProposalContentType struct {
	Cart []Item `json:"cart"`
}

func (p ProposalContentType) GetCart() *[]Item {
	return &p.Cart
}

type RejectContentType struct {
	Cart []Item `json:"cart"`
}

func (r RejectContentType) GetCart() *[]Item {
	return &r.Cart
}
