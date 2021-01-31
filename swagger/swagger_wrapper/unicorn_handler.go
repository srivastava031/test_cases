
package swagger_wrapper

import (
	"fmt"
	"net/http"
	"unicorn/helper"
	"strings"
)

//wrapperListSUTs - Lists all the SUTs available
func WrapperUnicornOperationHandler(w http.ResponseWriter, r *http.Request) {
	splittedURLPath := strings.Split(r.URL.Path, "/")
	operation := splittedURLPath[2]
	if operation == "GET_VERSION" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		version, err := helper.GetUnicornVersion()
		if err != nil {
			// Internal server error
			w.WriteHeader(http.StatusInternalServerError)
			response := generateJsonBodyForFailureCause(err.Error())
			w.Write(response)
			return
		}
		response := generateJsonBodyForUnicornVersion(version)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err := fmt.Sprintf("Invalid Operation - %s. Available operations - [LIST_SUT]", operation)
		response := generateJsonBodyForFailureCause(err)
		w.Write(response)
	}
}
