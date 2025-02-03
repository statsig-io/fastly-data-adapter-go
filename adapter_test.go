package statsig_fastly_adapter

import (
	"fmt"
	"testing"
)

func TestFastlyDataAdapter_Get(t *testing.T) {
	adapter := NewFastlyDataAdapter("a", "b", "c")
	response, _ := adapter.getData()
	fmt.Println("Error:", response)
	testCases := []struct {
		key      string
		expected bool
	}{
		{"statsig.cache", true},
		{"statsig|/v1/download_config_specs|random|key", true},
		{"statsig|/v2/download_config_specs|random|key", true},
		{"invalid-key", false},
	}

	for _, tc := range testCases {
		result := adapter.isConfigSpecKey(tc.key)

		if result != tc.expected {
			t.Errorf("For key %s, expected %v but got %v", tc.key, tc.expected, result)
		}
	}
}
