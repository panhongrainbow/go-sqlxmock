package sqlmock

import "testing"

func Test_Check_ConvertStringFormats(t *testing.T) {
	tests := []struct {
		input    string
		caseFlag uint8
		expected string
	}{
		{"convertThisString", Case_Upper, "CONVERTTHISSTRING"},
		{"convertThisString", Case_Lower, "convertthisstring"},
		{"convertThisString", Case_Snake, "convert_This_String"},
	}

	for _, test := range tests {
		converted := ConvertStringFormats(test.input, test.caseFlag)
		if converted != test.expected {
			t.Errorf("Conversion result is incorrect for input '%s', got: '%s', want: '%s'",
				test.input, converted, test.expected)
		}
	}
}
