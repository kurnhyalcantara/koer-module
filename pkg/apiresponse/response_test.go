package apiresponse

import (
	"testing"
)

func TestOK(t *testing.T) {
	r := OK("data", "success")
	if !r.Success {
		t.Error("expected success=true")
	}
	if r.Error != nil {
		t.Error("expected no error")
	}
}

func TestFail(t *testing.T) {
	r := Fail("ERR_001", "something went wrong", nil)
	if r.Success {
		t.Error("expected success=false")
	}
	if r.Error == nil || r.Error.Code != "ERR_001" {
		t.Error("expected error code ERR_001")
	}
}

func TestOKWithMeta(t *testing.T) {
	r := OKWithMeta([]int{1, 2, 3}, "ok", Meta{Page: 1, PerPage: 10, Total: 3, TotalPages: 1})
	if !r.Success || r.Meta == nil {
		t.Error("expected success with meta")
	}
}

func TestNotFound(t *testing.T) {
	r := NotFound("resource not found")
	if r.Error.Code != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %s", r.Error.Code)
	}
}
