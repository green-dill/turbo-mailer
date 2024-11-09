package postfix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetExternalIP(t *testing.T) {
	ip, err := getExternalIP()
	require.NoError(t, err)
	t.Log(ip)
}
