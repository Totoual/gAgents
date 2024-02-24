package fipa

type BusinessLogic interface {
	Apply(Performative, *FipaContent) *BusinessLogicContext
}

type BusinessLogicContext struct {
	Performative   Performative
	IsAccepted     bool        // Specific for Accept/Reject
	AdditionalInfo interface{} // Can hold any additional data specific to a context
}
