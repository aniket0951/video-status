package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	CurrentCount int   `json:"current_count"`
	PrevCount    int   `json:"prev_count"`
	TotalCount   int64 `json:"total_count"`
}

func PaginationData() Pagination {
	return Pagination{
		CurrentCount: CURRENT_COUNT,
		PrevCount:    PREV_COUNT,
		TotalCount:   int64(TOTAL_COUNT),
	}
}

type EmptyObj struct{}

func ResponseBuilder(msg string, err string, status bool, dataName string, data interface{}, enablePagination bool) map[string]interface{} {
	response := map[string]interface{}{}

	response["status"] = status
	response["message"] = msg
	response["error"] = err
	response[dataName] = data

	if enablePagination {
		paginationData := PaginationData()
		response["pagination"] = paginationData
	}

	return response
}

func BuildSuccessResponse(msg string, dataName string, data interface{}) map[string]interface{} {
	response := ResponseBuilder(msg, "", true, dataName, data, false)
	return response
}

func BuildFailedResponse(msg string, err string, dataName string) map[string]interface{} {
	response := ResponseBuilder(msg, err, false, dataName, EmptyObj{}, false)
	return response
}

func BuildResponseWithPagination(msg string, err string, status bool, dataName string, data interface{}, enablePagination bool) map[string]interface{} {
	response := ResponseBuilder(msg, err, false, dataName, data, true)
	return response
}

func RequestBodyEmptyResponse(ctx *gin.Context) {
	response := ResponseBuilder(FAILED_PROCESS, REQUIRED_PARAMS, false, DATA, EmptyObj{}, false)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)

}

func BuildUnprocessableEntityResponse(ctx *gin.Context, err error) {
	response := ResponseBuilder(FAILED_PROCESS, err.Error(), false, DATA, EmptyObj{}, false)
	ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
}
