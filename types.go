package astvalidator

type Condition struct {
	Operator   string       `json:"operator,omitempty"`
	Attribute  *Attribute   `json:"attribute,omitempty"`
	Conditions []*Condition `json:"conditions,omitempty"`
}

type Attribute struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type TokenAttribute struct {
	value     string
	hasCalled bool
}
