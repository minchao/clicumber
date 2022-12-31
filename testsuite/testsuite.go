/*
Copyright (C) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testsuite

import (
	"context"
	"fmt"
	"os"

	"github.com/code-ready/clicumber/util"
	"github.com/cucumber/godog"
)

var (
	testDir         string
	testRunDir      string
	testResultsDir  string
	testDefaultHome string
	testWithShell   string

	GodogFormat              string
	GodogTags                string
	GodogShowStepDefinitions bool
	GodogStopOnFailure       bool
	GodogNoColors            bool
	GodogPaths               string
)

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		err := PrepareForE2eTest()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
	ctx.AfterSuite(func() {
		_ = util.LogMessage("info", "----- Cleaning Up -----")
		err := util.CloseLog()
		if err != nil {
			fmt.Println("Error closing the log:", err)
		}
	})
}

// InitializeScenario defines godog.ScenarioContext steps for the test scenario.
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Executing commands
	ctx.Step(`^executing "(.*)"$`,
		ExecuteCommand)
	ctx.Step(`^executing "(.*)" (succeeds|fails)$`,
		ExecuteCommandSucceedsOrFails)

	// Command output verification
	ctx.Step(`^(stdout|stderr|exitcode) (?:should contain|contains) "(.*)"$`,
		CommandReturnShouldContain)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should contain|contains)$`,
		CommandReturnShouldContainContent)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does) not contain "(.*)"$`,
		CommandReturnShouldNotContain)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does not) contain$`,
		CommandReturnShouldNotContainContent)

	ctx.Step(`^(stdout|stderr|exitcode) (?:should equal|equals) "(.*)"$`,
		CommandReturnShouldEqual)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should equal|equals)$`,
		CommandReturnShouldEqualContent)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does) not equal "(.*)"$`,
		CommandReturnShouldNotEqual)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does) not equal$`,
		CommandReturnShouldNotEqualContent)

	ctx.Step(`^(stdout|stderr|exitcode) (?:should match|matches) "(.*)"$`,
		CommandReturnShouldMatch)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should match|matches)`,
		CommandReturnShouldMatchContent)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does) not match "(.*)"$`,
		CommandReturnShouldNotMatch)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should|does) not match`,
		CommandReturnShouldNotMatchContent)

	ctx.Step(`^(stdout|stderr|exitcode) (?:should be|is) empty$`,
		CommandReturnShouldBeEmpty)
	ctx.Step(`^(stdout|stderr|exitcode) (?:should not be|is not) empty$`,
		CommandReturnShouldNotBeEmpty)

	ctx.Step(`^(stdout|stderr|exitcode) (?:should be|is) valid "([^"]*)"$`,
		ShouldBeInValidFormat)

	// Command output and execution: extra steps
	ctx.Step(`^with up to "(\d*)" retries with wait period of "(\d*(?:ms|s|m))" command "(.*)" output (should contain|contains|should not contain|does not contain) "(.*)"$`,
		ExecuteCommandWithRetry)
	ctx.Step(`^evaluating stdout of the previous command succeeds$`,
		ExecuteStdoutLineByLine)

	// Scenario variables
	// allows to set a scenario variable to the output values of minishift and oc commands
	// and then refer to it by $(NAME_OF_VARIABLE) directly in the text of feature file
	ctx.Step(`^setting scenario variable "(.*)" to the stdout from executing "(.*)"$`,
		SetScenarioVariableExecutingCommand)

	// Filesystem operations
	ctx.Step(`^creating directory "([^"]*)" succeeds$`,
		CreateDirectory)
	ctx.Step(`^creating file "([^"]*)" succeeds$`,
		CreateFile)
	ctx.Step(`^deleting directory "([^"]*)" succeeds$`,
		DeleteDirectory)
	ctx.Step(`^deleting file "([^"]*)" succeeds$`,
		DeleteFile)
	ctx.Step(`^directory "([^"]*)" should not exist$`,
		DirectoryShouldNotExist)
	ctx.Step(`^file "([^"]*)" should not exist$`,
		FileShouldNotExist)
	ctx.Step(`^file "([^"]*)" exists$`,
		FileExist)
	ctx.Step(`^file from "(.*)" is downloaded into location "(.*)"$`,
		DownloadFileIntoLocation)
	ctx.Step(`^writing text "([^"]*)" to file "([^"]*)" succeeds$`,
		WriteToFile)

	// File content checks
	ctx.Step(`^content of file "([^"]*)" should contain "([^"]*)"$`,
		FileContentShouldContain)
	ctx.Step(`^content of file "([^"]*)" should not contain "([^"]*)"$`,
		FileContentShouldNotContain)
	ctx.Step(`^content of file "([^"]*)" should equal "([^"]*)"$`,
		FileContentShouldEqual)
	ctx.Step(`^content of file "([^"]*)" should not equal "([^"]*)"$`,
		FileContentShouldNotEqual)
	ctx.Step(`^content of file "([^"]*)" should match "([^"]*)"$`,
		FileContentShouldMatchRegex)
	ctx.Step(`^content of file "([^"]*)" should not match "([^"]*)"$`,
		FileContentShouldNotMatchRegex)
	ctx.Step(`^content of file "([^"]*)" (?:should be|is) valid "([^"]*)"$`,
		FileContentIsInValidFormat)

	// Config file content, JSON and YAML
	ctx.Step(`"(JSON|YAML)" config file "(.*)" (contains|does not contain) key "(.*)" with value matching "(.*)"$`,
		ConfigFileContainsKeyMatchingValue)
	ctx.Step(`"(JSON|YAML)" config file "(.*)" (contains|does not contain) key "(.*)"$`,
		ConfigFileContainsKey)

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		_ = util.LogMessage("info", fmt.Sprintf("----- Scenario: %s -----", sc.Name))
		_ = util.LogMessage("info", fmt.Sprintf("----- Scenario Outline: %s -----", sc))
		_ = StartHostShellInstance(testWithShell)
		util.ClearScenarioVariables()
		err := CleanTestRunDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return ctx, nil
	})

	ctx.StepContext().Before(func(ctx context.Context, st *godog.Step) (context.Context, error) {
		st.Text = util.ProcessScenarioVariables(st.Text)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_ = util.LogMessage("info", "----- Cleaning after scenario -----")
		_ = CloseHostShellInstance()
		return ctx, nil
	})
}
