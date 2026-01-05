package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nikitalobanov12/cppinit/internal/templates"
)

// Generate creates the project structure based on the configuration
func Generate(config *Config) error {
	// Create base directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create directory structure
	dirs := []string{
		"src",
		"include/" + config.ProjectName,
		"cmake",
	}

	if config.TestFramework != "none" {
		dirs = append(dirs, "tests")
	}

	if config.IncludeBenchmark && config.ProjectType != "executable" {
		dirs = append(dirs, "benchmarks")
	}

	if config.IncludeVSCode {
		dirs = append(dirs, ".vscode")
	}

	if config.UseDocker {
		dirs = append(dirs, ".devcontainer")
	}

	if config.IncludeCI {
		dirs = append(dirs, ".github/workflows")
	}

	for _, dir := range dirs {
		path := filepath.Join(config.OutputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate all files
	files := make(map[string]string)

	// Core CMake files
	files["CMakeLists.txt"] = generateRootCMakeLists(config)
	files["cmake/CompilerWarnings.cmake"] = templates.CompilerWarningsCMake()

	// CMake presets
	files["CMakePresets.json"] = templates.CMakePresets(
		config.ProjectName,
		config.PackageManager,
		config.UseSanitizers,
		config.UseCoverage,
	)

	// Additional CMake modules
	if config.UseSanitizers {
		files["cmake/Sanitizers.cmake"] = templates.SanitizersCMake()
	}
	if config.UseCoverage {
		files["cmake/Coverage.cmake"] = templates.CoverageCMake()
	}
	if config.UseClangTidy {
		files["cmake/StaticAnalysis.cmake"] = templates.StaticAnalysisCMake()
	}
	if config.UseDoxygen {
		files["cmake/Doxygen.cmake"] = templates.DoxygenCMake()
	}
	if config.PackageManager == "cpm" {
		files["cmake/CPM.cmake"] = templates.CPMCMake()
	}

	// Source files
	if config.ProjectType == "executable" {
		files["src/main.cpp"] = templates.MainCpp(config.ProjectName)
	} else if config.ProjectType == "static" {
		files["src/"+config.ProjectName+".cpp"] = templates.LibraryCpp(config.ProjectName)
		files["include/"+config.ProjectName+"/"+config.ProjectName+".hpp"] = templates.LibraryHpp(config.ProjectName)
	} else if config.ProjectType == "header-only" {
		files["include/"+config.ProjectName+"/"+config.ProjectName+".hpp"] = templates.HeaderOnlyHpp(config.ProjectName)
	}

	// Test files
	if config.TestFramework != "none" {
		files["tests/CMakeLists.txt"] = templates.TestsCMakeLists(config.ProjectName, config.ProjectType, config.TestFramework)
		files["tests/test_main.cpp"] = templates.TestMainCpp(config.ProjectName, config.ProjectType, config.TestFramework)
	}

	// Benchmark files
	if config.IncludeBenchmark && config.ProjectType != "executable" {
		files["benchmarks/CMakeLists.txt"] = templates.BenchmarkCMake(config.ProjectName)
		files["benchmarks/benchmark_main.cpp"] = templates.BenchmarkMain(config.ProjectName, config.ProjectType)
	}

	// Package manager files
	switch config.PackageManager {
	case "vcpkg":
		files["vcpkg.json"] = templates.VcpkgJson(config.ProjectName, config.TestFramework)
	case "conan":
		files["conanfile.txt"] = templates.ConanfileTxt(config.TestFramework)
	}

	// Tooling configs
	if config.UseClangFormat {
		files[".clang-format"] = templates.ClangFormat()
	}
	if config.UseClangTidy {
		files[".clang-tidy"] = templates.ClangTidy(config.CppStandard)
	}
	files[".editorconfig"] = templates.EditorConfig()

	// License
	if config.License != "none" {
		year := time.Now().Format("2006")
		files["LICENSE"] = templates.License(config.License, config.AuthorName, year)
	}

	// Git files
	files[".gitignore"] = templates.GitIgnore()

	// Documentation
	files["README.md"] = generateReadme(config)

	// VSCode configuration
	if config.IncludeVSCode {
		files[".vscode/settings.json"] = templates.VSCodeSettings()
		files[".vscode/extensions.json"] = templates.VSCodeExtensions()
		files[".vscode/launch.json"] = templates.VSCodeLaunch(config.ProjectName, config.ProjectType)
		files[".vscode/tasks.json"] = templates.VSCodeTasks()
	}

	// Docker
	if config.UseDocker {
		if config.ProjectType == "executable" {
			files["Dockerfile"] = templates.Dockerfile(config.ProjectName, config.CppStandard)
		}
		files[".dockerignore"] = templates.DockerIgnore()
		files[".devcontainer/devcontainer.json"] = templates.DevContainer(config.ProjectName)
	}

	// Pre-commit
	if config.UsePreCommit {
		files[".pre-commit-config.yaml"] = templates.PreCommitConfig()
	}

	// CI
	if config.IncludeCI {
		files[".github/workflows/ci.yml"] = templates.GitHubActionsCIFull(
			config.ProjectName,
			config.PackageManager,
			config.TestFramework,
			config.UseSanitizers,
			config.UseCoverage,
		)
		files[".github/dependabot.yml"] = templates.GitHubDependabot()
	}

	// Write all files
	for filename, content := range files {
		if content == "" {
			continue
		}
		path := filepath.Join(config.OutputDir, filename)

		// Ensure parent directory exists
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filename, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// generateRootCMakeLists creates the main CMakeLists.txt with all features
func generateRootCMakeLists(config *Config) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`cmake_minimum_required(VERSION 3.21)

project(%s
    VERSION 0.1.0
    DESCRIPTION "%s"
    LANGUAGES CXX
)

# Prevent in-source builds
if(CMAKE_SOURCE_DIR STREQUAL CMAKE_BINARY_DIR)
    message(FATAL_ERROR "In-source builds are not allowed. Please use a separate build directory.")
endif()

# Set C++ standard
set(CMAKE_CXX_STANDARD %s)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for IDE/tooling support
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Include custom CMake modules
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")

`, config.ProjectName, config.Description, config.CppStandard))

	// Include CMake modules
	sb.WriteString("# Include CMake modules\n")
	sb.WriteString("include(CompilerWarnings)\n")

	if config.UseSanitizers {
		sb.WriteString("include(Sanitizers)\n")
	}
	if config.UseCoverage {
		sb.WriteString("include(Coverage)\n")
	}
	if config.UseClangTidy {
		sb.WriteString("include(StaticAnalysis)\n")
	}
	if config.UseDoxygen {
		sb.WriteString("include(Doxygen)\n")
	}
	if config.PackageManager == "cpm" {
		sb.WriteString("include(CPM)\n")
	}

	sb.WriteString("\n")

	// Add target based on project type
	switch config.ProjectType {
	case "executable":
		sb.WriteString(fmt.Sprintf(`# Main executable
add_executable(${PROJECT_NAME}
    src/main.cpp
)

target_include_directories(${PROJECT_NAME}
    PRIVATE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
)

`))
	case "static":
		sb.WriteString(fmt.Sprintf(`# Library target
add_library(${PROJECT_NAME} STATIC
    src/%s.cpp
)

# Create alias for use with FetchContent/subdirectory
add_library(${PROJECT_NAME}::${PROJECT_NAME} ALIAS ${PROJECT_NAME})

target_include_directories(${PROJECT_NAME}
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)

`, config.ProjectName))
	case "header-only":
		sb.WriteString(`# Header-only library
add_library(${PROJECT_NAME} INTERFACE)

# Create alias for use with FetchContent/subdirectory
add_library(${PROJECT_NAME}::${PROJECT_NAME} ALIAS ${PROJECT_NAME})

target_include_directories(${PROJECT_NAME}
    INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)

`)
	}

	// Apply compiler warnings
	sb.WriteString("# Apply compiler warnings\n")
	sb.WriteString("set_project_warnings(${PROJECT_NAME})\n\n")

	// Apply sanitizers
	if config.UseSanitizers {
		sb.WriteString("# Apply sanitizers (if enabled)\n")
		sb.WriteString("enable_sanitizers(${PROJECT_NAME})\n\n")
	}

	// Apply coverage
	if config.UseCoverage {
		sb.WriteString("# Apply code coverage (if enabled)\n")
		sb.WriteString("enable_coverage(${PROJECT_NAME})\n\n")
	}

	// Apply static analysis
	if config.UseClangTidy {
		sb.WriteString("# Apply static analysis (if enabled)\n")
		sb.WriteString("enable_static_analysis(${PROJECT_NAME})\n\n")
	}

	// Testing
	if config.TestFramework != "none" {
		sb.WriteString(`# Testing
option(BUILD_TESTS "Build the tests" ON)
if(BUILD_TESTS)
    enable_testing()
    add_subdirectory(tests)
endif()

`)
	}

	// Benchmarks
	if config.IncludeBenchmark && config.ProjectType != "executable" {
		sb.WriteString(`# Benchmarks
option(BUILD_BENCHMARKS "Build the benchmarks" OFF)
if(BUILD_BENCHMARKS)
    add_subdirectory(benchmarks)
endif()

`)
	}

	// Documentation
	if config.UseDoxygen {
		sb.WriteString("# Documentation\n")
		sb.WriteString("enable_docs()\n\n")
	}

	// Coverage target
	if config.UseCoverage {
		sb.WriteString("# Coverage report target\n")
		sb.WriteString("add_coverage_target()\n\n")
	}

	// Install rules for libraries
	if config.ProjectType != "executable" {
		sb.WriteString(`# Installation rules
include(GNUInstallDirs)
install(TARGETS ${PROJECT_NAME}
    EXPORT ${PROJECT_NAME}Targets
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
    INCLUDES DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}
)

install(DIRECTORY include/
    DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}
)

install(EXPORT ${PROJECT_NAME}Targets
    FILE ${PROJECT_NAME}Targets.cmake
    NAMESPACE ${PROJECT_NAME}::
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/${PROJECT_NAME}
)
`)
	}

	return sb.String()
}

// generateReadme creates a comprehensive README.md
func generateReadme(config *Config) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", config.ProjectName))
	sb.WriteString(fmt.Sprintf("%s\n\n", config.Description))

	// Badges
	if config.IncludeCI {
		sb.WriteString("![CI](https://github.com/USERNAME/" + config.ProjectName + "/workflows/CI/badge.svg)\n")
	}
	if config.License != "none" {
		sb.WriteString(fmt.Sprintf("![License](https://img.shields.io/badge/license-%s-blue.svg)\n", config.License))
	}
	sb.WriteString(fmt.Sprintf("![C++%s](https://img.shields.io/badge/C%%2B%%2B-%s-blue.svg)\n\n", config.CppStandard, config.CppStandard))

	// Features
	sb.WriteString("## Features\n\n")
	sb.WriteString(fmt.Sprintf("- Modern C++%s\n", config.CppStandard))
	sb.WriteString("- CMake 3.21+ with presets\n")
	if config.TestFramework != "none" {
		sb.WriteString(fmt.Sprintf("- %s testing framework\n", config.TestFramework))
	}
	if config.UseClangFormat {
		sb.WriteString("- clang-format for code formatting\n")
	}
	if config.UseClangTidy {
		sb.WriteString("- clang-tidy for static analysis\n")
	}
	if config.UseSanitizers {
		sb.WriteString("- Address, UB, and Thread sanitizers\n")
	}
	if config.UseCoverage {
		sb.WriteString("- Code coverage support\n")
	}
	if config.IncludeCI {
		sb.WriteString("- GitHub Actions CI/CD\n")
	}
	sb.WriteString("\n")

	// Requirements
	sb.WriteString("## Requirements\n\n")
	sb.WriteString("- CMake 3.21 or higher\n")
	sb.WriteString(fmt.Sprintf("- C++%s compatible compiler (GCC 10+, Clang 12+, MSVC 2019+)\n", config.CppStandard))
	if config.PackageManager == "vcpkg" {
		sb.WriteString("- vcpkg (optional, for dependency management)\n")
	} else if config.PackageManager == "conan" {
		sb.WriteString("- Conan (optional, for dependency management)\n")
	}
	sb.WriteString("\n")

	// Building
	sb.WriteString("## Building\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("# Configure (debug build)\n")
	sb.WriteString("cmake --preset debug\n\n")
	sb.WriteString("# Build\n")
	sb.WriteString("cmake --build --preset debug\n\n")
	sb.WriteString("# Or for release\n")
	sb.WriteString("cmake --preset release\n")
	sb.WriteString("cmake --build --preset release\n")
	sb.WriteString("```\n\n")

	// Testing
	if config.TestFramework != "none" {
		sb.WriteString("## Testing\n\n")
		sb.WriteString("```bash\n")
		sb.WriteString("# Run tests\n")
		sb.WriteString("ctest --preset debug\n\n")
		sb.WriteString("# Or with verbose output\n")
		sb.WriteString("ctest --preset debug --output-on-failure\n")
		sb.WriteString("```\n\n")
	}

	// Sanitizers
	if config.UseSanitizers {
		sb.WriteString("## Sanitizers\n\n")
		sb.WriteString("```bash\n")
		sb.WriteString("# AddressSanitizer\n")
		sb.WriteString("cmake --preset asan\n")
		sb.WriteString("cmake --build --preset asan\n\n")
		sb.WriteString("# UndefinedBehaviorSanitizer\n")
		sb.WriteString("cmake --preset ubsan\n")
		sb.WriteString("cmake --build --preset ubsan\n\n")
		sb.WriteString("# ThreadSanitizer\n")
		sb.WriteString("cmake --preset tsan\n")
		sb.WriteString("cmake --build --preset tsan\n")
		sb.WriteString("```\n\n")
	}

	// Coverage
	if config.UseCoverage {
		sb.WriteString("## Code Coverage\n\n")
		sb.WriteString("```bash\n")
		sb.WriteString("cmake --preset coverage\n")
		sb.WriteString("cmake --build --preset coverage\n")
		sb.WriteString("ctest --preset debug\n")
		sb.WriteString("cmake --build --preset coverage --target coverage\n")
		sb.WriteString("# Open build/coverage/coverage_report/index.html\n")
		sb.WriteString("```\n\n")
	}

	// Docker
	if config.UseDocker {
		sb.WriteString("## Docker\n\n")
		sb.WriteString("```bash\n")
		sb.WriteString("# Build image\n")
		sb.WriteString(fmt.Sprintf("docker build -t %s .\n\n", config.ProjectName))
		sb.WriteString("# Run container\n")
		sb.WriteString(fmt.Sprintf("docker run --rm %s\n", config.ProjectName))
		sb.WriteString("```\n\n")
		sb.WriteString("### VS Code Dev Container\n\n")
		sb.WriteString("Open the project in VS Code and click \"Reopen in Container\" when prompted.\n\n")
	}

	// Project structure
	sb.WriteString("## Project Structure\n\n")
	sb.WriteString("```\n")
	sb.WriteString(config.ProjectName + "/\n")
	sb.WriteString("├── CMakeLists.txt          # Main CMake configuration\n")
	sb.WriteString("├── CMakePresets.json       # CMake presets for easy building\n")
	sb.WriteString("├── cmake/                  # CMake modules\n")
	sb.WriteString("│   ├── CompilerWarnings.cmake\n")
	if config.UseSanitizers {
		sb.WriteString("│   ├── Sanitizers.cmake\n")
	}
	if config.UseCoverage {
		sb.WriteString("│   ├── Coverage.cmake\n")
	}
	sb.WriteString("├── include/                # Public headers\n")
	sb.WriteString(fmt.Sprintf("│   └── %s/\n", config.ProjectName))
	sb.WriteString("├── src/                    # Source files\n")
	if config.TestFramework != "none" {
		sb.WriteString("├── tests/                  # Test files\n")
	}
	if config.IncludeVSCode {
		sb.WriteString("├── .vscode/                # VS Code configuration\n")
	}
	if config.UseDocker {
		sb.WriteString("├── .devcontainer/          # Dev container configuration\n")
		sb.WriteString("├── Dockerfile\n")
	}
	sb.WriteString("└── README.md\n")
	sb.WriteString("```\n\n")

	// License
	if config.License != "none" {
		sb.WriteString("## License\n\n")
		licenseName := map[string]string{
			"mit":     "MIT",
			"apache2": "Apache 2.0",
			"gpl3":    "GPL 3.0",
			"bsd3":    "BSD 3-Clause",
		}[config.License]
		sb.WriteString(fmt.Sprintf("This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.\n", licenseName))
	}

	return sb.String()
}
