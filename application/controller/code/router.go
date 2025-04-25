package code

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.POST("/saveDraft", SaveDraft)                    // 请求路径是/code/saveDraft
	router.GET("/draftBox", QueryDraftBox)                  // 请求路径是/code/draftBox
	router.PUT("/updateDraftMetaData", UpdateDraftMetaData) // 请求路径是/code/updateDraftMetaData
	router.POST("/updateDraftFile", UpdateDraftFile)        // 请求路径是/code/updateDraftFile
	router.DELETE("/deleteDraft", DeleteDraft)              // 请求路径是/code/deleteDraft
	router.POST("/queryDetails", QueryDetails)              // 请求路径是/code/queryDetails
	router.POST("/queryCipInfo", QueryCipInfo)              // 请求路径是/code/queryCipInfo
	router.POST("/register", Register)                      // 请求路径是/code/register
	router.PUT("/updateCip", UpdateCip)                     // 请求路径是/code/updateCip
	router.DELETE("/deleteCip/:cpi", DeleteCIP)             // 请求路径是/code/deleteCip/:cpi
	router.POST("/queryCipByFilename", QueryCipByFilename)  // 请求路径是/code/queryCipByFilename
	router.GET("/queryByAuthor", QueryByAuthor)             // 请求路径是/code/queryByAuthor
	router.POST("/queryByDesc", QueryByDesc)                // 请求路径是/code/queryByDesc
	router.GET("/queryAllCips", QueryAllCips)               // 请求路径是/code/queryAllCips
	router.GET("/count", Count)                             // 请求路径是/code/count
	router.POST("/check", Check)                            // 请求路径是/code/check
	//router.POST("/draftDetails", QueryDraftDetails)         // 请求路径是/code/draftDetails
	//router.POST("/getCipDetails", QueryCipRentDetails)      // 请求路径是/code/getCipDetails
}
