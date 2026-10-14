package pg17

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

var files = []string{"postgres.wasm.zst", "share.tar.zst", "cluster.tar.zst", "SHA256SUMS", "INPUTS.SHA256SUMS"}

func TestEmbeddedFilesExist(t *testing.T) {
	for _, name := range files {
		if _, err := FS.ReadFile(name); err != nil {
			t.Errorf("%s is not embedded: %v", name, err)
		}
	}
}

func TestChecksumsMatch(t *testing.T) {
	sums, err := FS.ReadFile("SHA256SUMS")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	sc := bufio.NewScanner(strings.NewReader(string(sums)))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 {
			t.Fatalf("bad SHA256SUMS line %q", sc.Text())
		}
		want, name := fields[0], strings.TrimPrefix(fields[1], "*")
		data, err := FS.ReadFile(name)
		if err != nil {
			t.Errorf("SHA256SUMS lists %s, which is not embedded: %v", name, err)
			continue
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != want {
			t.Errorf("%s: sha256 is %s, SHA256SUMS says %s", name, got, want)
		}
		seen[name] = true
	}
	for _, name := range files[:3] {
		if !seen[name] {
			t.Errorf("SHA256SUMS has no entry for %s", name)
		}
	}
}

func TestVersion(t *testing.T) {
	if !strings.HasPrefix(Version, "17.") {
		t.Errorf("Version = %q, want a 17.x release", Version)
	}
}
