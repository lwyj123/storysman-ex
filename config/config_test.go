package config

import (
	"fmt"
	"os"
	"testing"

	"code.byted.org/gopkg/logs"
)

func TestConfig(t *testing.T) {
	// test NewConfig
	conf, err := NewConfig("../conf/conf_open_dev.json")
	if err != nil {
		t.Fatal(err)
	}
	if conf.Env != "dev" {
		t.Fatal("config env invalid")
	}
	// test Init
	err = Init("../conf/conf_open_dev.json")
	if err != nil {
		t.Fatal(err)
	}
	if Instance == nil {
		t.Fatal("config init failed")
	}
	if Instance.Product() {
		t.Fatal("config env invalid")
	}
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	fmt.Println("config tests set up.")
	configFile := "../conf/conf_open_dev.json"
	err := Init(configFile)
	if err != nil {
		logs.Error("init config error: %v", err)
	}
}

func tearDown() {
	fmt.Println("config tests tear down")
}
