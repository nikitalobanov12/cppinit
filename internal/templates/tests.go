package templates

import "fmt"

// TestsCMakeLists generates the tests/CMakeLists.txt
func TestsCMakeLists(projectName, projectType, testFramework string, isC bool) string {
	// Only link against library if it's a library project
	linkLib := ""
	if projectType == "static" || projectType == "header-only" {
		linkLib = fmt.Sprintf("\n        %s", projectName)
	}

	// C Unity framework
	if testFramework == "unity" {
		srcExt := ".c"
		return fmt.Sprintf(`include(FetchContent)

FetchContent_Declare(
    unity
    GIT_REPOSITORY https://github.com/ThrowTheSwitch/Unity.git
    GIT_TAG v2.6.0
)
FetchContent_MakeAvailable(unity)

add_executable(tests
    test_main%s
)

target_link_libraries(tests
    PRIVATE
        unity%s
)

target_include_directories(tests
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
)

add_test(NAME tests COMMAND tests)
`, srcExt, linkLib)
	}

	if testFramework == "googletest" {
		return fmt.Sprintf(`include(FetchContent)

FetchContent_Declare(
    googletest
    GIT_REPOSITORY https://github.com/google/googletest.git
    GIT_TAG v1.14.0
)

# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

add_executable(tests
    test_main.cpp
)

target_link_libraries(tests
    PRIVATE
        GTest::gtest_main%s
)

target_include_directories(tests
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
)

include(GoogleTest)
gtest_discover_tests(tests)
`, linkLib)
	}

	if testFramework == "doctest" {
		return fmt.Sprintf(`include(FetchContent)

FetchContent_Declare(
    doctest
    GIT_REPOSITORY https://github.com/doctest/doctest.git
    GIT_TAG v2.4.11
)
FetchContent_MakeAvailable(doctest)

add_executable(tests
    test_main.cpp
)

target_link_libraries(tests
    PRIVATE
        doctest::doctest%s
)

target_include_directories(tests
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
)

include(CTest)
include(${doctest_SOURCE_DIR}/scripts/cmake/doctest.cmake)
doctest_discover_tests(tests)
`, linkLib)
	}

	// Catch2 (default for C++)
	return fmt.Sprintf(`include(FetchContent)

FetchContent_Declare(
    Catch2
    GIT_REPOSITORY https://github.com/catchorg/Catch2.git
    GIT_TAG v3.5.2
)
FetchContent_MakeAvailable(Catch2)

add_executable(tests
    test_main.cpp
)

target_link_libraries(tests
    PRIVATE
        Catch2::Catch2WithMain%s
)

target_include_directories(tests
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
)

include(CTest)
include(Catch)
catch_discover_tests(tests)
`, linkLib)
}

// TestMainCpp generates the test file
func TestMainCpp(projectName, projectType, testFramework string) string {
	// For executable projects, just provide a basic test without library includes
	if projectType == "executable" {
		if testFramework == "googletest" {
			return fmt.Sprintf(`#include <gtest/gtest.h>

TEST(%sTest, BasicAssertion) {
    EXPECT_EQ(1, 1);
}

TEST(%sTest, SampleTest) {
    // Add your tests here
    EXPECT_TRUE(true);
}
`, projectName, projectName)
		}
		if testFramework == "doctest" {
			return fmt.Sprintf(`#define DOCTEST_CONFIG_IMPLEMENT_WITH_MAIN
#include <doctest/doctest.h>

TEST_CASE("%s basic tests") {
    SUBCASE("Basic assertion") {
        CHECK(1 == 1);
    }

    SUBCASE("Sample test") {
        // Add your tests here
        CHECK(true);
    }
}
`, projectName)
		}
		// Catch2 for executable
		return fmt.Sprintf(`#include <catch2/catch_test_macros.hpp>

TEST_CASE("%s basic tests", "[%s]") {
    SECTION("Basic assertion") {
        REQUIRE(1 == 1);
    }

    SECTION("Sample test") {
        // Add your tests here
        REQUIRE(true);
    }
}
`, projectName, projectName)
	}

	// For library projects, include the library header
	if testFramework == "googletest" {
		return fmt.Sprintf(`#include <gtest/gtest.h>
#include "%s/%s.hpp"

TEST(%sTest, BasicAssertion) {
    EXPECT_EQ(1, 1);
}

TEST(%sTest, AddFunction) {
    EXPECT_EQ(%s::add(2, 3), 5);
    EXPECT_EQ(%s::add(-1, 1), 0);
}
`, projectName, projectName, projectName, projectName, projectName, projectName)
	}

	if testFramework == "doctest" {
		return fmt.Sprintf(`#define DOCTEST_CONFIG_IMPLEMENT_WITH_MAIN
#include <doctest/doctest.h>
#include "%s/%s.hpp"

TEST_CASE("%s basic tests") {
    SUBCASE("Basic assertion") {
        CHECK(1 == 1);
    }

    SUBCASE("Add function") {
        CHECK(%s::add(2, 3) == 5);
        CHECK(%s::add(-1, 1) == 0);
    }
}
`, projectName, projectName, projectName, projectName, projectName)
	}

	// Catch2 for library
	return fmt.Sprintf(`#include <catch2/catch_test_macros.hpp>
#include "%s/%s.hpp"

TEST_CASE("%s basic tests", "[%s]") {
    SECTION("Basic assertion") {
        REQUIRE(1 == 1);
    }

    SECTION("Add function") {
        REQUIRE(%s::add(2, 3) == 5);
        REQUIRE(%s::add(-1, 1) == 0);
    }
}
`, projectName, projectName, projectName, projectName, projectName, projectName)
}

// TestMainC generates the C test file (Unity framework)
func TestMainC(projectName, projectType, testFramework string) string {
	if projectType == "executable" {
		return `#include "unity.h"

void setUp(void) {
    // Set up code here (runs before each test)
}

void tearDown(void) {
    // Tear down code here (runs after each test)
}

void test_basic_assertion(void) {
    TEST_ASSERT_EQUAL(1, 1);
}

void test_sample(void) {
    // Add your tests here
    TEST_ASSERT_TRUE(1);
}

int main(void) {
    UNITY_BEGIN();
    RUN_TEST(test_basic_assertion);
    RUN_TEST(test_sample);
    return UNITY_END();
}
`
	}

	// For library projects, include the library header
	return fmt.Sprintf(`#include "unity.h"
#include "%s/%s.h"

void setUp(void) {
    // Set up code here (runs before each test)
}

void tearDown(void) {
    // Tear down code here (runs after each test)
}

void test_basic_assertion(void) {
    TEST_ASSERT_EQUAL(1, 1);
}

void test_add_function(void) {
    TEST_ASSERT_EQUAL(5, %s_add(2, 3));
    TEST_ASSERT_EQUAL(0, %s_add(-1, 1));
}

int main(void) {
    UNITY_BEGIN();
    RUN_TEST(test_basic_assertion);
    RUN_TEST(test_add_function);
    return UNITY_END();
}
`, projectName, projectName, projectName, projectName)
}

// VcpkgJson generates vcpkg.json manifest
func VcpkgJson(projectName, testFramework string) string {
	deps := ""
	if testFramework == "googletest" {
		deps = `
    "dependencies": [
        "gtest"
    ]`
	} else if testFramework == "catch2" {
		deps = `
    "dependencies": [
        "catch2"
    ]`
	}

	return fmt.Sprintf(`{
    "name": "%s",
    "version-string": "0.1.0",
    "description": "A C++ project"%s
}
`, projectName, deps)
}

// CMakePresetsVcpkg generates CMakePresets.json for vcpkg
func CMakePresetsVcpkg() string {
	return `{
    "version": 6,
    "cmakeMinimumRequired": {
        "major": 3,
        "minor": 20,
        "patch": 0
    },
    "configurePresets": [
        {
            "name": "default",
            "hidden": true,
            "binaryDir": "${sourceDir}/build/${presetName}",
            "toolchainFile": "$env{VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake"
        },
        {
            "name": "debug",
            "inherits": "default",
            "cacheVariables": {
                "CMAKE_BUILD_TYPE": "Debug"
            }
        },
        {
            "name": "release",
            "inherits": "default",
            "cacheVariables": {
                "CMAKE_BUILD_TYPE": "Release"
            }
        }
    ],
    "buildPresets": [
        {
            "name": "debug",
            "configurePreset": "debug"
        },
        {
            "name": "release",
            "configurePreset": "release"
        }
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
`
}

// ConanfileTxt generates conanfile.txt
func ConanfileTxt(testFramework string) string {
	deps := ""
	if testFramework == "googletest" {
		deps = "gtest/1.14.0"
	} else if testFramework == "catch2" {
		deps = "catch2/3.5.2"
	}

	return fmt.Sprintf(`[requires]
%s

[generators]
CMakeDeps
CMakeToolchain

[layout]
cmake_layout
`, deps)
}

// GitHubActionsCI generates GitHub Actions workflow
func GitHubActionsCI(packageManager, testFramework string) string {
	testStep := ""
	if testFramework != "none" {
		testStep = `
      - name: Test
        working-directory: build
        run: ctest --output-on-failure`
	}

	vcpkgSetup := ""
	cmakeArgs := ""
	if packageManager == "vcpkg" {
		vcpkgSetup = `
      - name: Setup vcpkg
        uses: lukka/run-vcpkg@v11
        with:
          vcpkgGitCommitId: 'a34c873a9717a888f58dc05268dea15592c2f0ff'`
		cmakeArgs = " --preset debug"
	}

	return fmt.Sprintf(`name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  build:
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]

    steps:
      - uses: actions/checkout@v4
%s
      - name: Configure
        run: cmake -B build -S .%s

      - name: Build
        run: cmake --build build
%s
`, vcpkgSetup, cmakeArgs, testStep)
}
