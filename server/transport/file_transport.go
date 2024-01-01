package transport

import (
	"context"
	"errors"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/kanosc/mysite/server/endpoints"
)

func DecodeFileListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	values := r.URL.Query()
	fmt.Printf("request params [%v]\n", values)
	dirName := values.Get(endpoints.FileListRequestKey)
	if dirName == "" {
		return nil, errors.New("No dir name found in request")
	}
	return dirName, nil
}

func EncodeFileListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp, ok := response.(endpoints.FileListResponse)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("Server type encode error")
	}

	if rsp.Err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return rsp.Err
	}

	return json.NewEncoder(w).Encode(response)
}
