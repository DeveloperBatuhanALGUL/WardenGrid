package capture

import "testing"

func TestListInterfacesDoesNotError(t *testing.T) {
	_, err := ListInterfaces()
	if err != nil {
		t.Fatalf("unexpected error listing interfaces: %v", err)
	}
}
