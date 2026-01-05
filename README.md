# cppinit

Create modern C++ projects with best practices, tooling, and CMake presets.

[![CI](https://github.com/nikitalobanov12/cppinit/workflows/CI/badge.svg)](https://github.com/nikitalobanov12/cppinit/actions)
[![Release](https://img.shields.io/github/v/release/nikitalobanov12/cppinit)](https://github.com/nikitalobanov12/cppinit/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Interactive wizard** - create-next-app style experience
- **Modern CMake** - CMake 3.21+ with presets
- **Multiple project types** - Executable, static library, header-only library
- **C++ standards** - C++11, 14, 17, 20, 23
- **Testing frameworks** - GoogleTest, Catch2, doctest
- **Package managers** - vcpkg, Conan, CPM.cmake
- **Code quality** - clang-format, clang-tidy, pre-commit hooks
- **Sanitizers** - AddressSanitizer, UBSan, ThreadSanitizer, MemorySanitizer
- **Code coverage** - gcov/lcov support
- **Documentation** - Doxygen integration
- **IDE support** - VSCode configuration (settings, launch, tasks, extensions)
- **CI/CD** - GitHub Actions with matrix builds
- **Containers** - Dockerfile and VS Code devcontainer
- **Licenses** - MIT, Apache 2.0, GPL 3.0, BSD 3-Clause

## Installation

### Homebrew (macOS/Linux)

```bash
brew install nikitalobanov12/tap/cppinit
```

### Go Install

```bash
go install github.com/nikitalobanov12/cppinit/cmd/cppinit@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/nikitalobanov12/cppinit/releases).

### From Source

```bash
git clone https://github.com/nikitalobanov12/cppinit.git
cd cppinit
go build -o cppinit ./cmd/cppinit
sudo mv cppinit /usr/local/bin/
```

## Usage

### Interactive Mode

Run the wizard to configure your project interactively:

```bash
cppinit
```

### Non-Interactive Mode

Create a project with command-line flags:

```bash
# Basic executable
cppinit -name myapp

# C++20 library with tests and full tooling
cppinit -name mylib -type static -std 20 -tests googletest -full

# Minimal header-only library
cppinit -name myheader -type header-only -minimal

# Executable with specific features
cppinit -name myapp -tests catch2 -sanitizers -ci -vscode
```

### CLI Options

```
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
  -full                Enable all features
  -minimal             Minimal project with no extra tooling
```

## Generated Project Structure

```
myproject/
├── CMakeLists.txt              # Main CMake configuration
├── CMakePresets.json           # Build presets (debug, release, asan, etc.)
├── cmake/
│   ├── CompilerWarnings.cmake  # Compiler warning flags
│   ├── Sanitizers.cmake        # Sanitizer configuration
│   ├── Coverage.cmake          # Code coverage setup
│   ├── StaticAnalysis.cmake    # clang-tidy integration
│   └── Doxygen.cmake           # Documentation generation
├── src/
│   └── main.cpp                # (or library source)
├── include/
│   └── myproject/              # Public headers
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
├── .vscode/
│   ├── settings.json
│   ├── launch.json
│   ├── tasks.json
│   └── extensions.json
├── .devcontainer/
│   └── devcontainer.json
├── .github/
│   ├── workflows/ci.yml
│   └── dependabot.yml
├── .clang-format
├── .clang-tidy
├── .editorconfig
├── .pre-commit-config.yaml
├── .gitignore
├── Dockerfile
├── LICENSE
└── README.md
```

## Building Generated Projects

```bash
cd myproject

# Configure and build (debug)
cmake --preset debug
cmake --build --preset debug

# Release build
cmake --preset release
cmake --build --preset release

# Run tests
ctest --preset debug

# Build with AddressSanitizer
cmake --preset asan
cmake --build --preset asan

# Build with code coverage
cmake --preset coverage
cmake --build --preset coverage
ctest --preset debug
cmake --build --preset coverage --target coverage
```

## Development

### Prerequisites

- Go 1.22+
- [GoReleaser](https://goreleaser.com/) (for releases)

### Building

```bash
go build -o cppinit ./cmd/cppinit
```

### Testing

```bash
go test ./...
```

### Releasing

Releases are automated via GitHub Actions when a tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
