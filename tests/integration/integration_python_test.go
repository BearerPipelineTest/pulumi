// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
// +build python all

package ints

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/pulumi/pulumi/sdk/v2/go/common/resource"
	"github.com/stretchr/testify/assert"
)

// TestEmptyPython simply tests that we can run an empty Python project.
func TestEmptyPython(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("empty", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
	})
}

// TestEmptyPythonVenv simply tests that we can run an empty Python project using automatic virtual environment support.
func TestEmptyPythonVenv(t *testing.T) {
	t.Skip("Temporarily skipping test - pulumi/pulumi#4849")
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("empty", "python_venv"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick:                  true,
		UseAutomaticVirtualEnv: true,
	})
}

func TestStackOutputsPython(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("stack_outputs", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Ensure the checkpoint contains a single resource, the Stack, with two outputs.
			fmt.Printf("Deployment: %v", stackInfo.Deployment)
			assert.NotNil(t, stackInfo.Deployment)
			if assert.Equal(t, 1, len(stackInfo.Deployment.Resources)) {
				stackRes := stackInfo.Deployment.Resources[0]
				assert.NotNil(t, stackRes)
				assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
				assert.Equal(t, 0, len(stackRes.Inputs))
				assert.Equal(t, 2, len(stackRes.Outputs))
				assert.Equal(t, "ABC", stackRes.Outputs["xyz"])
				assert.Equal(t, float64(42), stackRes.Outputs["foo"])
			}
		},
	})
}

// Tests basic configuration from the perspective of a Pulumi program.
func TestConfigBasicPython(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("config_basic", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Config: map[string]string{
			"aConfigValue": "this value is a Pythonic value",
		},
		Secrets: map[string]string{
			"bEncryptedSecret": "this super Pythonic secret is encrypted",
		},
		OrderedConfig: []integration.ConfigValue{
			{Key: "outer.inner", Value: "value", Path: true},
			{Key: "names[0]", Value: "a", Path: true},
			{Key: "names[1]", Value: "b", Path: true},
			{Key: "names[2]", Value: "c", Path: true},
			{Key: "names[3]", Value: "super secret name", Path: true, Secret: true},
			{Key: "servers[0].port", Value: "80", Path: true},
			{Key: "servers[0].host", Value: "example", Path: true},
			{Key: "a.b[0].c", Value: "true", Path: true},
			{Key: "a.b[1].c", Value: "false", Path: true},
			{Key: "tokens[0]", Value: "shh", Path: true, Secret: true},
			{Key: "foo.bar", Value: "don't tell", Path: true, Secret: true},
		},
	})
}

func TestConfigBasicPythonVenv(t *testing.T) {
	t.Skip("Temporarily skipping test - pulumi/pulumi#4849")
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("config_basic", "python_venv"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Config: map[string]string{
			"aConfigValue": "this value is a Pythonic value",
		},
		Secrets: map[string]string{
			"bEncryptedSecret": "this super Pythonic secret is encrypted",
		},
		OrderedConfig: []integration.ConfigValue{
			{Key: "outer.inner", Value: "value", Path: true},
			{Key: "names[0]", Value: "a", Path: true},
			{Key: "names[1]", Value: "b", Path: true},
			{Key: "names[2]", Value: "c", Path: true},
			{Key: "names[3]", Value: "super secret name", Path: true, Secret: true},
			{Key: "servers[0].port", Value: "80", Path: true},
			{Key: "servers[0].host", Value: "example", Path: true},
			{Key: "a.b[0].c", Value: "true", Path: true},
			{Key: "a.b[1].c", Value: "false", Path: true},
			{Key: "tokens[0]", Value: "shh", Path: true, Secret: true},
			{Key: "foo.bar", Value: "don't tell", Path: true, Secret: true},
		},
		UseAutomaticVirtualEnv: true,
	})
}

func TestStackReferencePython(t *testing.T) {
	if runtime.GOOS == WindowsOS {
		t.Skip("Temporarily skipping test on Windows - pulumi/pulumi#3811")
	}
	if owner := os.Getenv("PULUMI_TEST_OWNER"); owner == "" {
		t.Skipf("Skipping: PULUMI_TEST_OWNER is not set")
	}

	opts := &integration.ProgramTestOptions{
		Dir: filepath.Join("stack_reference", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Config: map[string]string{
			"org": os.Getenv("PULUMI_TEST_OWNER"),
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step1",
				Additive: true,
			},
			{
				Dir:      "step2",
				Additive: true,
			},
		},
	}
	integration.ProgramTest(t, opts)
}

func TestMultiStackReferencePython(t *testing.T) {
	if runtime.GOOS == WindowsOS {
		t.Skip("Temporarily skipping test on Windows - pulumi/pulumi#3811")
	}
	if owner := os.Getenv("PULUMI_TEST_OWNER"); owner == "" {
		t.Skipf("Skipping: PULUMI_TEST_OWNER is not set")
	}

	// build a stack with an export
	exporterOpts := &integration.ProgramTestOptions{
		Dir: filepath.Join("stack_reference_multi", "python", "exporter"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Config: map[string]string{
			"org": os.Getenv("PULUMI_TEST_OWNER"),
		},
		NoParallel: true,
	}

	// we're going to manually initialize and then defer the deletion of this stack
	exporterPt := integration.ProgramTestManualLifeCycle(t, exporterOpts)
	exporterPt.TestFinished = false
	err := exporterPt.TestLifeCyclePrepare()
	assert.NoError(t, err)
	err = exporterPt.TestLifeCycleInitialize()
	assert.NoError(t, err)

	defer func() {
		destroyErr := exporterPt.TestLifeCycleDestroy()
		assert.NoError(t, destroyErr)
		exporterPt.TestFinished = true
		exporterPt.TestCleanUp()
	}()

	err = exporterPt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	exporterStackName := exporterOpts.GetStackName().String()

	importerOpts := &integration.ProgramTestOptions{
		Dir: filepath.Join("stack_reference_multi", "python", "importer"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Config: map[string]string{
			"org":                 os.Getenv("PULUMI_TEST_OWNER"),
			"exporter_stack_name": exporterStackName,
		},
		NoParallel: true,
	}
	integration.ProgramTest(t, importerOpts)
}

// Tests that we issue an error if we fail to locate the Python command when running
// a Python example.
func TestPython3NotInstalled(t *testing.T) {
	stderr := &bytes.Buffer{}
	badPython := "python3000"
	expectedError := fmt.Sprintf(
		"error: Failed to locate any of %q on your PATH.  Have you installed Python 3.6 or greater?",
		[]string{badPython})
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("empty", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Quick: true,
		Env: []string{
			// Note: we use PULUMI_PYTHON_CMD to override the default behavior of searching
			// for Python 3, since anyone running tests surely already has Python 3 installed on their
			// machine. The code paths are functionally the same.
			fmt.Sprintf("PULUMI_PYTHON_CMD=%s", badPython),
		},
		ExpectFailure: true,
		Stderr:        stderr,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			output := stderr.String()
			assert.Contains(t, output, expectedError)
		},
	})
}

// Tests dynamic provider in Python.
func TestDynamicPython(t *testing.T) {
	var randomVal string
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("dynamic", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			randomVal = stack.Outputs["random_val"].(string)
		},
		EditDirs: []integration.EditDir{{
			Dir:      "step1",
			Additive: true,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.Equal(t, randomVal, stack.Outputs["random_val"].(string))
			},
		}},
	})
}

// Tests dynamic provider in Python using automatic virtual environment support.
func TestDynamicPythonVenv(t *testing.T) {
	t.Skip("Temporarily skipping test - pulumi/pulumi#4849")
	var randomVal string
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("dynamic", "python_venv"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			randomVal = stack.Outputs["random_val"].(string)
		},
		EditDirs: []integration.EditDir{{
			Dir:      "step1",
			Additive: true,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.Equal(t, randomVal, stack.Outputs["random_val"].(string))
			},
		}},
		UseAutomaticVirtualEnv: true,
	})
}

func TestPartialValuesPython(t *testing.T) {
	if runtime.GOOS == WindowsOS {
		t.Skip("Temporarily skipping test on Windows - pulumi/pulumi#3811")
	}
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir: filepath.Join("partial_values", "python"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		AllowEmptyPreviewChanges: true,
	})
}

// Tests a resource with a large (>4mb) string prop in Python
func TestLargeResourcePython(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		Dir: filepath.Join("large_resource", "python"),
	})
}

// Test to ensure Pylint is clean.
func TestPythonPylint(t *testing.T) {
	t.Skip("Temporarily skipping test - pulumi/pulumi#4849")
	var opts *integration.ProgramTestOptions
	opts = &integration.ProgramTestOptions{
		Dir: filepath.Join("python", "pylint"),
		Dependencies: []string{
			filepath.Join("..", "..", "sdk", "python", "env", "src"),
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			randomURN := stack.Outputs["random_urn"].(string)
			assert.NotEmpty(t, randomURN)

			randomID := stack.Outputs["random_id"].(string)
			randomVal := stack.Outputs["random_val"].(string)
			assert.Equal(t, randomID, randomVal)

			cwd := stack.Outputs["cwd"].(string)
			assert.NotEmpty(t, cwd)

			pylint := filepath.Join("venv", "bin", "pylint")
			if runtime.GOOS == WindowsOS {
				pylint = filepath.Join("venv", "Scripts", "pylint")
			}

			err := integration.RunCommand(t, "pylint", []string{pylint, "__main__.py"}, cwd, opts)
			assert.NoError(t, err)
		},
		Quick:                  true,
		UseAutomaticVirtualEnv: true,
	}
	integration.ProgramTest(t, opts)
}