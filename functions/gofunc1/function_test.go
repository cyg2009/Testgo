package main

import (
	"testing"
	"plugin"
    "time"
)

func TestInvokeFunction(t *testing.T) {
	p, err := plugin.Open("code/gofunc1.so.1.0.0")
	if err != nil {
		t.Errorf("Fail to load the plugin: %s", err)
		return
	}

	greetSymbol, err := p.Lookup("Handler")
	if err != nil {
		t.Errorf("Fail to find the handler: %s", err)
		return
	}

	trigger := greetSymbol.(func(string) string)
	
	evt := "{\"name\":\"king\", \"timestamp\":\""  +  time.Now().String() +  "\"}"
	ret := trigger(evt)

	t.Log(ret)	
}
