package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nikitalobanov12/cppinit/internal/scaffold"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse flags
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("help", false, "Show help")

	// Quick create flags (non-interactive)
	name := flag.String("name", "", "Project name (enables non-interactive mode)")
	description := flag.String("desc", "", "Project description")
	author := flag.String("author", "", "Author name")
	language := flag.String("lang", "c++", "Language (c, c++)")
	std := flag.String("std", "", "Standard (C: 89, 99, 11, 17, 23 | C++: 11, 14, 17, 20, 23)")
	projectType := flag.String("type", "executable", "Project type (executable, static, header-only)")
	testFw := flag.String("tests", "none", "Test framework (none, googletest, catch2, doctest for C++; none, unity for C)")
	pkgMgr := flag.String("pkg", "none", "Package manager (none, vcpkg, conan, cpm)")
	license := flag.String("license", "mit", "License (none, mit, apache2, gpl3, bsd3)")

	// Feature flags
	clangFormat := flag.Bool("clang-format", true, "Include clang-format configuration")
	clangTidy := flag.Bool("clang-tidy", true, "Include clang-tidy configuration")
	sanitizers := flag.Bool("sanitizers", false, "Include sanitizer support")
	coverage := flag.Bool("coverage", false, "Include code coverage support")
	doxygen := flag.Bool("doxygen", false, "Include Doxygen documentation")
	docker := flag.Bool("docker", false, "Include Docker/devcontainer support")
	precommit := flag.Bool("precommit", false, "Include pre-commit hooks")
	ci := flag.Bool("ci", false, "Include GitHub Actions CI")
	vscode := flag.Bool("vscode", false, "Include VSCode configuration")
	benchmark := flag.Bool("benchmark", false, "Include Google Benchmark")

	// Preset flags
	full := flag.Bool("full", false, "Include all features (same as --all)")
	minimal := flag.Bool("minimal", false, "Minimal project (no extra tools)")

	flag.Parse()

	if *showVersion {
		fmt.Printf("cppinit %s\n", version)
		return nil
	}

	if *showHelp {
		printHelp()
		return nil
	}

	var config *scaffold.Config
	var err error

	// Non-interactive mode if name is provided
	if *name != "" {
		// Set default standard based on language
		standard := *std
		if standard == "" {
			if *language == "c" {
				standard = "11" // C11 default
			} else {
				standard = "17" // C++17 default
			}
		}

		// Set default description based on language
		desc := *description
		if desc == "" {
			if *language == "c" {
				desc = "A modern C project"
			} else {
				desc = "A modern C++ project"
			}
		}

		config = &scaffold.Config{
			ProjectName:      *name,
			Description:      desc,
			AuthorName:       *author,
			Language:         *language,
			Standard:         standard,
			ProjectType:      *projectType,
			TestFramework:    *testFw,
			PackageManager:   *pkgMgr,
			License:          *license,
			UseClangFormat:   *clangFormat,
			UseClangTidy:     *clangTidy,
			UseSanitizers:    *sanitizers,
			UseCoverage:      *coverage,
			UseDoxygen:       *doxygen,
			UseDocker:        *docker,
			UsePreCommit:     *precommit,
			IncludeCI:        *ci,
			IncludeVSCode:    *vscode,
			IncludeBenchmark: *benchmark,
			OutputDir:        *name,
		}

		// Apply presets
		if *full {
			config.UseClangFormat = true
			config.UseClangTidy = true
			config.UseSanitizers = true
			config.UseCoverage = true
			config.UseDoxygen = true
			config.UseDocker = true
			config.UsePreCommit = true
			config.IncludeCI = true
			config.IncludeVSCode = true
			if config.TestFramework == "none" {
				config.TestFramework = "googletest"
			}
		}

		if *minimal {
			config.UseClangFormat = false
			config.UseClangTidy = false
			config.UseSanitizers = false
			config.UseCoverage = false
			config.UseDoxygen = false
			config.UseDocker = false
			config.UsePreCommit = false
			config.IncludeCI = false
			config.IncludeVSCode = false
		}
	} else {
		// Interactive mode
		config, err = scaffold.RunPrompts()
		if err != nil {
			return err
		}
	}

	if err := scaffold.Generate(config); err != nil {
		return err
	}

	scaffold.PrintSuccess(config)
	return nil
}

func printHelp() {
	fmt.Println(`cppinit - Create C/C++ projects with modern CMake

Usage:
  cppinit                    Run interactive project wizard
  cppinit -name <name>       Create project with specified options (non-interactive)

Project Options:
  -name string         Project name (required for non-interactive mode)
  -desc string         Project description
  -author string       Author name for license
  -lang string         Language: c, c++ (default "c++")
  -std string          Standard (C: 89, 99, 11, 17, 23 | C++: 11, 14, 17, 20, 23)
                       Defaults to C11 for C, C++17 for C++
  -type string         Project type: executable, static, header-only (default "executable")
  -license string      License: none, mit, apache2, gpl3, bsd3 (default "mit")

Dependencies:
  -tests string        Test framework:
                         C++: none, googletest, catch2, doctest (default "none")
                         C: none, unity (default "none")
  -pkg string          Package manager: none, vcpkg, conan, cpm (default "none")
  -benchmark           Include Google Benchmark for performance testing

Code Quality:
  -clang-format        Include clang-format config (default true)
  -clang-tidy          Include clang-tidy config (default true)
  -sanitizers          Include Address/UB/Thread sanitizers
  -coverage            Include code coverage support

DevOps & Tooling:
  -ci                  Include GitHub Actions CI
  -vscode              Include VSCode configuration
  -docker              Include Dockerfile and devcontainer
  -precommit           Include pre-commit hooks
  -doxygen             Include Doxygen documentation setup

Presets:
  -full                Enable all features (tests, sanitizers, coverage, CI, etc.)
  -minimal             Minimal project with no extra tooling

Other:
  -version             Show version
  -help                Show this help message

Examples:
  # Interactive wizard
  cppinit

  # Quick C++ executable with defaults
  cppinit -name myapp

  # C project
  cppinit -name myapp -lang c -std 11

  # Full-featured C++ library
  cppinit -name mylib -type static -std 20 -tests googletest -full

  # Minimal header-only library
  cppinit -name myheader -type header-only -minimal

  # C library with Unity tests
  cppinit -name myclib -lang c -type static -tests unity

  # Executable with specific features
  cppinit -name myapp -tests catch2 -sanitizers -ci -vscode`)
}
