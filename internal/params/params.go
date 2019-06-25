package params

import (
	"reflect"

	"github.com/gojuno/minimock"
)

type Matchers struct {
	matchFuncs []interface{}
}

func (m *Matchers) AddMatchFunc(matchFunc interface{}) {
	m.matchFuncs = append(m.matchFuncs, matchFunc)
}

func (m *Matchers) Expectation(t minimock.Tester, methodDisplayName string) *Expectation {
	return &Expectation{
		t:          t,
		methodName: methodDisplayName,
		matchers:   m,
	}
}

type Expectation struct {
	t                minimock.Tester
	methodName       string
	matchers         *Matchers
	matcherTypesUsed []reflect.Type
}

func (e *Expectation) Next(parameterValue, matchFuncPtr, expValuePtr interface{}) *Expectation {
	zeroValueOfParameterType := reflect.ValueOf(expValuePtr).Elem().Interface()
	if parameterValue != nil {
		reflect.ValueOf(expValuePtr).Elem().Set(reflect.ValueOf(parameterValue))
	}
	if !reflect.DeepEqual(parameterValue, zeroValueOfParameterType) {
		return e
	}

	matchFuncType := reflect.ValueOf(matchFuncPtr).Elem().Type()
	noMatchFunc := len(e.matchers.matchFuncs) == 0 || reflect.TypeOf(e.matchers.matchFuncs[len(e.matchers.matchFuncs)-1]) != matchFuncType
	prevMatchFuncTypeIsTheSame := len(e.matcherTypesUsed) != 0 && matchFuncType == e.matcherTypesUsed[len(e.matcherTypesUsed)-1]
	if noMatchFunc && prevMatchFuncTypeIsTheSame {
		e.t.Fatal("Ambiguous combination of minimock argument matchers and values for method '" + e.methodName + "'")
		return e
	}

	reflect.ValueOf(matchFuncPtr).Elem().Set(reflect.ValueOf(e.matchers.matchFuncs[len(e.matchers.matchFuncs)-1]))
	e.matcherTypesUsed = append(e.matcherTypesUsed, matchFuncType)
	e.matchers.matchFuncs = e.matchers.matchFuncs[1:]
	return e
}

func (e *Expectation) Done() {
	if len(e.matchers.matchFuncs) != 0 {
		e.t.Fatal("Incorrect usage of minimock argument matchers around expectations for method '" + e.methodName + "'")
	}
}
