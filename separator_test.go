package yaml

import (
	"bytes"
	"testing"
)

func TestSplitYAMLDocument(t *testing.T) {
	testCases := []struct {
		input  string
		atEOF  bool
		expect string
		adv    int
	}{
		{"foo", true, "foo", 3},
		{"fo", false, "", 0},

		{"---", true, "---", 3},
		{"---\n", true, "---\n", 4},
		{"---\n", false, "", 0},

		{"\n---\n", false, "", 5},
		{"\n---\n", true, "", 5},

		{"abc\n---\ndef", true, "abc", 8},
		{"def", true, "def", 3},
		{"", true, "", 0},
	}
	for i, testCase := range testCases {
		adv, token, err := splitYAMLDocument([]byte(testCase.input), testCase.atEOF)
		if err != nil {
			t.Errorf("%d: unexpected error: %v", i, err)
			continue
		}
		if adv != testCase.adv {
			t.Errorf("%d: advance did not match: %d %d", i, testCase.adv, adv)
		}
		if testCase.expect != string(token) {
			t.Errorf("%d: token did not match: %q %q", i, testCase.expect, string(token))
		}
	}
}

func TestYAMLSeparatorNext(t *testing.T) {
	reader := bytes.NewReader([]byte(testYAMLDocFull))
	separator := NewYAMLDocumentSeparator(reader)
	doc1, err := separator.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !bytes.Equal([]byte(testYAMLDoc1), doc1) {
		t.Errorf("%s does not match %s", testYAMLDoc1, doc1)
	}
	doc2, err := separator.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !bytes.Equal([]byte(testYAMLDoc2), doc2) {
		t.Errorf("%s does not match %s", testYAMLDoc2, doc2)
	}
	extra, err := separator.Next()
	if err == nil {
		t.Error("expected EOF error but got none")
	}
	if len(extra) > 0 {
		t.Errorf("%v is not empty after the entire document has been scanned", extra)
	}
}

const testYAMLDoc1 = `
a: TestFieldA
b: TestFieldB
c:
  c1: TestFieldC1
  c2: TestFieldC2`

const testYAMLDoc2 = `
d: TestFieldD
e: TestFieldE
f:
  - name: TestFieldFListItem1
  - name: TestFieldFListItem2`

const testYAMLDocFull = testYAMLDoc1 + `
---
` + testYAMLDoc2
