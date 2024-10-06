package main

import "testing"

func TestRegex(t *testing.T) {
	if !re_command.MatchString(`{{- extends "layout/default" }}
{{- block body() }}`) {
		t.Error("re_command is broken")
	}

	if stripComments("awawa {* awawawawawa*} awawawa") != `awawa  awawawa` {
		t.Error("re_comment is broken")
	}
}
