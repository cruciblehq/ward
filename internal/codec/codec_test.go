package codec

import (
	"bytes"
	"testing"
)

func TestByteRoundTrip(t *testing.T) {
	for _, v := range []byte{0x00, 0x42, 0xFF} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteByte(v); err != nil {
			t.Fatalf("WriteByte(%d): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadByte()
		if err != nil {
			t.Fatalf("ReadByte: %v", err)
		}
		if got != v {
			t.Errorf("byte round-trip: got %d, want %d", got, v)
		}
	}
}

func TestBoolRoundTrip(t *testing.T) {
	for _, v := range []bool{true, false} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteBool(v); err != nil {
			t.Fatalf("WriteBool(%v): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadBool()
		if err != nil {
			t.Fatalf("ReadBool: %v", err)
		}
		if got != v {
			t.Errorf("bool round-trip: got %v, want %v", got, v)
		}
	}
}

func TestUint16RoundTrip(t *testing.T) {
	for _, v := range []uint16{0, 1, 0x00FF, 0xFF00, 0xFFFF} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteUint16(v); err != nil {
			t.Fatalf("WriteUint16(%d): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadUint16()
		if err != nil {
			t.Fatalf("ReadUint16: %v", err)
		}
		if got != v {
			t.Errorf("uint16 round-trip: got %d, want %d", got, v)
		}
	}
}

func TestUint32RoundTrip(t *testing.T) {
	for _, v := range []uint32{0, 1, 0xDEADBEEF, 0xFFFFFFFF} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteUint32(v); err != nil {
			t.Fatalf("WriteUint32(%d): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadUint32()
		if err != nil {
			t.Fatalf("ReadUint32: %v", err)
		}
		if got != v {
			t.Errorf("uint32 round-trip: got %d, want %d", got, v)
		}
	}
}

func TestUint64RoundTrip(t *testing.T) {
	for _, v := range []uint64{0, 1, 0xDEADBEEFCAFEBABE, 0xFFFFFFFFFFFFFFFF} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteUint64(v); err != nil {
			t.Fatalf("WriteUint64(%d): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadUint64()
		if err != nil {
			t.Fatalf("ReadUint64: %v", err)
		}
		if got != v {
			t.Errorf("uint64 round-trip: got %d, want %d", got, v)
		}
	}
}

func TestStringRoundTrip(t *testing.T) {
	for _, v := range []string{"", "hello", "café", "\x00\x01\x02"} {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteString(v); err != nil {
			t.Fatalf("WriteString(%q): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadString()
		if err != nil {
			t.Fatalf("ReadString: %v", err)
		}
		if got != v {
			t.Errorf("string round-trip: got %q, want %q", got, v)
		}
	}
}

func TestStringsRoundTrip(t *testing.T) {
	cases := [][]string{
		nil,
		{},
		{"one"},
		{"alpha", "bravo", "charlie"},
		{"", "non-empty", ""},
	}
	for _, v := range cases {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).WriteStrings(v); err != nil {
			t.Fatalf("WriteStrings(%v): %v", v, err)
		}
		got, err := NewDecoder(&buf).ReadStrings()
		if err != nil {
			t.Fatalf("ReadStrings: %v", err)
		}
		if len(v) == 0 && len(got) == 0 {
			continue
		}
		if len(got) != len(v) {
			t.Fatalf("strings round-trip: got len %d, want %d", len(got), len(v))
		}
		for i := range v {
			if got[i] != v[i] {
				t.Errorf("strings[%d]: got %q, want %q", i, got[i], v[i])
			}
		}
	}
}

func TestSequentialFields(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.WriteUint32(42); err != nil {
		t.Fatal(err)
	}
	if err := enc.WriteString("hello"); err != nil {
		t.Fatal(err)
	}
	if err := enc.WriteBool(true); err != nil {
		t.Fatal(err)
	}
	if err := enc.WriteUint64(999); err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(&buf)

	u32, err := dec.ReadUint32()
	if err != nil {
		t.Fatal(err)
	}
	if u32 != 42 {
		t.Errorf("got %d, want 42", u32)
	}

	s, err := dec.ReadString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" {
		t.Errorf("got %q, want %q", s, "hello")
	}

	b, err := dec.ReadBool()
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("got false, want true")
	}

	u64, err := dec.ReadUint64()
	if err != nil {
		t.Fatal(err)
	}
	if u64 != 999 {
		t.Errorf("got %d, want 999", u64)
	}
}

func TestEmptyStringRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteString(""); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 4 {
		t.Errorf("empty string wire size = %d, want 4", buf.Len())
	}
	got, err := NewDecoder(&buf).ReadString()
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestNilStringsRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteStrings(nil); err != nil {
		t.Fatal(err)
	}
	got, err := NewDecoder(&buf).ReadStrings()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got len %d, want 0", len(got))
	}
}
