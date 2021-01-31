package swagger_wrapper

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
	"unicorn/config"
	"unicorn/exectest"
	"unicorn/helper"
	"unicorn/testconfig"
)

//wrapperListSUTs - Lists all the SUTs available
func WrapperListSUTs(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	operation := splittedURLPath[2]
	if operation == "LIST_SUT" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		sutList, err := helper.ListDirectories(config.TestPath)
		if err != nil {
			// Internal server error
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForTests(sutList)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [LIST_SUT]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}

func WrapperListTestCases(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	sut, testsuite, operation := splittedURLPath[2], splittedURLPath[3], splittedURLPath[4]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !helper.IsPathExist(path.Join(config.TestPath, sut, testsuite)) {
		errorMessage := fmt.Sprintf("Test Suite '%s-%s' does't exist", sut, testsuite)
		w.WriteHeader(http.StatusInternalServerError)
		response := generateJsonBodyForFailureCause(errorMessage)
		w.Write(response)
		return
	}
	if operation == "LIST_TEST_CASES" {
		testCaseList, err := helper.ListDirectories(path.Join(config.TestPath, sut, testsuite))
		if err != nil {
			// Internal server error
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForTests(testCaseList)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [LIST_TEST_CASES]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}

func WrapperListTestSuites(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	sut, operation := splittedURLPath[2], splittedURLPath[3]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !helper.IsPathExist(path.Join(config.TestPath, sut)) {
		errorMessage := fmt.Sprintf("Test Suite '%s' does't exist", sut)
		w.WriteHeader(http.StatusInternalServerError)
		response := generateJsonBodyForFailureCause(errorMessage)
		w.Write(response)
		return
	}
	if operation == "LIST_TEST_SUITES" {
		testSuiteList, err := helper.ListDirectories(path.Join(config.TestPath, sut))
		if err != nil {
			// Internal server error
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForTests(testSuiteList)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [LIST_TEST_SUITES]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}

func WrapperStatsHandler(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	sut, testsuite, testcase, operation := splittedURLPath[2], splittedURLPath[3], splittedURLPath[4], splittedURLPath[5]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !helper.IsPathExist(path.Join(config.TestPath, sut, testsuite, testcase)) {
		errorMessage := fmt.Sprintf("Test case '%s-%s-%s' does't exist", sut, testsuite, testcase)
		w.WriteHeader(http.StatusInternalServerError)
		response := generateJsonBodyForFailureCause(errorMessage)
		w.Write(response)
		return
	}
	if operation == "TEST_STATUS" {
		testCaseStatus, failurCause, err := exectest.GetTestCaseStatus(sut, testsuite, testcase)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForTestStatus(testCaseStatus, failurCause)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else if operation == "TEST_STATISTICS" {
		if !exectest.IsTestRunning(sut, testsuite, testcase){
			w.WriteHeader(http.StatusInternalServerError)
			errorMessage := fmt.Sprintf("Test case '%s-%s-%s' is not running!", sut, testsuite, testcase)
                        response := generateJsonBodyForFailureCause(errorMessage)
                        w.Write(response)
                        return
		}
		stats, timeElapsed, err := exectest.GetCurrentTestStats(sut, testsuite, testcase)
		testDuration := testconfig.TestConf.Client.LoadParameters.TestDurationInSec
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForTestStats(stats, timeElapsed, testDuration)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [TEST_STATUS, TEST_STATISTICS]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}

func WrapperSutExecutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//TODO Have to implement
	w.WriteHeader(http.StatusOK)
}

//wrapperTestcaseExecutionHandler Handler for testcase execution related requests
func WrapperTestcaseExecutionHandler(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	sut, testsuite, testcase, operation := splittedURLPath[2], splittedURLPath[3], splittedURLPath[4], splittedURLPath[5]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !helper.IsPathExist(path.Join(config.TestPath, sut, testsuite, testcase)) {
		errorMessage := fmt.Sprintf("Test case '%s-%s-%s' does't exist", sut, testsuite, testcase)
		w.WriteHeader(http.StatusInternalServerError)
		response := generateJsonBodyForFailureCause(errorMessage)
		w.Write(response)
		return
	}
	if operation == "START_TEST" {
		if exectest.CheckIsAnyTestRunning() {
			errorMessage := fmt.Sprintf("Currently A test is running. Running multile test cases parallely not allowed")
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(errorMessage)
			w.Write(response)
			return
		}
		go exectest.ExecuteTestcase(sut, testsuite, testcase)
		time.Sleep(2 * time.Second)
		testCaseStatus, failurCause, err := exectest.GetTestCaseStatus(sut, testsuite, testcase)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		//fmt.Println(testCaseStatus)
		if testCaseStatus == "FAILED" {
			// Internal server error
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(failurCause)
			w.Write(response)
			return
		}
		message := fmt.Sprintf("Test case '%s-%s-%s' Started Executing", sut, testsuite, testcase)
		response := generateJsonBodyForAPIResponse(message)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else if operation == "STOP_TEST" {
		if !exectest.CheckIsAnyTestRunning() {
			errorMessage := fmt.Sprintf("No test is Running to stop")
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(errorMessage)
			w.Write(response)
			return
		}
		err := exectest.StopTestCase(sut, testsuite, testcase)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		message := fmt.Sprintf("Test case '%s-%s-%s' Stopped", sut, testsuite, testcase)
		response := generateJsonBodyForAPIResponse(message)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [START_TEST, STOP_TEST]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}

func WrapperTestsuiteExecutionHandler(w http.ResponseWriter, r *http.Request) {
	//TODO Have to implement
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
