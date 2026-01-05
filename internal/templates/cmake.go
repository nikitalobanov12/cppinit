package templates

import "fmt"

// RootCMakeLists generates the main CMakeLists.txt
func RootCMakeLists(projectName, cppStandard, projectType, testFramework, packageManager string) string {
	var vcpkgToolchain string
	if packageManager == "vcpkg" {
		vcpkgToolchain = `# vcpkg toolchain is handled via CMakePresets.json
`
	}

	testSection := ""
	if testFramework != "none" {
		testSection = `
# Testing
option(BUILD_TESTS "Build the tests" ON)
if(BUILD_TESTS)
    enable_testing()
    add_subdirectory(tests)
endif()
`
	}

	var targetSection string
	switch projectType {
	case "executable":
		targetSection = `
# Main executable
add_executable(${PROJECT_NAME} src/main.cpp)

target_include_directories(${PROJECT_NAME}
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)
`
	case "static":
		targetSection = fmt.Sprintf(`
# Library
add_library(${PROJECT_NAME} STATIC
    src/%s.cpp
)

target_include_directories(${PROJECT_NAME}
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)
`, projectName)
	case "header-only":
		targetSection = `
# Header-only library
add_library(${PROJECT_NAME} INTERFACE)

target_include_directories(${PROJECT_NAME}
    INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)
`
	}

	return fmt.Sprintf(`cmake_minimum_required(VERSION 3.20)
%s
project(%s
    VERSION 0.1.0
    DESCRIPTION "A C++ project"
    LANGUAGES CXX
)

# Set C++ standard
set(CMAKE_CXX_STANDARD %s)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for IDE support
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Include custom CMake modules
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(CompilerWarnings)
%s
# Apply compiler warnings
set_project_warnings(${PROJECT_NAME})
%s`, vcpkgToolchain, projectName, cppStandard, targetSection, testSection)
}

// CompilerWarningsCMake generates the compiler warnings CMake module
func CompilerWarningsCMake() string {
	return `# Set compiler warnings for a target
function(set_project_warnings target)
    set(MSVC_WARNINGS
        /W4          # Baseline reasonable warnings
        /w14242      # 'identifier': conversion from 'type1' to 'type2', possible loss of data
        /w14254      # 'operator': conversion from 'type1:field_bits' to 'type2:field_bits'
        /w14263      # 'function': member function does not override any base class virtual member function
        /w14265      # 'classname': class has virtual functions, but destructor is not virtual
        /w14287      # 'operator': unsigned/negative constant mismatch
        /we4289      # nonstandard extension used: 'variable': loop control variable declared in the for-loop is used outside the for-loop scope
        /w14296      # 'operator': expression is always 'boolean_value'
        /w14311      # 'variable': pointer truncation from 'type1' to 'type2'
        /w14545      # expression before comma evaluates to a function which is missing an argument list
        /w14546      # function call before comma missing argument list
        /w14547      # 'operator': operator before comma has no effect; expected operator with side-effect
        /w14549      # 'operator': operator before comma has no effect; did you intend 'operator'?
        /w14555      # expression has no effect; expected expression with side-effect
        /w14619      # pragma warning: there is no warning number 'number'
        /w14640      # Enable warning on thread un-safe static member initialization
        /w14826      # Conversion from 'type1' to 'type2' is sign-extended
        /w14905      # wide string literal cast to 'LPSTR'
        /w14906      # string literal cast to 'LPWSTR'
        /w14928      # illegal copy-initialization; more than one user-defined conversion has been implicitly applied
        /permissive- # standards conformance mode
    )

    set(CLANG_WARNINGS
        -Wall
        -Wextra              # reasonable and standard
        -Wshadow             # warn if a variable declaration shadows one from a parent context
        -Wnon-virtual-dtor   # warn if a class with virtual functions has a non-virtual destructor
        -Wold-style-cast     # warn for c-style casts
        -Wcast-align         # warn for potential performance problem casts
        -Wunused             # warn on anything being unused
        -Woverloaded-virtual # warn if you overload (not override) a virtual function
        -Wpedantic           # warn if non-standard C++ is used
        -Wconversion         # warn on type conversions that may lose data
        -Wsign-conversion    # warn on sign conversions
        -Wnull-dereference   # warn if a null dereference is detected
        -Wdouble-promotion   # warn if float is implicit promoted to double
        -Wformat=2           # warn on security issues around functions that format output
        -Wimplicit-fallthrough # warn on missing break in switch
    )

    set(GCC_WARNINGS
        ${CLANG_WARNINGS}
        -Wmisleading-indentation # warn if indentation implies blocks where blocks do not exist
        -Wduplicated-cond        # warn if if / else chain has duplicated conditions
        -Wduplicated-branches    # warn if if / else branches have duplicated code
        -Wlogical-op             # warn about logical operations being used where bitwise were probably wanted
        -Wuseless-cast           # warn if you perform a cast to the same type
    )

    if(MSVC)
        set(PROJECT_WARNINGS ${MSVC_WARNINGS})
    elseif(CMAKE_CXX_COMPILER_ID MATCHES ".*Clang")
        set(PROJECT_WARNINGS ${CLANG_WARNINGS})
    elseif(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        set(PROJECT_WARNINGS ${GCC_WARNINGS})
    else()
        message(AUTHOR_WARNING "No compiler warnings set for '${CMAKE_CXX_COMPILER_ID}' compiler.")
    endif()

    # Check if target is INTERFACE (header-only library)
    get_target_property(target_type ${target} TYPE)
    if(target_type STREQUAL "INTERFACE_LIBRARY")
        target_compile_options(${target} INTERFACE ${PROJECT_WARNINGS})
    else()
        target_compile_options(${target} PRIVATE ${PROJECT_WARNINGS})
    endif()
endfunction()
`
}

// GitIgnore generates a .gitignore file
func GitIgnore() string {
	return `# Build directories
build/
cmake-build-*/
out/

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# Compiled files
*.o
*.obj
*.exe
*.out
*.app
*.so
*.dylib
*.dll
*.a
*.lib

# CMake
CMakeCache.txt
CMakeFiles/
cmake_install.cmake
Makefile
compile_commands.json

# Package managers
vcpkg_installed/
conan/

# Testing
Testing/
CTestTestfile.cmake

# OS
.DS_Store
Thumbs.db
`
}

// Readme generates a README.md file
func Readme(projectName, cppStandard, projectType string) string {
	typeDesc := "application"
	if projectType == "static" || projectType == "header-only" {
		typeDesc = "library"
	}

	return fmt.Sprintf(`# %s

A C++%s %s.

## Building

`+"```bash"+`
mkdir build && cd build
cmake ..
cmake --build .
`+"```"+`

## Requirements

- CMake 3.20 or higher
- C++%s compatible compiler

## License

MIT
`, projectName, cppStandard, typeDesc, cppStandard)
}

// MainCpp generates main.cpp for executable projects
func MainCpp(projectName string) string {
	return fmt.Sprintf(`#include <iostream>

int main() {
    std::cout << "Hello from %s!" << std::endl;
    return 0;
}
`, projectName)
}

// LibraryCpp generates the library source file
func LibraryCpp(projectName string) string {
	return fmt.Sprintf(`#include "%s/%s.hpp"

namespace %s {

int add(int a, int b) {
    return a + b;
}

} // namespace %s
`, projectName, projectName, projectName, projectName)
}

// LibraryHpp generates the library header file
func LibraryHpp(projectName string) string {
	upperName := toUpperSnake(projectName)
	return fmt.Sprintf(`#ifndef %s_HPP
#define %s_HPP

namespace %s {

/// Adds two integers
/// @param a First operand
/// @param b Second operand
/// @return Sum of a and b
int add(int a, int b);

} // namespace %s

#endif // %s_HPP
`, upperName, upperName, projectName, projectName, upperName)
}

// HeaderOnlyHpp generates a header-only library template
func HeaderOnlyHpp(projectName string) string {
	upperName := toUpperSnake(projectName)
	return fmt.Sprintf(`#ifndef %s_HPP
#define %s_HPP

namespace %s {

/// Adds two integers
/// @param a First operand
/// @param b Second operand
/// @return Sum of a and b
template<typename T>
constexpr T add(T a, T b) {
    return a + b;
}

} // namespace %s

#endif // %s_HPP
`, upperName, upperName, projectName, projectName, upperName)
}

// Helper to convert to UPPER_SNAKE_CASE
func toUpperSnake(s string) string {
	result := ""
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result += "_"
			}
			result += string(r)
		} else if r >= 'a' && r <= 'z' {
			result += string(r - 32) // Convert to uppercase
		} else if r == '-' || r == ' ' {
			result += "_"
		} else {
			result += string(r)
		}
	}
	return result
}
