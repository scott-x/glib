/*
* @Author: scottxiong
* @Date:   2021-06-16 20:39:13
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-16 21:15:10
*/
package glib

import (
	"net/http"
	"time"
)

const X_Session_ID = "X-Session-ID"

type Session struct {
	ID string
	Name string
	TTL time.Time//ms
}

var session map[string]interface{}

func init() {
	session = make(map[string]interface{}, 0)
}

//set sessionID to cookie
func SetSessionIDToCookie(w http.ResponseWriter, sess Session){
	expire := time.Now().Add(sess.TTL)
	cookie := http.Cookie{
		Name: name,
		Value: value, //usually sessionId
		Expires: expire,
		Path:"/",
		HttpOnly: true,
	}
	http.SetCookie(w,&cookie)
}

func NewCookie(w http.ResponseWriter, name, value string, ttl time.Duration){
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name: name,
		Value: value, //usually sessionId
		Expires: expire,
		Path:"/",
		HttpOnly: true,
	}
	http.SetCookie(w,&cookie)
}

func isCookieExpired() bool{

}