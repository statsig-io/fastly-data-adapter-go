package statsig_fastly_adapter

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type FastlyDataAdapter struct {
	fastlyKey                string
	storeID                  string
	configSpecsKey           string
	supportConfigSpecPolling bool
	httpClient               *http.Client
}

func NewFastlyDataAdapter(fastlyKey, storeID, configSpecsKey string) *FastlyDataAdapter {
	client := &http.Client{}

	return &FastlyDataAdapter{
		fastlyKey:      fastlyKey,
		storeID:        storeID,
		configSpecsKey: configSpecsKey,
		httpClient:     client,
	}
}

func (f *FastlyDataAdapter) getData() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.fastly.com/resources/stores/config/%s/item/%s", f.storeID, f.configSpecsKey), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Fastly-Key", f.fastlyKey)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return string(body), err
}

func (f *FastlyDataAdapter) Initialize() error {
	data, err := f.getData()
	if err != nil {
		return err
	}

	if data != "" {
		f.supportConfigSpecPolling = true
	}
	return nil
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

func (f *FastlyDataAdapter) Set(key string, value string) error {
	// no-op. Statsig's Edge Config integration keeps config specs synced through Statsig's service
	return nil
}

func (f *FastlyDataAdapter) Shutdown() error {
	// no-op
	return nil
}

func (f *FastlyDataAdapter) SupportsPollingUpdatesFor(key string) bool {
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
