package remote

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

//go:generate mockgen -source=remote.go -destination=mocks/mock.go
type RecSysProvider interface {
	GetRecs(id uint64) ([]uint64, error)
}

type Response struct {
	Ids []uint64 `json:"ids"`
}

type recsysRemote struct {
	addr string
}

func New(addr string) RecSysProvider {
	return &recsysRemote{addr: addr}
}

func (r recsysRemote) GetRecs(id uint64) ([]uint64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", r.addr, id))
	if err != nil {
		return nil, errors.Wrap(err, "recsys.remote.GetRecs error")
	}
	defer resp.Body.Close()

	var respParsed Response
	if err := json.NewDecoder(resp.Body).Decode(&respParsed); err != nil {
		return nil, errors.Wrap(err, "recsys.remote.GetRecs error")
	}

	return respParsed.Ids, nil
}
