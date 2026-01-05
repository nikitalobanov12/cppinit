package templates

import "fmt"

// VSCodeSettings generates .vscode/settings.json
func VSCodeSettings() string {
	return `{
    "cmake.configureOnOpen": true,
    "cmake.buildDirectory": "${workspaceFolder}/build/debug",
    "cmake.configureSettings": {
        "CMAKE_EXPORT_COMPILE_COMMANDS": "ON"
    },
    "C_Cpp.default.configurationProvider": "ms-vscode.cmake-tools",
    "C_Cpp.default.compileCommands": "${workspaceFolder}/build/debug/compile_commands.json",
    "C_Cpp.clang_format_style": "file",
    "C_Cpp.codeAnalysis.clangTidy.enabled": true,
    "C_Cpp.codeAnalysis.clangTidy.useBuildPath": true,
    "editor.formatOnSave": true,
    "editor.tabSize": 4,
    "files.insertFinalNewline": true,
    "files.trimTrailingWhitespace": true,
    "files.associations": {
        "*.hpp": "cpp",
        "*.h": "cpp",
        "*.cpp": "cpp",
        "*.tpp": "cpp"
    },
    "[cpp]": {
        "editor.defaultFormatter": "ms-vscode.cpptools"
    }
}
`
}

// VSCodeExtensions generates .vscode/extensions.json
func VSCodeExtensions() string {
	return `{
    "recommendations": [
        "ms-vscode.cpptools",
        "ms-vscode.cmake-tools",
        "ms-vscode.cpptools-extension-pack",
        "twxs.cmake",
        "xaver.clang-format",
        "cschlosser.doxdocgen",
        "jeff-hykin.better-cpp-syntax",
        "vadimcn.vscode-lldb"
    ]
}
`
}

// VSCodeLaunch generates .vscode/launch.json
func VSCodeLaunch(projectName, projectType string) string {
	if projectType != "executable" {
		return `{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Run Tests (GDB)",
            "type": "cppdbg",
            "request": "launch",
            "program": "${workspaceFolder}/build/debug/tests/tests",
            "args": [],
            "stopAtEntry": false,
            "cwd": "${workspaceFolder}",
            "environment": [],
            "externalConsole": false,
            "MIMode": "gdb",
            "setupCommands": [
                {
                    "description": "Enable pretty-printing for gdb",
                    "text": "-enable-pretty-printing",
                    "ignoreFailures": true
                }
            ],
            "preLaunchTask": "CMake: build"
        },
        {
            "name": "Run Tests (LLDB)",
            "type": "lldb",
            "request": "launch",
            "program": "${workspaceFolder}/build/debug/tests/tests",
            "args": [],
            "cwd": "${workspaceFolder}",
            "preLaunchTask": "CMake: build"
        }
    ]
}
`
	}

	return fmt.Sprintf(`{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug (GDB)",
            "type": "cppdbg",
            "request": "launch",
            "program": "${workspaceFolder}/build/debug/%s",
            "args": [],
            "stopAtEntry": false,
            "cwd": "${workspaceFolder}",
            "environment": [],
            "externalConsole": false,
            "MIMode": "gdb",
            "setupCommands": [
                {
                    "description": "Enable pretty-printing for gdb",
                    "text": "-enable-pretty-printing",
                    "ignoreFailures": true
                },
                {
                    "description": "Set Disassembly Flavor to Intel",
                    "text": "-gdb-set disassembly-flavor intel",
                    "ignoreFailures": true
                }
            ],
            "preLaunchTask": "CMake: build"
        },
        {
            "name": "Debug (LLDB)",
            "type": "lldb",
            "request": "launch",
            "program": "${workspaceFolder}/build/debug/%s",
            "args": [],
            "cwd": "${workspaceFolder}",
            "preLaunchTask": "CMake: build"
        },
        {
            "name": "Run Tests (GDB)",
            "type": "cppdbg",
            "request": "launch",
            "program": "${workspaceFolder}/build/debug/tests/tests",
            "args": [],
            "stopAtEntry": false,
            "cwd": "${workspaceFolder}",
            "environment": [],
            "externalConsole": false,
            "MIMode": "gdb",
            "preLaunchTask": "CMake: build"
        }
    ]
}
`, projectName, projectName)
}

// VSCodeTasks generates .vscode/tasks.json
func VSCodeTasks() string {
	return `{
    "version": "2.0.0",
    "tasks": [
        {
            "type": "cmake",
            "label": "CMake: configure",
            "command": "configure",
            "preset": "${command:cmake.activeConfigurePresetName}",
            "problemMatcher": []
        },
        {
            "type": "cmake",
            "label": "CMake: build",
            "command": "build",
            "preset": "${command:cmake.activeBuildPresetName}",
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "problemMatcher": "$gcc"
        },
        {
            "label": "Run clang-format",
            "type": "shell",
            "command": "find src include tests -name '*.cpp' -o -name '*.hpp' | xargs clang-format -i",
            "problemMatcher": []
        },
        {
            "label": "Run clang-tidy",
            "type": "shell",
            "command": "run-clang-tidy -p build/debug",
            "problemMatcher": []
        },
        {
            "label": "Run tests",
            "type": "shell",
            "command": "ctest --preset debug --output-on-failure",
            "group": {
                "kind": "test",
                "isDefault": true
            },
            "problemMatcher": []
        },
        {
            "label": "Clean build",
            "type": "shell",
            "command": "rm -rf build",
            "problemMatcher": []
        }
    ]
}
`
}

// Dockerfile generates a multi-stage Dockerfile
func Dockerfile(projectName, cppStandard string) string {
	return fmt.Sprintf(`# syntax=docker/dockerfile:1

# Build stage
FROM gcc:13 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    cmake \
    ninja-build \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy source files
COPY . .

# Build the project
RUN cmake -B build -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_CXX_STANDARD=%s \
    -DBUILD_TESTS=OFF \
    && cmake --build build

# Runtime stage
FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built executable
COPY --from=builder /app/build/%s /app/%s

# Run as non-root user
RUN useradd -m -s /bin/bash appuser
USER appuser

ENTRYPOINT ["/app/%s"]
`, cppStandard, projectName, projectName, projectName)
}

// DockerIgnore generates .dockerignore
func DockerIgnore() string {
	return `# Build artifacts
build/
cmake-build-*/
out/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Git
.git/
.gitignore

# Documentation
docs/
*.md

# Testing
tests/
coverage/

# Package managers
vcpkg_installed/
conan/
`
}

// DevContainer generates .devcontainer/devcontainer.json
func DevContainer(projectName string) string {
	return fmt.Sprintf(`{
    "name": "%s Development",
    "image": "mcr.microsoft.com/devcontainers/cpp:1-debian-12",
    "features": {
        "ghcr.io/devcontainers/features/cmake:1": {
            "version": "latest"
        },
        "ghcr.io/devcontainers/features/ninja:1": {}
    },
    "customizations": {
        "vscode": {
            "settings": {
                "cmake.configureOnOpen": true,
                "C_Cpp.default.configurationProvider": "ms-vscode.cmake-tools"
            },
            "extensions": [
                "ms-vscode.cpptools",
                "ms-vscode.cmake-tools",
                "ms-vscode.cpptools-extension-pack",
                "twxs.cmake",
                "xaver.clang-format"
            ]
        }
    },
    "postCreateCommand": "cmake --preset debug",
    "remoteUser": "vscode"
}
`, projectName)
}

// PreCommitConfig generates .pre-commit-config.yaml
func PreCommitConfig() string {
	return `# Pre-commit hooks for C++ projects
# Install: pip install pre-commit && pre-commit install

repos:
  # General hooks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files
        args: ['--maxkb=1000']
      - id: check-merge-conflict
      - id: mixed-line-ending
        args: ['--fix=lf']

  # CMake formatting
  - repo: https://github.com/cheshirekow/cmake-format-precommit
    rev: v0.6.13
    hooks:
      - id: cmake-format
        args: ['--in-place']
      - id: cmake-lint

  # C++ formatting with clang-format
  - repo: https://github.com/pre-commit/mirrors-clang-format
    rev: v17.0.6
    hooks:
      - id: clang-format
        types_or: [c++, c]
        args: ['-style=file', '-i']

  # Markdown linting
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.38.0
    hooks:
      - id: markdownlint
        args: ['--fix']

  # YAML formatting
  - repo: https://github.com/macisamuele/language-formatters-pre-commit-hooks
    rev: v2.12.0
    hooks:
      - id: pretty-format-yaml
        args: ['--autofix', '--indent', '2']

# Local hooks for project-specific checks
  - repo: local
    hooks:
      - id: cmake-build-check
        name: CMake Build Check
        entry: bash -c 'cmake --preset debug && cmake --build --preset debug'
        language: system
        pass_filenames: false
        stages: [push]
`
}

// GitHubActionsCIFull generates a comprehensive CI workflow
func GitHubActionsCIFull(projectName, packageManager, testFramework string, useSanitizers, useCoverage bool) string {
	testJob := ""
	if testFramework != "none" {
		testJob = `
  test:
    needs: build
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        build_type: [Debug, Release]

    steps:
      - uses: actions/checkout@v4

      - name: Download build artifacts
        uses: actions/download-artifact@v4
        with:
          name: build-${{ matrix.os }}-${{ matrix.build_type }}
          path: build

      - name: Run tests
        run: ctest --test-dir build --output-on-failure
`
	}

	sanitizerJob := ""
	if useSanitizers {
		sanitizerJob = `
  sanitizers:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        sanitizer: [asan, ubsan, tsan]

    steps:
      - uses: actions/checkout@v4

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y ninja-build

      - name: Configure with ${{ matrix.sanitizer }}
        run: cmake --preset ${{ matrix.sanitizer }}

      - name: Build
        run: cmake --build --preset ${{ matrix.sanitizer }}

      - name: Test
        run: ctest --preset debug --output-on-failure
        env:
          ASAN_OPTIONS: detect_leaks=1:strict_string_checks=1
          UBSAN_OPTIONS: print_stacktrace=1
          TSAN_OPTIONS: second_deadlock_stack=1
`
	}

	coverageJob := ""
	if useCoverage {
		coverageJob = `
  coverage:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y ninja-build lcov

      - name: Configure with coverage
        run: cmake --preset coverage

      - name: Build
        run: cmake --build --preset coverage

      - name: Run tests
        run: ctest --preset debug --output-on-failure

      - name: Generate coverage report
        run: |
          lcov --directory . --capture --output-file coverage.info
          lcov --remove coverage.info '/usr/*' '*/tests/*' '*/build/*' --output-file coverage.info

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          files: coverage.info
          fail_ci_if_error: true
`
	}

	vcpkgSetup := ""
	if packageManager == "vcpkg" {
		vcpkgSetup = `
      - name: Setup vcpkg
        uses: lukka/run-vcpkg@v11
        with:
          vcpkgGitCommitId: 'a34c873a9717a888f58dc05268dea15592c2f0ff'`
	}

	return fmt.Sprintf(`name: CI

on:
  push:
    branches: [main, master, develop]
  pull_request:
    branches: [main, master]

env:
  CMAKE_VERSION: '3.28'
  NINJA_VERSION: '1.11.1'

jobs:
  build:
    runs-on: ${{ matrix.os }}

    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        build_type: [Debug, Release]
        compiler:
          - { cc: gcc, cxx: g++ }
          - { cc: clang, cxx: clang++ }
        exclude:
          - os: windows-latest
            compiler: { cc: clang, cxx: clang++ }

    steps:
      - uses: actions/checkout@v4
%s
      - name: Install Ninja
        uses: seanmiddleditch/gha-setup-ninja@v4

      - name: Configure CMake
        run: >
          cmake -B build -G Ninja
          -DCMAKE_BUILD_TYPE=${{ matrix.build_type }}
          -DCMAKE_C_COMPILER=${{ matrix.compiler.cc }}
          -DCMAKE_CXX_COMPILER=${{ matrix.compiler.cxx }}

      - name: Build
        run: cmake --build build --config ${{ matrix.build_type }}

      - name: Upload build artifacts
        uses: actions/upload-artifact@v4
        with:
          name: build-${{ matrix.os }}-${{ matrix.build_type }}
          path: build
%s%s%s
  lint:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Install clang-format
        run: sudo apt-get install -y clang-format

      - name: Check formatting
        run: |
          find src include tests -name '*.cpp' -o -name '*.hpp' -o -name '*.h' | \
            xargs clang-format --dry-run --Werror

      - name: Install cmake-format
        run: pip install cmake-format

      - name: Check CMake formatting
        run: cmake-format --check CMakeLists.txt cmake/*.cmake
`, vcpkgSetup, testJob, sanitizerJob, coverageJob)
}

// GitHubDependabot generates .github/dependabot.yml
func GitHubDependabot() string {
	return `version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
`
}

// CPMCMake generates cmake/CPM.cmake bootstrap
func CPMCMake() string {
	return `# CPM.cmake - Package Manager
# https://github.com/cpm-cmake/CPM.cmake

set(CPM_DOWNLOAD_VERSION 0.38.7)

if(CPM_SOURCE_CACHE)
    set(CPM_DOWNLOAD_LOCATION "${CPM_SOURCE_CACHE}/cpm/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
elseif(DEFINED ENV{CPM_SOURCE_CACHE})
    set(CPM_DOWNLOAD_LOCATION "$ENV{CPM_SOURCE_CACHE}/cpm/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
else()
    set(CPM_DOWNLOAD_LOCATION "${CMAKE_BINARY_DIR}/cmake/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
endif()

get_filename_component(CPM_DOWNLOAD_LOCATION ${CPM_DOWNLOAD_LOCATION} ABSOLUTE)

function(download_cpm)
    message(STATUS "Downloading CPM.cmake to ${CPM_DOWNLOAD_LOCATION}")
    file(DOWNLOAD
        https://github.com/cpm-cmake/CPM.cmake/releases/download/v${CPM_DOWNLOAD_VERSION}/CPM.cmake
        ${CPM_DOWNLOAD_LOCATION}
    )
endfunction()

if(NOT (EXISTS ${CPM_DOWNLOAD_LOCATION}))
    download_cpm()
endif()

include(${CPM_DOWNLOAD_LOCATION})
`
}

// BenchmarkCMake generates benchmark setup
func BenchmarkCMake(projectName string) string {
	return fmt.Sprintf(`include(FetchContent)

FetchContent_Declare(
    googlebenchmark
    GIT_REPOSITORY https://github.com/google/benchmark.git
    GIT_TAG v1.8.3
)

set(BENCHMARK_ENABLE_TESTING OFF CACHE BOOL "" FORCE)
set(BENCHMARK_ENABLE_GTEST_TESTS OFF CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googlebenchmark)

add_executable(benchmarks
    benchmark_main.cpp
)

target_link_libraries(benchmarks
    PRIVATE
        benchmark::benchmark
        %s
)

target_include_directories(benchmarks
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
)
`, projectName)
}

// BenchmarkMain generates benchmarks/benchmark_main.cpp
func BenchmarkMain(projectName, projectType string) string {
	if projectType == "executable" {
		return `#include <benchmark/benchmark.h>

static void BM_Example(benchmark::State& state) {
    for (auto _ : state) {
        // Benchmark code here
        benchmark::DoNotOptimize(1 + 1);
    }
}
BENCHMARK(BM_Example);

BENCHMARK_MAIN();
`
	}

	return fmt.Sprintf(`#include <benchmark/benchmark.h>
#include "%s/%s.hpp"

static void BM_Add(benchmark::State& state) {
    for (auto _ : state) {
        benchmark::DoNotOptimize(%s::add(state.range(0), state.range(0)));
    }
}
BENCHMARK(BM_Add)->Range(8, 8 << 10);

BENCHMARK_MAIN();
`, projectName, projectName, projectName)
}
