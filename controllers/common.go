package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func paramIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func paramIsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsErrorCode(err error, errcode string) bool {

	if pgerr, ok := err.(*pgconn.PgError); ok {
		return pgerr.Code == errcode
	}
	return false
}
