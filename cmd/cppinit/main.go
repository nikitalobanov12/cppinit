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
	description := flag.String("desc", "A modern C++ project", "Project description")
	author := flag.String("author", "", "Author name")
	cppStd := flag.String("std", "17", "C++ standard (11, 14, 17, 20, 23)")
	projectType := flag.String("type", "executable", "Project type (executable, static, header-only)")
	testFw := flag.String("tests", "none", "Test framework (none, googletest, catch2, doctest)")
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
		config = &scaffold.Config{
			ProjectName:      *name,
			Description:      *description,
			AuthorName:       *author,
			CppStandard:      *cppStd,
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
	fmt.Println(`cppinit - Create C++ projects with modern CMake

Usage:
  cppinit                    Run interactive project wizard
  cppinit -name <name>       Create project with specified options (non-interactive)

Project Options:
  -name string         Project name (required for non-interactive mode)
  -desc string         Project description (default "A modern C++ project")
  -author string       Author name for license
  -std string          C++ standard: 11, 14, 17, 20, 23 (default "17")
  -type string         Project type: executable, static, header-only (default "executable")
  -license string      License: none, mit, apache2, gpl3, bsd3 (default "mit")

Dependencies:
  -tests string        Test framework: none, googletest, catch2, doctest (default "none")
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

  # Quick executable with defaults
  cppinit -name myapp

  # Full-featured library
  cppinit -name mylib -type static -std 20 -tests googletest -full

  # Minimal header-only library
  cppinit -name myheader -type header-only -minimal

  # Executable with specific features
  cppinit -name myapp -tests catch2 -sanitizers -ci -vscode`)
}
