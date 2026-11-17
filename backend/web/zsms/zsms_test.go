package zsms

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBuildAliyunTemplateParam(t *testing.T) {
	raw, err := buildAliyunTemplateParam(Message{
		Code: "123456",
		Params: []TemplateParam{
			{Key: "time", Value: "5"},
		},
	})
	if err != nil {
		t.Fatalf("build params: %v", err)
	}
	var got map[string]string
	if err = json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatalf("unmarshal params: %v", err)
	}
	want := map[string]string{"code": "123456", "time": "5"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestBuildTencentParamValues(t *testing.T) {
	values := buildTencentParamValues(Message{
		Code: "123456",
		Params: []TemplateParam{
			{Key: "time", Value: "5"},
		},
	})
	want := []string{"123456", "5"}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("got %v, want %v", values, want)
	}
}
