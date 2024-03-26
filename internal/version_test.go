package internal

import (
	"os"
	"testing"

	"github.com/go-quicktest/qt"
)

func TestVersion(t *testing.T) {
	a, err := NewVersion("1.2")
	if err != nil {
		t.Fatal(err)
	}

	b, err := NewVersion("2.2.1")
	if err != nil {
		t.Fatal(err)
	}

	if !a.Less(b) {
		t.Error("A should be less than B")
	}

	if b.Less(a) {
		t.Error("B shouldn't be less than A")
	}

	v200 := Version{2, 0, 0}
	if !a.Less(v200) {
		t.Error("1.2.1 should not be less than 2.0.0")
	}

	if v200.Less(a) {
		t.Error("2.0.0 should not be less than 1.2.1")
	}
}

func TestKernelVersion(t *testing.T) {
	// Kernels 4.4 and 4.9 have a SUBLEVEL of over 255 and clamp it to 255.
	// In our implementation, the other version segments are truncated.
	if v, want := (Version{256, 256, 256}), uint32(255); v.Kernel() != want {
		t.Errorf("256.256.256 should result in a kernel version of %d, got: %d", want, v.Kernel())
	}

	// Known good version.
	if v, want := (Version{4, 9, 128}), uint32(264576); v.Kernel() != want {
		t.Errorf("4.9.1 should result in a kernel version of %d, got: %d", want, v.Kernel())
	}
}

func TestCurrentKernelVersion(t *testing.T) {
	v, err := KernelVersion()
	if err != nil {
		t.Fatal(err)
	}
	if evStr := os.Getenv("KERNEL_VERSION"); evStr != "" {
		ev, err := NewVersion(evStr)
		if err != nil {
			t.Fatal(err)
		}
		if ev[0] != v[0] || ev[1] != v[1] {
			t.Errorf("expected kernel version %d.%d, got %d.%d", ev[0], ev[1], v[0], v[1])
		}
	}
}

func TestVersionFromCode(t *testing.T) {
	var tests = []struct {
		name string
		code uint32
		v    Version
	}{
		{"0.0.0", 0, Version{0, 0, 0}},
		{"1.0.0", 0x10000, Version{1, 0, 0}},
		{"4.4.255", 0x404ff, Version{4, 4, 255}},
		{"255.255.255", 0xffffff, Version{255, 255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVersionFromCode(tt.code)
			if v != tt.v {
				t.Errorf("unexpected version for code '%d'. got: %v, want: %v", tt.code, v, tt.v)
			}
		})
	}
}

func TestNewVersionFromKernelRelease(t *testing.T) {
	tests := []struct {
		release string
		version Version
		err     string
	}{
		{"6.7.9-200.fc39.x86_64", Version{6, 7, 9}, ""},
		{"5.4.0-65-generic", Version{5, 4, 0}, ""},
		{"5.4", Version{}, "invalid version: 5.4"},
	}

	for _, test := range tests {
		t.Run(test.release, func(t *testing.T) {
			v, err := NewVersionFromKernelRelease(test.release)
			if test.err != "" {
				qt.Assert(t, qt.ErrorMatches(err, test.err))
				return
			}
			qt.Assert(t, qt.IsNil(err))
			qt.Assert(t, qt.DeepEquals(v, test.version))
		})
	}
}

func TestKernelRelease(t *testing.T) {
	r, err := KernelRelease()
	if err != nil {
		t.Fatal(err)
	}

	if r == "" {
		t.Fatal("unexpected empty kernel release")
	}
}
