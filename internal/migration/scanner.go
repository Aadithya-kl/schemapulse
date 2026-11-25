package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	Path      string
	Name      string
	Sequence  int64
	HasSeq    bool
	Content   string
	EmptyFile bool
}

type SequenceIssue struct {
	File    string
	Message string
}

type ScanResult struct {
	Files  []File
	Issues []SequenceIssue
}

var prefixDigitsRe = regexp.MustCompile(`^(\d+)[_-]`)

func Scan(targetPath string) (*ScanResult, error) {
	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path %q: %w", targetPath, err)
	}

	var paths []string
	if !info.IsDir() {
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".sql") {
			return nil, fmt.Errorf("target file %q must end with .sql", targetPath)
		}
		paths = append(paths, targetPath)
	} else {
		entries, err := os.ReadDir(targetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %q: %w", targetPath, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
				paths = append(paths, filepath.Join(targetPath, entry.Name()))
			}
		}
		sort.Strings(paths)
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no .sql migration files found at %q", targetPath)
	}

	result := &ScanResult{
		Files: make([]File, 0, len(paths)),
	}

	seenSeq := make(map[int64]string)
	var prevSeq int64
	hasPrevSeq := false

	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %q: %w", p, err)
		}
		content := string(raw)
		name := filepath.Base(p)
		mf := File{
			Path:      p,
			Name:      name,
			Content:   content,
			EmptyFile: strings.TrimSpace(content) == "",
		}

		if m := prefixDigitsRe.FindStringSubmatch(name); len(m) == 2 {
			seq, parseErr := strconv.ParseInt(m[1], 10, 64)
			if parseErr == nil {
				mf.Sequence = seq
				mf.HasSeq = true

				if existing, dup := seenSeq[seq]; dup {
					result.Issues = append(result.Issues, SequenceIssue{
						File:    p,
						Message: fmt.Sprintf("duplicate migration sequence prefix %d (already used by %s)", seq, existing),
					})
				} else {
					seenSeq[seq] = name
				}

				if hasPrevSeq && seq < prevSeq {
					result.Issues = append(result.Issues, SequenceIssue{
						File:    p,
						Message: fmt.Sprintf("out-of-order migration sequence %d follows %d", seq, prevSeq),
					})
				}
				prevSeq = seq
				hasPrevSeq = true
			}
		} else if info.IsDir() && len(paths) > 1 {
			result.Issues = append(result.Issues, SequenceIssue{
				File:    p,
				Message: fmt.Sprintf("migration filename %q does not start with a numeric prefix (e.g., 0001_init.sql)", name),
			})
		}

		if mf.EmptyFile {
			result.Issues = append(result.Issues, SequenceIssue{
				File:    p,
				Message: "migration file is empty or contains only whitespace",
			})
		}

		result.Files = append(result.Files, mf)
	}

	return result, nil
}
