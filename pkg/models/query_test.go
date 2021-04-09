package models_test

import (
	"testing"
	"time"

	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSQuery(t *testing.T) {
	f0 := models.ExpressionConfig{
		BaseSignalField: models.BaseSignalField{
			Name: "Hello",
		},
		Expr: "x",
	}

	query := &models.SignalQuery{
		MaxDataPoints: 10,
		Signal: models.SignalConfig{
			Name: "test",
			Time: models.TimeFieldConfig{
				Period: "10s",
			},
			Fields: []models.ExpressionConfig{f0},
		},
		Interval: time.Millisecond * 100,
	}

	key0 := models.GetStreamKey(query)
	require.Equal(t, "dc9e1a2dbc6bc0d5dc891cacbad45c5fd9d7e02d4fcca13c3f400eb1f6094011", key0)

	query.Interval = time.Minute * 2
	key1 := models.GetStreamKey(query)
	require.NotEqual(t, key1, key0)

	query.Interval = time.Millisecond * 100
	key1 = models.GetStreamKey(query)
	require.Equal(t, key1, key0)
}
