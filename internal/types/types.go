package types

// FrameworkPreset defines scanning rules for a specific framework.
type FrameworkPreset struct {
	Name                      string          `json:"name"`
	DisplayName               string          `json:"displayName"`
	FileExtensions            []string        `json:"fileExtensions"`
	I18nFunctionPatterns      []string        `json:"i18nFunctionPatterns"`
	// DynamicPrefixPatterns are regex patterns (with a named "prefix" group) that
	// detect dynamic key constructions like 'SECTION.' + variable | translate.
	// Any key in the i18n files whose string prefix matches a found prefix is
	// excluded from the unused-key report and logged as a dynamic key.
	DynamicPrefixPatterns     []string        `json:"dynamicPrefixPatterns"`
	HardcodedStringPatterns   []string        `json:"hardcodedStringPatterns"`
	HardcodedStringExclusions []string        `json:"hardcodedStringExclusions"`
	I18nFilePatterns          []string        `json:"i18nFilePatterns"`
	I18nFileFormat            string          `json:"i18nFileFormat"`
	KeyStyle                  string          `json:"keyStyle"`
	KeySeparator              string          `json:"keySeparator"`
	ProjectMarkers            []ProjectMarker `json:"projectMarkers"`
	IgnorePatterns            []string        `json:"ignorePatterns"`
}

// ProjectMarker identifies a framework by checking for specific files and content.
type ProjectMarker struct {
	File           string   `json:"file"`
	ContainsAny    []string `json:"containsAny,omitempty"`
	NotContainsAny []string `json:"notContainsAny,omitempty"`
}

// UsedKey represents an i18n key found in source code.
type UsedKey struct {
	Key      string
	File     string
	Line     int
	Column   int
	RawMatch string
}

// DynamicKeyWarning represents a key that cannot be statically resolved.
type DynamicKeyWarning struct {
	RawExpression string
	File          string
	Line          int
}

// HardcodedString represents a user-facing string not wrapped in an i18n function.
type HardcodedString struct {
	Value   string
	File    string
	Line    int
	Column  int
	Context string
}

// I18nEntry represents a single key-value pair from a translation file.
type I18nEntry struct {
	Key    string
	Value  string
	File   string
	Locale string
}

// MissingKeyIssue represents a key used in code but absent from translation files.
type MissingKeyIssue struct {
	Key                string
	UsedIn             []Location
	MissingFromLocales []string
}

// UnusedKeyIssue represents a key in translation files but never used in code.
type UnusedKeyIssue struct {
	Key       string
	DefinedIn []LocaleLocation
}

// HardcodedStringIssue represents a hardcoded string grouped by value.
type HardcodedStringIssue struct {
	Value        string
	SuggestedKey string
	Occurrences  []Location
}

// Location is a file:line:column reference.
type Location struct {
	File   string
	Line   int
	Column int
}

// LocaleLocation is a file + locale reference.
type LocaleLocation struct {
	File   string
	Locale string
}

// LocaleCoverage represents translation completeness for a single locale.
type LocaleCoverage struct {
	Locale     string
	TotalKeys  int
	HasKeys    int
	Percentage float64
}

// DuplicateKeyIssue represents a key with conflicting values in the same locale.
type DuplicateKeyIssue struct {
	Key    string
	Locale string
	Values []DuplicateValue
}

// DuplicateValue is a value + file pair for a duplicate key.
type DuplicateValue struct {
	Value string
	File  string
}

// KeyNamingIssue represents a key that violates the naming convention.
type KeyNamingIssue struct {
	Key        string
	Expected   string
	File       string
	Locale     string
}

// AuditResult holds all findings from a scan.
type AuditResult struct {
	Summary          AuditSummary
	MissingKeys      []MissingKeyIssue
	UnusedKeys       []UnusedKeyIssue
	HardcodedStrings []HardcodedStringIssue
	DynamicKeys      []DynamicKeyWarning
	LocaleCoverage   []LocaleCoverage
	DuplicateKeys    []DuplicateKeyIssue
	KeyNamingIssues  []KeyNamingIssue
	Metadata         AuditMetadata
}

// AuditSummary provides counts of all findings.
type AuditSummary struct {
	FilesScanned         int
	I18nFilesLoaded      int
	TotalKeys            int
	Locales              []string
	MissingKeyCount      int
	UnusedKeyCount       int
	HardcodedStringCount int
	DynamicKeyCount      int
	DuplicateKeyCount    int
	KeyNamingIssueCount  int
	OverallCompleteness  float64
}

// AuditMetadata provides context about the scan.
type AuditMetadata struct {
	Tool           string
	Version        string
	Timestamp      string
	Preset         string
	RootDir        string
	I18nFileFormat string
	Duration       int64
}
