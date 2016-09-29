package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RunHttpServer() {
	router := gin.New()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/sortasc", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))
		//startTime := currentTimeMillis()
		arr := randInt(int32(limit))
		//r := []int32{}
		qsort(arr, 0, len(arr)-1, true)
		//		for i := 0; i < 10; i++ {
		//			r = append(r, arr[i])
		//		}
		//		endTime := currentTimeMillis()
		//		fmt.Println("Http time->", endTime, startTime, (endTime - startTime))
		jsonBytes, _ := json.Marshal(arr)
		c.String(http.StatusOK, "%s", string(jsonBytes))
	})

	router.GET("/sortdesc", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))
		arr := randInt(int32(limit))
		r := []int32{}
		qsort(arr, 0, len(arr)-1, false)
		for i := 0; i < 10; i++ {
			r = append(r, arr[i])
		}
		jsonBytes, _ := json.Marshal(r)
		c.String(http.StatusOK, "%s", string(jsonBytes))
	})

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "http hello world")
	})
	router.Run(":23456")
}
