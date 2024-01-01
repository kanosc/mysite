package endpoints

import (
	"context"
	"errors"

	"github.com/kanosc/mysite/server/service"

	"github.com/go-kit/kit/endpoint"
)

const (
	FileListRequestKey = "dirname"
)

type FileListRequest struct {
	DirName string `json:"DirName"`
}

type FileListResponse struct {
	FileList []service.FileInfo `json:"FileList"`
	Err      error              `json:"Err"`
}

func MakeFileListEndpoint(svc service.FileAccessor) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		dirname, ok := request.(string)
		if !ok {
			return nil, errors.New("Type Error")
		}

		v, err := svc.GetFileList(dirname)
		if err != nil {
			return FileListResponse{Err: errors.New("Dir name doesn't exsist")}, nil
		}
		return FileListResponse{FileList: v}, nil
	}
}
