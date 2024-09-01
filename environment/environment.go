// environment/environment.go
package environment

type Environment struct {
	variables map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]interface{}),
	}
}

func (env *Environment) Set(name string, value interface{}) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (interface{}, bool) {
	val, ok := env.variables[name]
	return val, ok
}

func (env *Environment) Remove(name string) {
	delete(env.variables, name)
}

func (env *Environment) Exists(name string) bool {
	_, ok := env.variables[name]
	return ok
}

func (env *Environment) Variables() map[string]interface{} {
	return env.variables
}
