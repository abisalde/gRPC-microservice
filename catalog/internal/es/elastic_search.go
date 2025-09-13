package es

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

var (
	client     *elasticsearch.TypedClient
	clientOnce sync.Once
	initErr    error
)

type ElasticClient struct {
	Client *elasticsearch.TypedClient
	Index  string
}

func Connect() (*ElasticClient, error) {
	index := "catalog"
	clientOnce.Do(func() {
		var err error

		client, err = initElastic()

		if err != nil {
			initErr = fmt.Errorf("🛑 Elasticsearch initialization failed: %w", err)
			return
		}

		if err := migrate(context.Background(), client, index); err != nil {
			initErr = fmt.Errorf("🛠️ Elasticsearch migration failed: %w", err)
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return &ElasticClient{
		Client: client,
		Index:  index,
	}, nil
}

func migrate(ctx context.Context, client *elasticsearch.TypedClient, index string) error {

	exists, err := client.Indices.Exists(index).IsSuccess(ctx)
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}

	if exists {
		log.Printf("▶️ Index %s already exists", index)
		return nil
	}

	_, err = client.Indices.Create(index).
		Mappings(&types.TypeMapping{
			Properties: map[string]types.Property{
				"id": types.KeywordProperty{},
				"name": map[string]interface{}{
					"type": "text",
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"price": map[string]interface{}{
					"type": "float",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"updated_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	log.Printf("✅ Created index: %s", index)

	return nil
}

func formatURL() string {
	return "http://catalog_db:9200"
}

func initElastic() (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			formatURL(),
		},
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
			MaxIdleConnsPerHost:   10,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("📡 Failed to create ES client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = es.Ping().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("⚙️ Elasticsearch ping failed: %w", err)
	}

	info, err := es.Info().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping ES: %w", err)
	}
	log.Printf("✅ Connected to Elasticsearch: %s", info.ClusterName)

	return es, nil
}

func (e *ElasticClient) Close() error {
	client = nil
	clientOnce = sync.Once{}
	return nil
}

func (e *ElasticClient) HealthCheck(ctx context.Context) error {
	if e.Client == nil {
		return fmt.Errorf("🎛️ Elasticsearch client not initialized")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}

	_, err := e.Client.Ping().Do(ctx)
	if err != nil {
		return fmt.Errorf("🕹️ Elasticsearch ping failed: %w", err)
	}

	health, err := e.Client.Cluster.Health().Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check ES health: %w", err)
	}
	if health.Status.String() == "red" {
		return fmt.Errorf("ES cluster is unhealthy (red status)")
	}
	return nil
}
