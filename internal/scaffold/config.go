package scaffold

// Config holds all the project configuration options
type Config struct {
	ProjectName    string
	Description    string
	CppStandard    string // "11", "14", "17", "20", "23"
	ProjectType    string // "executable", "static", "header-only"
	TestFramework  string // "none", "googletest", "catch2", "doctest"
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
		CppStandard:    "17",
		ProjectType:    "executable",
		TestFramework:  "none",
		PackageManager: "none",
		License:        "mit",
		UseClangFormat: true,
		UseClangTidy:   true,
	}
}
