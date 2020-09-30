package nodeman_test

import (
	"testing"

	"github.com/rdaniels6813/cli-manager/pkg/nodeman"
	"github.com/stretchr/testify/assert"
)

func TestGetBinsString(t *testing.T) {
	p := nodeman.NpmViewResponse{
		Bin: "./bin/test-app",
	}
	bins := p.GetBins()
	assert.Len(t, bins, 1)
	assert.Equal(t, "./bin/test-app", bins["test-app"])
}

func TestGetBinsMapInterface(t *testing.T) {
	p := nodeman.NpmViewResponse{
		Bin: map[string]interface{}{"test-app": "./bin/test-app"},
	}
	bins := p.GetBins()
	assert.Len(t, bins, 1)
	assert.Equal(t, "./bin/test-app", bins["test-app"])
}

func TestGetBinsMap(t *testing.T) {
	p := nodeman.NpmViewResponse{
		Bin: map[string]string{"test-app": "./bin/test-app"},
	}
	bins := p.GetBins()
	assert.Len(t, bins, 1)
	assert.Equal(t, "./bin/test-app", bins["test-app"])
}
