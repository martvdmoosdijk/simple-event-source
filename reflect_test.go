package simpleventsrc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNameOf(t *testing.T) {
	type someInterface interface{}
	type someStruct struct{}
	var s someInterface = someStruct{}

	require.Equal(t, nameOf(someStruct{}), "someStruct")
	require.Equal(t, nameOf(&someStruct{}), "someStruct")
	require.Equal(t, nameOf(s), "someStruct")
}
