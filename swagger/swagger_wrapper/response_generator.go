/*
 * Unicorn Rest Interface
 *
 * This is to provide rest interface to Unicorn. The rest API's can be used to execute and manage the tests.
 *
 * API version: 1.0.1
 * Contact: rakeshkumbi@telaverge.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger_wrapper

import (
	"encoding/json"
)

func generateJsonBodyForTests(tests []string) []byte {
	response, _ := json.Marshal(map[string][]string{"tests": tests})
	return response
}

func generateJsonBodyForFailureCause(failureCause string) []byte {
	response, _ := json.Marshal(map[string]string{"failureCause": failureCause})
	return response
}

func generateJsonBodyForAPIResponse(taskResponse string) []byte {
	response, _ := json.Marshal(map[string]string{"testStatus": taskResponse})
	return response
}

func generateJsonBodyForTestStatus(testStatus, additionalInfo string) []byte {
	response, _ := json.Marshal(map[string]string{
		"testStatus":     testStatus,
		"additionalInfo": additionalInfo,
	})
	return response
}

func generateJsonBodyForUnicornVersion(version string) []byte {
	response, _ := json.Marshal(map[string]string{
		"version": version,
	})
	return response
}

func generateJsonBodyForTestStats(stats string, timeElapsedInSec float64, testDuration int) []byte {
	response, _ := json.Marshal(map[string]interface{}{
		"testDuration": testDuration,
		"timeElapsedInSec":  timeElapsedInSec,
		"stats":        stats,
	})
	return response
}