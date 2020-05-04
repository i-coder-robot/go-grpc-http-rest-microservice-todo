package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func AddLogger(logger *zap.Logger,h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Header.Get("X-Liveness-Probe") == "Healthz"{
			h.ServeHTTP(w,r)
			return
		}
		id:=GetReqID(ctx)
		var scheme string
		if r.TLS !=nil{
			scheme = "https"
		}else{
			scheme = "http"
		}
		proto := r.Proto
		method := r.Method
		remoteAddr := r.RemoteAddr
		userAgent := r.UserAgent()
		uri := strings.Join([]string{scheme,"://",r.Host,r.RequestURI},"")

		logger.Debug("request started",
			zap.String("request-id",id),
			zap.String("http-scheme",scheme),
			zap.String("http-proto",proto),
			zap.String("http-method",method),
			zap.String("remote-addr",remoteAddr),
			zap.String("user-agent",userAgent),
			zap.String("uri",uri),
			)

		t1 := time.Now()
		h.ServeHTTP(w,r)

		logger.Debug("request completed",
			zap.String("request-id",id),
			zap.String("http-scheme",scheme),
			zap.String("http-proto",proto),
			zap.String("http-method",method),
			zap.String("remote-addr",remoteAddr),
			zap.String("user-agent",userAgent),
			zap.String("uri",uri),
			zap.Float64("elapsed-ms",float64(time.Since(t1).Nanoseconds())/1000000.0),
			)
	})
}