package scaffold

// Config holds all the project configuration options
type Config struct {
	ProjectName    string
	Description    string
	Language       string // "c" or "c++"
	Standard       string // C: "89", "99", "11", "17", "23" | C++: "11", "14", "17", "20", "23"
	ProjectType    string // "executable", "static", "header-only"
	TestFramework  string // "none", "googletest", "catch2", "doctest" (C++ only), "unity" (C only)
	PackageManager string // "none", "vcpkg", "conan", "cpm"
	License        string // "none", "mit", "apache2", "gpl3", "bsd3"

	// Feature flags
	UseClangFormat   bool
	UseClangTidy     bool
	UseSanitizers    bool
	UseCoverage      bool
	UseDoxygen       bool
	UseDocker        bool
	UsePreCommit     bool
	IncludeCI        bool
	IncludeVSCode    bool
	IncludeBenchmark bool

	// Metadata
	AuthorName  string
	AuthorEmail string
	GitRepo     string

	OutputDir string
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Language:       "c++",
		Standard:       "17",
		ProjectType:    "executable",
		TestFramework:  "none",
		PackageManager: "none",
		License:        "mit",
		UseClangFormat: true,
		UseClangTidy:   true,
	}
}

// IsC returns true if the project is a C project
func (c *Config) IsC() bool {
	return c.Language == "c"
}

// IsCpp returns true if the project is a C++ project
func (c *Config) IsCpp() bool {
	return c.Language == "c++" || c.Language == ""
}
