package main

import (
	"fmt"
	"gotest/session/session"
	"net/http"
	"strconv"
	"time"

	_ "gotest/session/memory"
)

var globalSessions *session.Manager

func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	fmt.Println(ct)
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	w.Write([]byte("count:" + strconv.Itoa(sess.Get("countnum").(int))))
}

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func main() {
	http.HandleFunc("/count", count)
	http.ListenAndServe(":3000", nil)
}
