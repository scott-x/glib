/*
* @Author: scottxiong
* @Date:   2021-06-16 20:25:00
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-16 20:39:02
*/
package glib

import (
	"github.com/gin-gonic/gin"
)

func AllowCrossOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
        w := c.Writer
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}