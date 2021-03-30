package orm

import "testing"

func TestConnection(t *testing.T) {
	t.Run("connection - OK", func(t *testing.T) {
		conn := Connection()

		conn.NowFunc()
	})
}
