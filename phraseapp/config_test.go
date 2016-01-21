package phraseapp

import (
	"fmt"
	"testing"
)

func TestValidateIsType(t *testing.T) {
	var t1 string = "foobar";
	var t2 int = 1;
	var t3 bool = true;
	expErrT1 := fmt.Sprintf(cfgValueErrStr, "a", t1)
	expErrT2 := fmt.Sprintf(cfgValueErrStr, "a", t2)
	expErrT3 := fmt.Sprintf(cfgValueErrStr, "a", t3)

	switch res, err := ValidateIsString("a", t1); {
	case err != nil:
		t.Errorf("didn't expect an error, got %q", err)
	case res != t1:
		t.Errorf("expected value to be %q, got %q", t1, res)
	}

	switch _, err := ValidateIsString("a", t2); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT2:
		t.Errorf("expected error to be %q, got %q", expErrT2, err)
	}

	switch _, err := ValidateIsString("a", t3); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT3:
		t.Errorf("expected error to be %q, got %q", expErrT3, err)
	}

	switch _, err := ValidateIsInt("a", t1); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT1:
		t.Errorf("expected error to be %q, got %q", expErrT1, err)
	}

	switch res, err := ValidateIsInt("a", t2); {
	case err != nil:
		t.Errorf("didn't expect an error, got %q", err)
	case res != t2:
		t.Errorf("expected value to be %q, got %q", t2, res)
	}

	switch _, err := ValidateIsInt("a", t3); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT3:
		t.Errorf("expected error to be %q, got %q", expErrT3, err)
	}

	switch _, err := ValidateIsBool("a", t1); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT1:
		t.Errorf("expected error to be %q, got %q", expErrT1, err)
	}

	switch _, err := ValidateIsBool("a", t2); {
	case err == nil:
		t.Errorf("expect an error, got none")
	case err.Error() != expErrT2:
		t.Errorf("expected error to be %q, got %q", expErrT2, err)
	}

	switch res, err := ValidateIsBool("a", t3); {
	case err != nil:
		t.Errorf("didn't expect an error, got %q", err)
	case res != t3:
		t.Errorf("expected value to be %t, got %t", t3, res)
	}
}

func TestValidateIsRawMapHappyPath(t *testing.T) {
	m := map[interface{}]interface{}{
		"foo": "bar",
		"fuu": 1,
		"few": true,
	}

	res, err := ValidateIsRawMap("a", m)
	if err != nil {
		t.Errorf("didn't expect an error, got %q", err)
	}

	if len(m) != len(res) {
		t.Errorf("expected %d elements, got %d", len(m), len(res))
	}

	for k, v := range res {
		if value, found := m[k]; !found {
			t.Errorf("expected key %q to be in source set, it wasn't", k)
		} else if value != v {
			t.Errorf("expected value of %q to be %q, got %q", k, value, v)
		}
	}
}

func TestValidateIsRawMapWithErrors(t *testing.T) {
	m := map[interface{}]interface{}{
		4: "should be error",
	}

	expErr := fmt.Sprintf(cfgKeyErrStr, "a.4", 4)
	_, err := ValidateIsRawMap("a", m)
	if err == nil {
		t.Errorf("expect an error, got none")
	} else if err.Error() != expErr {
		t.Errorf("expected error %q, got %q", expErr, err)
	}
}

func TestParseYAMLToMap(t *testing.T) {
	var a string
	var b int
	var c bool
	var d []byte
	e := map[string]interface{}{}

	err := ParseYAMLToMap(func(raw interface{}) error {
		m, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid type received")
		}
		m["a"] = "foo"
		m["b"] = 1
		m["c"] = true
		m["d"] = &struct {
			A string
			B int
		}{A: "bar", B: 2}
		m["e"] = map[interface{}]interface{}{"c": "baz", "d": 4}
		return nil
	}, map[string]interface{}{
		"a": &a,
		"b": &b,
		"c": &c,
		"d": &d,
		"e": &e,
	})
	if err != nil {
		t.Fatalf("didn't expect an error, got %q", err)
	}

	if a != "foo" {
		t.Errorf("expected %q, got %q", "foo", a)
	}

	if b != 1 {
		t.Errorf("expected %d, got %d", 1, b)
	}

	if c != true {
		t.Errorf("expected %t, got %t", true, c)
	}

	if string(d) != "a: bar\nb: 2\n" {
		t.Errorf("expected %s, got %s", "a: bar\nb: 2\n", string(d))
	}


	if val, found := e["c"]; !found {
		t.Errorf("expected e to contain key %q, it didn't", "c")
	} else if val != "baz" {
		t.Errorf("expected e['c'] to have value %q, got %q", "baz", val)
	}

	if val, found := e["d"]; !found {
		t.Errorf("expected e to contain key %q, it didn't", "d")
	} else if val != 4 {
		t.Errorf("expected e['d'] to have value %d, got %d", 4, val)
	}
}