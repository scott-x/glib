/*
* @Author: scottxiong
* @Date:   2021-06-16 20:25:00
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-22 17:40:44
 */
package glib

import (
	"github.com/gin-gonic/gin"
	"log"
)

func AllowCrossOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}

//session auth
func SessionAuthMiddleware(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := GetSession(c.Request)
		log.Printf("sess: %v", sess)
		if err != nil || len(sess.ID) == 0 || !IsUserLogin(c.Writer, c.Request, domain) {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":    2001,
				"message": "user not login",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
