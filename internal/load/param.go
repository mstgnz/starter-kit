package load

type Param struct {
	Table      string
	Query      string
	Model      any
	Offset     int
	Limit      int
	LangCode   string
	Fields     map[string]any
	Conditions map[string]any
}
