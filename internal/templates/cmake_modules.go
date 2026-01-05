package templates

import "fmt"

// CMakePresets generates a comprehensive CMakePresets.json
func CMakePresets(projectName, packageManager string, useSanitizers, useCoverage bool) string {
	toolchainFile := ""
	if packageManager == "vcpkg" {
		toolchainFile = `
            "toolchainFile": "$env{VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake",`
	} else if packageManager == "conan" {
		toolchainFile = `
            "toolchainFile": "${sourceDir}/build/conan_toolchain.cmake",`
	}

	sanitizerPresets := ""
	if useSanitizers {
		sanitizerPresets = `,
        {
            "name": "asan",
            "displayName": "AddressSanitizer",
            "inherits": "debug",
            "cacheVariables": {
                "ENABLE_SANITIZER_ADDRESS": "ON"
            }
        },
        {
            "name": "ubsan",
            "displayName": "UndefinedBehaviorSanitizer",
            "inherits": "debug",
            "cacheVariables": {
                "ENABLE_SANITIZER_UNDEFINED": "ON"
            }
        },
        {
            "name": "tsan",
            "displayName": "ThreadSanitizer",
            "inherits": "debug",
            "cacheVariables": {
                "ENABLE_SANITIZER_THREAD": "ON"
            }
        },
        {
            "name": "msan",
            "displayName": "MemorySanitizer (Clang only)",
            "inherits": "debug",
            "cacheVariables": {
                "ENABLE_SANITIZER_MEMORY": "ON"
            }
        }`
	}

	coveragePreset := ""
	coverageBuildPreset := ""
	if useCoverage {
		coveragePreset = `,
        {
            "name": "coverage",
            "displayName": "Code Coverage",
            "inherits": "debug",
            "cacheVariables": {
                "ENABLE_COVERAGE": "ON"
            }
        }`
		coverageBuildPreset = `,
        {
            "name": "coverage",
            "configurePreset": "coverage"
        }`
	}

	sanitizerBuildPresets := ""
	if useSanitizers {
		sanitizerBuildPresets = `,
        {
            "name": "asan",
            "configurePreset": "asan"
        },
        {
            "name": "ubsan",
            "configurePreset": "ubsan"
        },
        {
            "name": "tsan",
            "configurePreset": "tsan"
        }`
	}

	return fmt.Sprintf(`{
    "version": 6,
    "cmakeMinimumRequired": {
        "major": 3,
        "minor": 21,
        "patch": 0
    },
    "configurePresets": [
        {
            "name": "base",
            "hidden": true,
            "binaryDir": "${sourceDir}/build/${presetName}",
            "installDir": "${sourceDir}/install/${presetName}",%s
            "cacheVariables": {
                "CMAKE_EXPORT_COMPILE_COMMANDS": "ON"
            }
        },
        {
            "name": "debug",
            "displayName": "Debug",
            "inherits": "base",
            "cacheVariables": {
                "CMAKE_BUILD_TYPE": "Debug"
            }
        },
        {
            "name": "release",
            "displayName": "Release",
            "inherits": "base",
            "cacheVariables": {
                "CMAKE_BUILD_TYPE": "Release"
            }
        },
        {
            "name": "relwithdebinfo",
            "displayName": "Release with Debug Info",
            "inherits": "base",
            "cacheVariables": {
                "CMAKE_BUILD_TYPE": "RelWithDebInfo"
            }
        }%s%s
    ],
    "buildPresets": [
        {
            "name": "debug",
            "configurePreset": "debug"
        },
        {
            "name": "release",
            "configurePreset": "release"
        },
        {
            "name": "relwithdebinfo",
            "configurePreset": "relwithdebinfo"
        }%s%s
    ],
    "testPresets": [
        {
            "name": "debug",
            "configurePreset": "debug",
            "output": {
                "outputOnFailure": true
            }
        },
        {
            "name": "release",
            "configurePreset": "release",
            "output": {
                "outputOnFailure": true
            }
        }
    ]
}
`, toolchainFile, sanitizerPresets, coveragePreset, sanitizerBuildPresets, coverageBuildPreset)
}

// SanitizersCMake generates cmake/Sanitizers.cmake
func SanitizersCMake() string {
	return `# Sanitizer configuration module
# Provides Address, Memory, Thread, and Undefined Behavior sanitizers

function(enable_sanitizers target)
    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU" OR CMAKE_CXX_COMPILER_ID MATCHES ".*Clang")
        set(SANITIZERS "")

        option(ENABLE_SANITIZER_ADDRESS "Enable address sanitizer" OFF)
        if(ENABLE_SANITIZER_ADDRESS)
            list(APPEND SANITIZERS "address")
        endif()

        option(ENABLE_SANITIZER_LEAK "Enable leak sanitizer" OFF)
        if(ENABLE_SANITIZER_LEAK)
            list(APPEND SANITIZERS "leak")
        endif()

        option(ENABLE_SANITIZER_UNDEFINED "Enable undefined behavior sanitizer" OFF)
        if(ENABLE_SANITIZER_UNDEFINED)
            list(APPEND SANITIZERS "undefined")
        endif()

        option(ENABLE_SANITIZER_THREAD "Enable thread sanitizer" OFF)
        if(ENABLE_SANITIZER_THREAD)
            if("address" IN_LIST SANITIZERS OR "leak" IN_LIST SANITIZERS)
                message(WARNING "Thread sanitizer cannot be used with Address or Leak sanitizer")
            else()
                list(APPEND SANITIZERS "thread")
            endif()
        endif()

        option(ENABLE_SANITIZER_MEMORY "Enable memory sanitizer (Clang only)" OFF)
        if(ENABLE_SANITIZER_MEMORY AND CMAKE_CXX_COMPILER_ID MATCHES ".*Clang")
            if("address" IN_LIST SANITIZERS
               OR "thread" IN_LIST SANITIZERS
               OR "leak" IN_LIST SANITIZERS)
                message(WARNING "Memory sanitizer cannot be used with Address, Thread, or Leak sanitizer")
            else()
                list(APPEND SANITIZERS "memory")
            endif()
        endif()

        if(SANITIZERS)
            list(JOIN SANITIZERS "," LIST_OF_SANITIZERS)
            message(STATUS "Enabling sanitizers: ${LIST_OF_SANITIZERS}")

            # Get target type to determine INTERFACE vs PRIVATE
            get_target_property(target_type ${target} TYPE)
            if(target_type STREQUAL "INTERFACE_LIBRARY")
                target_compile_options(${target} INTERFACE
                    -fsanitize=${LIST_OF_SANITIZERS}
                    -fno-omit-frame-pointer
                    -fno-optimize-sibling-calls
                )
                target_link_options(${target} INTERFACE -fsanitize=${LIST_OF_SANITIZERS})
            else()
                target_compile_options(${target} PRIVATE
                    -fsanitize=${LIST_OF_SANITIZERS}
                    -fno-omit-frame-pointer
                    -fno-optimize-sibling-calls
                )
                target_link_options(${target} PRIVATE -fsanitize=${LIST_OF_SANITIZERS})
            endif()
        endif()
    elseif(MSVC)
        option(ENABLE_SANITIZER_ADDRESS "Enable address sanitizer" OFF)
        if(ENABLE_SANITIZER_ADDRESS)
            message(STATUS "Enabling AddressSanitizer for MSVC")
            target_compile_options(${target} PRIVATE /fsanitize=address)
        endif()
    endif()
endfunction()
`
}

// CoverageCMake generates cmake/Coverage.cmake
func CoverageCMake() string {
	return `# Code coverage configuration module
# Supports GCC (gcov) and Clang (llvm-cov)

option(ENABLE_COVERAGE "Enable code coverage" OFF)

function(enable_coverage target)
    if(NOT ENABLE_COVERAGE)
        return()
    endif()

    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        message(STATUS "Enabling code coverage for GCC")
        target_compile_options(${target} PRIVATE --coverage -fprofile-arcs -ftest-coverage)
        target_link_options(${target} PRIVATE --coverage)
    elseif(CMAKE_CXX_COMPILER_ID MATCHES ".*Clang")
        message(STATUS "Enabling code coverage for Clang")
        target_compile_options(${target} PRIVATE -fprofile-instr-generate -fcoverage-mapping)
        target_link_options(${target} PRIVATE -fprofile-instr-generate -fcoverage-mapping)
    else()
        message(WARNING "Code coverage is not supported for ${CMAKE_CXX_COMPILER_ID}")
    endif()
endfunction()

# Custom target to generate coverage report
function(add_coverage_target)
    if(NOT ENABLE_COVERAGE)
        return()
    endif()

    find_program(LCOV lcov)
    find_program(GENHTML genhtml)
    find_program(LLVM_COV llvm-cov)
    find_program(LLVM_PROFDATA llvm-profdata)

    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU" AND LCOV AND GENHTML)
        add_custom_target(coverage
            COMMAND ${LCOV} --directory . --capture --output-file coverage.info
            COMMAND ${LCOV} --remove coverage.info '/usr/*' '*/tests/*' '*/build/*' --output-file coverage.info
            COMMAND ${GENHTML} coverage.info --output-directory coverage_report
            WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
            COMMENT "Generating code coverage report..."
        )
        message(STATUS "Coverage target available: cmake --build build --target coverage")
    elseif(CMAKE_CXX_COMPILER_ID MATCHES ".*Clang" AND LLVM_COV AND LLVM_PROFDATA)
        add_custom_target(coverage
            COMMAND ${LLVM_PROFDATA} merge -sparse default.profraw -o default.profdata
            COMMAND ${LLVM_COV} show ./tests -instr-profile=default.profdata -format=html -output-dir=coverage_report
            WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
            COMMENT "Generating code coverage report..."
        )
        message(STATUS "Coverage target available: cmake --build build --target coverage")
    else()
        message(WARNING "Coverage tools not found. Install lcov/genhtml (GCC) or llvm-cov/llvm-profdata (Clang)")
    endif()
endfunction()
`
}

// StaticAnalysisCMake generates cmake/StaticAnalysis.cmake
func StaticAnalysisCMake() string {
	return `# Static analysis configuration module
# Integrates clang-tidy, cppcheck, and include-what-you-use

option(ENABLE_CLANG_TIDY "Enable clang-tidy static analysis" OFF)
option(ENABLE_CPPCHECK "Enable cppcheck static analysis" OFF)
option(ENABLE_IWYU "Enable include-what-you-use" OFF)

function(enable_static_analysis target)
    # Clang-Tidy
    if(ENABLE_CLANG_TIDY)
        find_program(CLANG_TIDY clang-tidy)
        if(CLANG_TIDY)
            message(STATUS "Enabling clang-tidy for ${target}")
            set_target_properties(${target} PROPERTIES
                CXX_CLANG_TIDY "${CLANG_TIDY};--config-file=${CMAKE_SOURCE_DIR}/.clang-tidy"
            )
        else()
            message(WARNING "clang-tidy not found")
        endif()
    endif()

    # Cppcheck
    if(ENABLE_CPPCHECK)
        find_program(CPPCHECK cppcheck)
        if(CPPCHECK)
            message(STATUS "Enabling cppcheck for ${target}")
            set_target_properties(${target} PROPERTIES
                CXX_CPPCHECK "${CPPCHECK};--enable=all;--suppress=missingIncludeSystem;--inline-suppr;--inconclusive"
            )
        else()
            message(WARNING "cppcheck not found")
        endif()
    endif()

    # Include-what-you-use
    if(ENABLE_IWYU)
        find_program(IWYU include-what-you-use)
        if(IWYU)
            message(STATUS "Enabling include-what-you-use for ${target}")
            set_target_properties(${target} PROPERTIES
                CXX_INCLUDE_WHAT_YOU_USE "${IWYU}"
            )
        else()
            message(WARNING "include-what-you-use not found")
        endif()
    endif()
endfunction()
`
}

// DoxygenCMake generates cmake/Doxygen.cmake
func DoxygenCMake() string {
	return `# Doxygen documentation configuration

option(BUILD_DOCS "Build documentation" OFF)

function(enable_docs)
    if(NOT BUILD_DOCS)
        return()
    endif()

    find_package(Doxygen REQUIRED OPTIONAL_COMPONENTS dot)

    if(DOXYGEN_FOUND)
        set(DOXYGEN_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/docs")
        set(DOXYGEN_GENERATE_HTML YES)
        set(DOXYGEN_GENERATE_MAN NO)
        set(DOXYGEN_EXTRACT_ALL YES)
        set(DOXYGEN_EXTRACT_PRIVATE YES)
        set(DOXYGEN_EXTRACT_STATIC YES)
        set(DOXYGEN_RECURSIVE YES)
        set(DOXYGEN_USE_MDFILE_AS_MAINPAGE "${CMAKE_SOURCE_DIR}/README.md")
        set(DOXYGEN_EXCLUDE_PATTERNS "*/build/*" "*/tests/*" "*/_deps/*")

        # Modern theme settings
        set(DOXYGEN_HTML_COLORSTYLE_HUE 209)
        set(DOXYGEN_HTML_COLORSTYLE_SAT 255)
        set(DOXYGEN_HTML_COLORSTYLE_GAMMA 113)

        doxygen_add_docs(docs
            ${CMAKE_SOURCE_DIR}/include
            ${CMAKE_SOURCE_DIR}/src
            ${CMAKE_SOURCE_DIR}/README.md
            COMMENT "Generating API documentation with Doxygen"
        )

        message(STATUS "Doxygen documentation target available: cmake --build build --target docs")
    else()
        message(WARNING "Doxygen not found. Documentation will not be generated.")
    endif()
endfunction()
`
}

// Doxyfile generates a Doxyfile.in template
func Doxyfile(projectName, description string) string {
	return fmt.Sprintf(`# Doxyfile configuration for %s

PROJECT_NAME           = "%s"
PROJECT_BRIEF          = "%s"
PROJECT_NUMBER         = @PROJECT_VERSION@

OUTPUT_DIRECTORY       = @CMAKE_BINARY_DIR@/docs
INPUT                  = @CMAKE_SOURCE_DIR@/include @CMAKE_SOURCE_DIR@/src @CMAKE_SOURCE_DIR@/README.md
RECURSIVE              = YES
EXCLUDE_PATTERNS       = */build/* */tests/* */_deps/*

# Build settings
EXTRACT_ALL            = YES
EXTRACT_PRIVATE        = YES
EXTRACT_STATIC         = YES

# Output settings
GENERATE_HTML          = YES
GENERATE_LATEX         = NO
GENERATE_MAN           = NO

# HTML settings
HTML_OUTPUT            = html
HTML_COLORSTYLE_HUE    = 209
HTML_COLORSTYLE_SAT    = 255
HTML_COLORSTYLE_GAMMA  = 113

# Source browsing
SOURCE_BROWSER         = YES
INLINE_SOURCES         = NO
REFERENCED_BY_RELATION = YES
REFERENCES_RELATION    = YES

# Preprocessing
ENABLE_PREPROCESSING   = YES
MACRO_EXPANSION        = YES

# Diagrams (requires graphviz)
HAVE_DOT               = YES
DOT_IMAGE_FORMAT       = svg
INTERACTIVE_SVG        = YES
CALL_GRAPH             = YES
CALLER_GRAPH           = YES

USE_MDFILE_AS_MAINPAGE = @CMAKE_SOURCE_DIR@/README.md
`, projectName, projectName, description)
}
