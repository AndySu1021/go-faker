package model

type Demo struct{}

func (m Demo) TableName() string {
	return "demo"
}

func (m Demo) Definition() map[string]interface{} {
	return map[string]interface{}{
		"column": "value",
	}
}
