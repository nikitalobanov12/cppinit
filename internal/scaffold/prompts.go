package scaffold

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// RunPrompts runs the interactive prompts and returns the configuration
func RunPrompts() (*Config, error) {
	fmt.Println(titleStyle.Render("🚀 Create C++ Project"))
	fmt.Println(subtitleStyle.Render("Configure your new C++ project with modern CMake"))
	fmt.Println()

	config := &Config{}

	// Get current directory name as default project name
	cwd, _ := os.Getwd()
	defaultName := filepath.Base(cwd)

	// Get current user for author name
	currentUser, _ := user.Current()
	defaultAuthor := currentUser.Username

	// Page 1: Basic Info
	basicForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("The name of your C++ project").
				Value(&config.ProjectName).
				Placeholder(defaultName).
				Validate(validateProjectName),

			huh.NewInput().
				Title("Description").
				Description("A short description of your project").
				Value(&config.Description).
				Placeholder("A modern C++ project"),

			huh.NewInput().
				Title("Author name").
				Value(&config.AuthorName).
				Placeholder(defaultAuthor),

			huh.NewSelect[string]().
				Title("C++ Standard").
				Description("Which C++ standard version to use").
				Options(
					huh.NewOption("C++11", "11"),
					huh.NewOption("C++14", "14"),
					huh.NewOption("C++17 (Recommended)", "17"),
					huh.NewOption("C++20", "20"),
					huh.NewOption("C++23", "23"),
				).
				Value(&config.CppStandard),

			huh.NewSelect[string]().
				Title("Project type").
				Description("What kind of project are you building?").
				Options(
					huh.NewOption("Executable", "executable"),
					huh.NewOption("Static Library", "static"),
					huh.NewOption("Header-only Library", "header-only"),
				).
				Value(&config.ProjectType),
		).Title("Project Basics"),
	)

	if err := basicForm.Run(); err != nil {
		return nil, err
	}

	// Page 2: Dependencies & Testing
	depsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Package manager").
				Description("How do you want to manage dependencies?").
				Options(
					huh.NewOption("None (FetchContent only)", "none"),
					huh.NewOption("vcpkg", "vcpkg"),
					huh.NewOption("Conan", "conan"),
					huh.NewOption("CPM.cmake", "cpm"),
				).
				Value(&config.PackageManager),

			huh.NewSelect[string]().
				Title("Testing framework").
				Description("Include a testing framework?").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("GoogleTest", "googletest"),
					huh.NewOption("Catch2", "catch2"),
					huh.NewOption("doctest", "doctest"),
				).
				Value(&config.TestFramework),

			huh.NewConfirm().
				Title("Include benchmarks?").
				Description("Add Google Benchmark for performance testing").
				Value(&config.IncludeBenchmark),
		).Title("Dependencies & Testing"),
	)

	if err := depsForm.Run(); err != nil {
		return nil, err
	}

	// Page 3: Tooling
	toolingForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Code quality tools").
				Description("Select the tools you want to include").
				Options(
					huh.NewOption("clang-format (code formatting)", "clang-format").Selected(true),
					huh.NewOption("clang-tidy (static analysis)", "clang-tidy").Selected(true),
					huh.NewOption("Sanitizers (ASan, UBSan, TSan)", "sanitizers"),
					huh.NewOption("Code coverage (gcov/lcov)", "coverage"),
					huh.NewOption("Doxygen (documentation)", "doxygen"),
					huh.NewOption("pre-commit hooks", "pre-commit"),
				).
				Value(&[]string{}).
				Filterable(false),
		).Title("Code Quality"),
	)

	var selectedTools []string
	toolingForm = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Code quality tools").
				Description("Select the tools you want to include").
				Options(
					huh.NewOption("clang-format (code formatting)", "clang-format").Selected(true),
					huh.NewOption("clang-tidy (static analysis)", "clang-tidy").Selected(true),
					huh.NewOption("Sanitizers (ASan, UBSan, TSan)", "sanitizers"),
					huh.NewOption("Code coverage (gcov/lcov)", "coverage"),
					huh.NewOption("Doxygen (documentation)", "doxygen"),
					huh.NewOption("pre-commit hooks", "pre-commit"),
				).
				Value(&selectedTools),
		).Title("Code Quality"),
	)

	if err := toolingForm.Run(); err != nil {
		return nil, err
	}

	// Parse selected tools
	for _, tool := range selectedTools {
		switch tool {
		case "clang-format":
			config.UseClangFormat = true
		case "clang-tidy":
			config.UseClangTidy = true
		case "sanitizers":
			config.UseSanitizers = true
		case "coverage":
			config.UseCoverage = true
		case "doxygen":
			config.UseDoxygen = true
		case "pre-commit":
			config.UsePreCommit = true
		}
	}

	// Page 4: DevOps & IDE
	devopsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("License").
				Options(
					huh.NewOption("MIT", "mit"),
					huh.NewOption("Apache 2.0", "apache2"),
					huh.NewOption("GPL 3.0", "gpl3"),
					huh.NewOption("BSD 3-Clause", "bsd3"),
					huh.NewOption("None", "none"),
				).
				Value(&config.License),

			huh.NewConfirm().
				Title("Include GitHub Actions CI?").
				Description("Automated builds, tests, and linting").
				Value(&config.IncludeCI),

			huh.NewConfirm().
				Title("Include VSCode configuration?").
				Description("Settings, launch configs, and recommended extensions").
				Value(&config.IncludeVSCode),

			huh.NewConfirm().
				Title("Include Docker support?").
				Description("Dockerfile and devcontainer for VS Code").
				Value(&config.UseDocker),
		).Title("DevOps & IDE"),
	)

	if err := devopsForm.Run(); err != nil {
		return nil, err
	}

	// Set defaults
	if config.ProjectName == "" {
		config.ProjectName = defaultName
	}
	if config.Description == "" {
		config.Description = "A modern C++ project"
	}
	if config.AuthorName == "" {
		config.AuthorName = defaultAuthor
	}
	config.OutputDir = config.ProjectName

	return config, nil
}

func validateProjectName(s string) error {
	if s == "" {
		return nil // Will use placeholder
	}
	if strings.ContainsAny(s, " /\\:*?\"<>|") {
		return fmt.Errorf("project name cannot contain special characters")
	}
	if strings.HasPrefix(s, ".") || strings.HasPrefix(s, "-") {
		return fmt.Errorf("project name cannot start with . or -")
	}
	return nil
}

// PrintSuccess prints the success message with next steps
func PrintSuccess(config *Config) {
	fmt.Println()
	fmt.Println(successStyle.Render("✓ Project created successfully!"))
	fmt.Println()

	// Show what was created
	fmt.Println("Created project with:")
	fmt.Printf("  • C++%s %s\n", config.CppStandard, config.ProjectType)
	if config.TestFramework != "none" {
		fmt.Printf("  • %s testing\n", config.TestFramework)
	}
	if config.PackageManager != "none" {
		fmt.Printf("  • %s package manager\n", config.PackageManager)
	}
	if config.UseClangFormat {
		fmt.Println("  • clang-format")
	}
	if config.UseClangTidy {
		fmt.Println("  • clang-tidy")
	}
	if config.UseSanitizers {
		fmt.Println("  • Address/UB/Thread sanitizers")
	}
	if config.UseCoverage {
		fmt.Println("  • Code coverage")
	}
	if config.IncludeCI {
		fmt.Println("  • GitHub Actions CI")
	}
	if config.UseDocker {
		fmt.Println("  • Docker & devcontainer")
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println()
	fmt.Printf("  %s\n", pathStyle.Render(fmt.Sprintf("cd %s", config.ProjectName)))
	fmt.Println()
	fmt.Println("  # Configure and build")
	fmt.Println("  cmake --preset debug")
	fmt.Println("  cmake --build --preset debug")
	fmt.Println()

	if config.TestFramework != "none" {
		fmt.Println("  # Run tests")
		fmt.Println("  ctest --preset debug")
		fmt.Println()
	}

	if config.UseSanitizers {
		fmt.Println("  # Run with sanitizers")
		fmt.Println("  cmake --preset asan && cmake --build --preset asan")
		fmt.Println()
	}

	if config.UsePreCommit {
		fmt.Println("  # Setup pre-commit hooks")
		fmt.Println("  pip install pre-commit && pre-commit install")
		fmt.Println()
	}

	if config.UseDocker {
		fmt.Println("  # Or use Docker")
		fmt.Println("  docker build -t", config.ProjectName, ".")
		fmt.Println()
	}

	fmt.Println(dimStyle.Render("Happy coding! 🎉"))
}
