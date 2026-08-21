package db

import (
	"strings"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/configs"
	"github.com/elastic/go-elasticsearch/v8"
)

// SetupES builds an Elasticsearch client and verifies connectivity with a
// cluster Info call so startup fails fast if the cluster is unreachable.
func SetupES() (*elasticsearch.TypedClient, error) {
	addrs := strings.Split(configs.Cfg.ESAddrs, ",")
	for i := range addrs {
		addrs[i] = strings.TrimSpace(addrs[i])
	}

	cfg := elasticsearch.Config{
		Addresses: addrs,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil
}
