package simpleventsrc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeOf(t *testing.T) {
	type someInterface interface{}
	type someStruct struct{}
	var s someInterface = someStruct{}

	require.Equal(t, typeOf(someStruct{}).Name(), "someStruct")
	require.Equal(t, typeOf(&someStruct{}).Name(), "someStruct")
	require.Equal(t, typeOf(s).Name(), "someStruct")

	require.NotNil(t, reflect.New(typeOf(someStruct{})).Interface())
	require.NotNil(t, reflect.New(typeOf(&someStruct{})).Interface())
	require.NotNil(t, reflect.New(typeOf(s)).Interface())
}
