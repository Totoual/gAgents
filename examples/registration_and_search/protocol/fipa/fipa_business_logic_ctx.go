package fipa

type BusinessLogic interface {
	Apply(Performative, *FipaContent) *BusinessLogicContext
}

type BusinessLogicContext interface {
	GetPerformative() Performative
	GetContent() *FipaContent
}
