package structs

import (
	"reflect"
	"testing"
)

func TestSiteConfig_DomainStructure(t *testing.T) {
	tables := []struct {
		input    string
		expected []string
	}{
		{"example.com", []string{"example.com"}},
		{"test.example.com", []string{"test", "example.com"}},
		{"", []string{}},
	}

	for _, table := range tables {
		cfg := SiteConfig{DomainName: table.input}
		result := cfg.DomainStructure()

		if !reflect.DeepEqual(result, table.expected) && len(table.expected) != 0 && len(result) != 0 {
			t.Errorf("Domain structure built incorrectly, expected %s, got %s", table.expected, result)
		}
	}
}

func TestReverseSlice(t *testing.T) {
	tables := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}},
		{[]string{"a"}, []string{"a"}},
		{[]string{}, []string{}},
	}

	for _, table := range tables {
		result := ReverseSlice(table.input)

		if !reflect.DeepEqual(result, table.expected) {
			t.Errorf("Slice reversed incorrectly, expected %s, got %s", table.expected, result)
		}
	}
}
