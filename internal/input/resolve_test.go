package input

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve_DefaultNoArgs(t *testing.T) {
	jobs, err := Resolve(nil)
	if err != nil {
		t.Fatalf("Resolve(nil): %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	wantIn := filepath.Join("images", "unknown.jpg")
	wantOut := filepath.Join("images", "result.jpg")
	if jobs[0].In != wantIn || jobs[0].Out != wantOut {
		t.Fatalf("got %+v, want In=%s Out=%s", jobs[0], wantIn, wantOut)
	}
}

func TestResolve_FileArgs(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.jpg")
	b := filepath.Join(dir, "sub", "unknown-01.jpg")
	if err := os.MkdirAll(filepath.Dir(b), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, a, []byte("a"))
	mustWrite(t, b, []byte("b"))

	jobs, err := Resolve([]string{a, b})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2", len(jobs))
	}
	if jobs[0].Out != filepath.Join(dir, "result_a.jpg") {
		t.Errorf("job0 Out = %s", jobs[0].Out)
	}
	if jobs[1].Out != filepath.Join(dir, "sub", "result_unknown-01.jpg") {
		t.Errorf("job1 Out = %s", jobs[1].Out)
	}
}

func TestResolve_DirectoryWalkFilters(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "album", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	photo := filepath.Join(root, "album", "photo.jpg")
	nestedJPEG := filepath.Join(nested, "shot.JPEG")
	avif := filepath.Join(root, "album", "shot.avif")
	prior := filepath.Join(root, "album", "result_photo.jpg")
	bareResult := filepath.Join(nested, "result.jpg")
	txt := filepath.Join(root, "album", "notes.txt")

	for _, p := range []string{photo, nestedJPEG, avif, prior, bareResult, txt} {
		mustWrite(t, p, []byte("x"))
	}

	jobs, err := Resolve([]string{filepath.Join(root, "album")})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	got := map[string]string{}
	for _, j := range jobs {
		got[j.In] = j.Out
	}

	if len(got) != 2 {
		t.Fatalf("got %d jobs %#v, want 2 eligible JPEGs", len(got), got)
	}
	if got[photo] != filepath.Join(root, "album", "result_photo.jpg") {
		t.Errorf("photo Out = %s", got[photo])
	}
	if got[nestedJPEG] != filepath.Join(nested, "result_shot.jpg") {
		t.Errorf("nested Out = %s", got[nestedJPEG])
	}
	for _, skipped := range []string{avif, prior, bareResult, txt} {
		if _, ok := got[skipped]; ok {
			t.Errorf("unexpected job for %s", skipped)
		}
	}
}

func TestResolve_MixedFilesAndDirs(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "solo.jpg")
	album := filepath.Join(root, "album")
	if err := os.MkdirAll(album, 0o755); err != nil {
		t.Fatal(err)
	}
	inAlbum := filepath.Join(album, "one.jpg")
	mustWrite(t, file, []byte("f"))
	mustWrite(t, inAlbum, []byte("a"))

	jobs, err := Resolve([]string{file, album})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2", len(jobs))
	}
}

func TestResolve_MissingPath(t *testing.T) {
	_, err := Resolve([]string{filepath.Join(t.TempDir(), "missing.jpg")})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestResolve_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := Resolve([]string{dir})
	if err == nil {
		t.Fatal("expected error when no eligible images")
	}
}

func mustWrite(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
