package mentalese

type ResultHandler struct {
	Condition RelationSet
	Preparation RelationSet
	Answer RelationSet
}

func (handler ResultHandler) Bind(binding Binding) ResultHandler {
	return ResultHandler{
		Condition: handler.Condition.BindSingle(binding),
		Preparation: handler.Preparation.BindSingle(binding),
		Answer: handler.Answer.BindSingle(binding),
	}
}