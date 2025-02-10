package statsig_fastly_adapter

import (
	"net/http"
	"regexp"

	"github.com/fastly/compute-sdk-go/kvstore"
)

type FastlyDataAdapter struct {
	fastlyKey                string
	storeName                string
	configSpecsKey           string
	supportConfigSpecPolling bool
	httpClient               *http.Client
}

func NewFastlyDataAdapter(fastlyKey, storeName, configSpecsKey string) *FastlyDataAdapter {
	client := &http.Client{}

	return &FastlyDataAdapter{
		fastlyKey:      fastlyKey,
		storeName:      storeName,
		configSpecsKey: configSpecsKey,
		httpClient:     client,
	}
}

func (f *FastlyDataAdapter) getData() (string, error) {
	kv, err := kvstore.Open(f.storeName)
	if err != nil {
		return "", err
	}

	value, err := kv.Lookup(f.configSpecsKey)
	if err != nil {
		return "", err
	}

	return value.String(), err
}

func (f *FastlyDataAdapter) Initialize() {
	data, err := f.getData()
	if err != nil {
		return
	}

	if data != "" {
		f.supportConfigSpecPolling = true
	}
}

func (f *FastlyDataAdapter) Get(key string) string {
	if !f.isConfigSpecKey(key) {
		return ""
	}

	data, err := f.getData()
	if err != nil || data == "" {
		return ""
	}

	return data
}

func (f *FastlyDataAdapter) Set(key string, value string) {
	// no-op. Statsig's Edge Config integration keeps config specs synced through Statsig's service
}

func (f *FastlyDataAdapter) Shutdown() {
	// no-op
}

func (f *FastlyDataAdapter) ShouldBeUsedForQueryingUpdates(key string) bool {
	if f.isConfigSpecKey(key) {
		return f.supportConfigSpecPolling
	}
	return false
}

func (f *FastlyDataAdapter) isConfigSpecKey(key string) bool {
	v2CacheKeyPattern := `^statsig\|/v[12]/download_config_specs\|.+\|.+`
	regex, _ := regexp.Compile(v2CacheKeyPattern)
	return key == "statsig.cache" || regex.MatchString(key)
}
