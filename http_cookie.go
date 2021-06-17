/*
* @Author: scottxiong
* @Date:   2021-06-16 20:39:13
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-17 20:05:06
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
	TTL time.Duration
}

var session map[string]interface{}

func init() {
	session = make(map[string]interface{}, 0)
}

func NewDefaultSession() {

}

//generate session
func NewSession(value string, ttl time.Duration) Session {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	return Session{
		ID: base64.URLEncoding.EncodeToString(b),
		Name: value,
		TTL: ttl,
	}
}

//set sessionID to cookie
func SetSessionIDToCookie(w http.ResponseWriter, sess Session){
	expire := time.Now().Add(sess.TTL)
	cookie := http.Cookie{
		Name: X_Session_ID,
		Value: sess.ID, //usually sessionId
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