// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fileconsumer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFinder(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		files          []string
		include        []string
		exclude        []string
		filterSortRule OrderingCriteria
		expected       []string
	}{
		{
			name:     "IncludeOne",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"a1.log"},
			exclude:  []string{},
			expected: []string{"a1.log"},
		},
		{
			name:     "IncludeNone",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"c*.log"},
			exclude:  []string{},
			expected: []string{},
		},
		{
			name:     "IncludeAll",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"*"},
			exclude:  []string{},
			expected: []string{"a1.log", "a2.log", "b1.log", "b2.log"},
		},
		{
			name:     "IncludeLogs",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"*.log"},
			exclude:  []string{},
			expected: []string{"a1.log", "a2.log", "b1.log", "b2.log"},
		},
		{
			name:     "IncludeA",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"a*.log"},
			exclude:  []string{},
			expected: []string{"a1.log", "a2.log"},
		},
		{
			name:     "Include2s",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"*2.log"},
			exclude:  []string{},
			expected: []string{"a2.log", "b2.log"},
		},
		{
			name:     "Exclude",
			files:    []string{"include.log", "exclude.log"},
			include:  []string{"*"},
			exclude:  []string{"exclude.log"},
			expected: []string{"include.log"},
		},
		{
			name:     "ExcludeMany",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"*"},
			exclude:  []string{"a*.log", "*2.log"},
			expected: []string{"b1.log"},
		},
		{
			name:     "ExcludeDuplicates",
			files:    []string{"a1.log", "a2.log", "b1.log", "b2.log"},
			include:  []string{"*1*", "a*"},
			exclude:  []string{"a*.log", "*2.log"},
			expected: []string{"b1.log"},
		},
		{
			name:     "IncludeMultipleDirectories",
			files:    []string{filepath.Join("a", "1.log"), filepath.Join("a", "2.log"), filepath.Join("b", "1.log"), filepath.Join("b", "2.log")},
			include:  []string{filepath.Join("a", "*.log"), filepath.Join("b", "*.log")},
			exclude:  []string{},
			expected: []string{filepath.Join("a", "1.log"), filepath.Join("a", "2.log"), filepath.Join("b", "1.log"), filepath.Join("b", "2.log")},
		},
		{
			name:     "IncludeMultipleDirectoriesVaryingDepth",
			files:    []string{"1.log", filepath.Join("a", "1.log"), filepath.Join("a", "b", "1.log"), filepath.Join("c", "1.log")},
			include:  []string{"*.log", filepath.Join("a", "*.log"), filepath.Join("a", "b", "*.log"), filepath.Join("c", "*.log")},
			exclude:  []string{},
			expected: []string{"1.log", filepath.Join("a", "1.log"), filepath.Join("a", "b", "1.log"), filepath.Join("c", "1.log")},
		},
		{
			name:     "DoubleStarSameDepth",
			files:    []string{filepath.Join("a", "1.log"), filepath.Join("b", "1.log"), filepath.Join("c", "1.log")},
			include:  []string{filepath.Join("**", "*.log")},
			exclude:  []string{},
			expected: []string{filepath.Join("a", "1.log"), filepath.Join("b", "1.log"), filepath.Join("c", "1.log")},
		},
		{
			name:     "DoubleStarVaryingDepth",
			files:    []string{"1.log", filepath.Join("a", "1.log"), filepath.Join("a", "b", "1.log"), filepath.Join("c", "1.log")},
			include:  []string{filepath.Join("**", "*.log")},
			exclude:  []string{},
			expected: []string{"1.log", filepath.Join("a", "1.log"), filepath.Join("a", "b", "1.log"), filepath.Join("c", "1.log")},
		},
		{
			name:     "SingleLevelFilesOnly",
			files:    []string{"a1.log", "a2.txt", "b/b1.log", "b/b2.txt"},
			include:  []string{"*"},
			expected: []string{"a1.log", "a2.txt"},
		},
		{
			name:     "MultiLevelFilesOnly",
			files:    []string{"a1.log", "a2.txt", "b/b1.log", "b/b2.txt", "b/c/c1.csv"},
			include:  []string{filepath.Join("**", "*")},
			expected: []string{"a1.log", "a2.txt", filepath.Join("b", "b1.log"), filepath.Join("b", "b2.txt"), filepath.Join("b", "c", "c1.csv")},
		},
		{
			name:    "Timestamp Sorting",
			files:   []string{"err.2023020611.log", "err.2023020612.log", "err.2023020610.log", "err.2023020609.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>\d{4}\d{2}\d{2}\d{2}).*log`,
				SortBy: []SortRuleImpl{
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
				},
			},
			expected: []string{"err.2023020612.log"},
		},
		{
			name:    "Timestamp Sorting Ascending",
			files:   []string{"err.2023020612.log", "err.2023020611.log", "err.2023020609.log", "err.2023020610.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>\d{4}\d{2}\d{2}\d{2}).*log`,
				SortBy: []SortRuleImpl{
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: true,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
				},
			},
			expected: []string{"err.2023020609.log"},
		},
		{
			name:    "Numeric Sorting",
			files:   []string{"err.123456788.log", "err.123456789.log", "err.123456787.log", "err.123456786.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>\d+).*log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: false,
							},
						},
					},
				},
			},
			expected: []string{"err.123456789.log"},
		},
		{
			name:    "Numeric Sorting Ascending",
			files:   []string{"err.123456789.log", "err.123456788.log", "err.123456786.log", "err.123456787.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>\d+).*log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: true,
							},
						},
					},
				},
			},
			expected: []string{"err.123456786.log"},
		},
		{
			name:    "Alphabetical Sorting",
			files:   []string{"err.a.log", "err.d.log", "err.b.log", "err.c.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>[a-zA-Z]+).*log`,
				SortBy: []SortRuleImpl{
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: false,
							},
						},
					},
				},
			},
			expected: []string{"err.d.log"},
		},
		{
			name:    "Alphabetical Sorting Ascending",
			files:   []string{"err.b.log", "err.a.log", "err.c.log", "err.d.log"},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<value>[a-zA-Z]+).*log`,
				SortBy: []SortRuleImpl{
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "value",
								Ascending: true,
							},
						},
					},
				},
			},
			expected: []string{"err.a.log"},
		},
		{
			name: "Multiple Sorting - timestamp priority sort",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: false,
							},
						},
					},
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: false,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
				},
			},
			expected: []string{"err.b.2.2023020602.log"},
		},
		{
			name: "Multiple Sorting - timestamp priority sort - numeric ascending",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: false,
							},
						},
					},
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: true,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
				},
			},
			expected: []string{"err.b.1.2023020602.log"},
		},
		{
			name: "Multiple Sorting - timestamp priority sort",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: false,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: false,
							},
						},
					},
				},
			},
			expected: []string{"err.b.2.2023020602.log"},
		},
		{
			name: "Multiple Sorting - alpha priority sort - alpha ascending",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: false,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: true,
							},
						},
					},
				},
			},
			expected: []string{"err.a.2.2023020602.log"},
		},
		{
			name: "Multiple Sorting - alpha priority sort - timestamp ascending",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: false,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: true,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: false,
							},
						},
					},
				},
			},
			expected: []string{"err.b.2.2023020601.log"},
		},
		{
			name: "Multiple Sorting - alpha priority sort - timestamp ascending",
			files: []string{
				"err.b.1.2023020601.log",
				"err.b.2.2023020601.log",
				"err.a.1.2023020601.log",
				"err.a.2.2023020601.log",
				"err.b.1.2023020602.log",
				"err.a.2.2023020602.log",
				"err.b.2.2023020602.log",
				"err.a.1.2023020602.log",
			},
			include: []string{"err.*.log"},
			exclude: []string{},
			filterSortRule: OrderingCriteria{
				Regex: `err\.(?P<alpha>[a-zA-Z])\.(?P<number>\d+)\.(?P<time>\d{10})\.log`,
				SortBy: []SortRuleImpl{
					{
						&NumericSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "number",
								Ascending: true,
							},
						},
					},
					{
						&TimestampSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "time",
								Ascending: false,
							},
							Location: "UTC",
							Layout:   `%Y%m%d%H`,
						},
					},
					{
						&AlphabeticalSortRule{
							BaseSortRule: BaseSortRule{
								RegexKey:  "alpha",
								Ascending: false,
							},
						},
					},
				},
			},
			expected: []string{"err.b.1.2023020602.log"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			files := absPath(tempDir, tc.files)
			include := absPath(tempDir, tc.include)
			exclude := absPath(tempDir, tc.exclude)
			expected := absPath(tempDir, tc.expected)

			for _, f := range files {
				require.NoError(t, os.MkdirAll(filepath.Dir(f), 0700))
				require.NoError(t, os.WriteFile(f, []byte(filepath.Base(f)), 0000))
			}

			finder := Finder{
				Include:          include,
				Exclude:          exclude,
				OrderingCriteria: tc.filterSortRule,
			}
			files, err := finder.FindFiles()
			require.NoError(t, err)
			require.Equal(t, expected, files)
		})
	}
}

func absPath(tempDir string, files []string) []string {
	absFiles := make([]string, 0, len(files))
	for _, f := range files {
		absFiles = append(absFiles, filepath.Join(tempDir, f))
	}
	return absFiles
}
