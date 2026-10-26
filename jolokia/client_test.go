package jolokia

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCQLStatsCassandra311(t *testing.T) {
	// Captured from Cassandra 3.11.19 with Jolokia 1.3.7 after executing
	// prepared statements. The cache size is a gauge (Value), while
	// execution and eviction counters use Count.
	body, err := ioutil.ReadFile("testdata/cql-stats-cassandra-3.11.19.json")
	require.NoError(t, err)

	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/jolokia/read/org.apache.cassandra.metrics:type=CQL,name=*", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	defer endpoint.Close()

	stats, err := Init(endpoint.URL, time.Second).CQLStats()
	require.NoError(t, err)
	assert.Equal(t, Gauge(5), stats.PreparedStatementsCount)
	assert.Equal(t, Counter(806), stats.PreparedStatementsExecuted)
	assert.Equal(t, Counter(16), stats.RegularStatementsExecuted)
	assert.Equal(t, Counter(0), stats.PreparedStatementsEvicted)
	assert.InDelta(t, 0.9805352798053528, float64(stats.PreparedStatementsRatio), 1e-12)
}
