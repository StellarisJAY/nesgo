package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func (s *AdminService) HandleUploadGame(ctx http.Context) error {
	file, header, err := ctx.Request().FormFile("data")
	if err != nil {
		return err
	}
	data := make([]byte, header.Size)
	n, err := file.Read(data)
	if err != nil {
		return err
	}
	fmt.Printf("file: %s, size: %d\n", header.Filename, n)
	err = s.gf.UploadGame(context.Background(), header.Filename, data[:n])
	return err
}
