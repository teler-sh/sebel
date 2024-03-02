package sebel

import "testing"

func TestIsBlacklist(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		if IsBlacklist(ErrSSLBlacklist) != true {
			t.Fail()
		}
	})

	t.Run("false", func(t *testing.T) {
		if IsBlacklist(ErrNoSSLBLData) == true {
			t.Fail()
		}
	})
}
