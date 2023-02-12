package contest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	v, err := json.Marshal(Config{
		Title:          map[string]string{"zh": "niii"},
		StartAt:        time.Now(),
		EndAt:          time.Now(),
		FrozenDuration: 200,
	})
	t.Error(string(v), err)

	config := Config{}
	err = json.Unmarshal(v, &config)
	t.Error(config, err)
}
