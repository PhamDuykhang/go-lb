package util

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func JSONWrite(w http.ResponseWriter, httpStatus int, mgs interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)
	if mgs != nil {
		if err := json.NewEncoder(w).Encode(mgs); err != nil {
			logrus.Errorf("encode json response error:", err)
		}
	}
}
