package recsys_client

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	url2 "net/url"
	"strconv"
)

//go:generate mockgen -source=recsys_client.go -destination=mocks/mock.go
type RecSysProvider interface {
	GetRecs(id uint64, page int, pageSize int) ([]uint64, error)
}

type RecSysResponse struct {
	Ids []uint64 `json:"ids"`
}

type recsysRemote struct {
	addr string
}

func NewRecSysClient(addr string) RecSysProvider {
	return &recsysRemote{addr: addr}
}

func (r recsysRemote) GetRecs(id uint64, page int, pageSize int) ([]uint64, error) {
	request, err := http.NewRequest("GET", r.addr, nil)
	if err != nil {
		return nil, errors.Wrap(err, "recsys.recsys_client.GetRecs error")
	}
	request.URL.RawQuery = url2.Values{
		"id":        {strconv.FormatUint(id, 10)},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}.Encode()

	client := http.DefaultClient
	respGot, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "recsys.recsys_client.GetRecs error")
	}
	defer respGot.Body.Close()

	var resp RecSysResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return nil, errors.Wrap(err, "recsys.recsys_client.GetRecs error")
	}

	return resp.Ids, nil
}
