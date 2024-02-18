package interf

type TypeInterface interface {
	ParseParams(message map[string]interface{}) (TypeInterface, error)
	ParseContent() string
}
