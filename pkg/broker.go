package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type BrokerImpl struct {
	Logger    lager.Logger
	Config    Config
	Instances map[string]brokerapi.GetInstanceDetailsSpec
	Bindings  map[string]brokerapi.GetBindingSpec
}

type Config struct {
	ServiceName    string
	ServicePlan    string
	BaseGUID       string
	Credentials    interface{}
	Tags           string
	ImageURL       string
	SysLogDrainURL string
	Free           bool

	FakeAsync    bool
	FakeStateful bool
}

func NewBrokerImpl(logger lager.Logger) (bkr *BrokerImpl) {
	var credentials interface{}
	json.Unmarshal([]byte(getEnvWithDefault("CREDENTIALS", "{\"port\": \"4000\"}")), &credentials)
	fmt.Printf("Credentials: %v\n", credentials)

	return &BrokerImpl{
		Logger:    logger,
		Instances: map[string]brokerapi.GetInstanceDetailsSpec{},
		Bindings:  map[string]brokerapi.GetBindingSpec{},
		Config: Config{
			BaseGUID:    getEnvWithDefault("BASE_GUID", "29140B3F-0E69-4C7E-8A35"),
			ServiceName: getEnvWithDefault("SERVICE_NAME", "some-service-name"),
			ServicePlan: getEnvWithDefault("SERVICE_PLAN_NAME", "shared"),
			Credentials: credentials,
			Tags:        getEnvWithDefault("TAGS", "shared,worlds-simplest-service-broker"),
			ImageURL:    os.Getenv("IMAGE_URL"),
			Free:        true,

			FakeAsync:    os.Getenv("FAKE_ASYNC") == "true",
			FakeStateful: os.Getenv("FAKE_STATEFUL") == "true",
		},
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}

