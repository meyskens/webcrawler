package main

import "testing"

func TestMakeLinkAbsolute(t *testing.T) {
	mainURL = "https://eyskens.me/"
	if link, _ := makeLinkAbsolute("https://eyskens.me/hello", "/test"); link != "https://eyskens.me/test" {
		t.Error("Fails on /")
	}
	if link, _ := makeLinkAbsolute("https://eyskens.me/hello/world.html", "../test"); link != "https://eyskens.me/test" {
		t.Error("Fails on ../")
	}
	if link, _ := makeLinkAbsolute("https://eyskens.me/hello/world.html", "./test"); link != "https://eyskens.me/hello/test" {
		t.Error("Fails on ./")
	}
	if link, _ := makeLinkAbsolute("https://eyskens.me/hello/world.html", "~/test"); link != "https://eyskens.me/test" {
		t.Error("Fails on ~/")
	}
	if link, _ := makeLinkAbsolute("https://eyskens.me/", "//eyskens.me/css/"); link != "https://eyskens.me/css/" {
		t.Error(makeLinkAbsolute("https://eyskens.me/", "//eyskens.me/css/"))
		t.Error("Fails on //")
	}
	if _, err := makeLinkAbsolute("https://eyskens.me/hello/world.html", ""); err == nil {
		t.Error("Fails on no empty link")
	}
	if link, _ := makeLinkAbsolute("https://eyskens.me/hello/world.html", "test"); link != "https://eyskens.me/hello/test" {
		t.Error("Fails on no prefix")
	}
}

func TestGetDirectory(t *testing.T) {
	if getDirectory("https://eyskens.me/test/", 0) != "https://eyskens.me/test/" {
		t.Error("Fails on only the directory")
	}

	if getDirectory("https://eyskens.me/test/go/test.html", 0) != "https://eyskens.me/test/go/" {
		t.Error("Fails on file in subdirectory")
	}

	if getDirectory("https://eyskens.me/test.html", 0) != "https://eyskens.me/" {
		t.Error("Fails on file on root")
	}

	if getDirectory("https://eyskens.me/", 0) != "https://eyskens.me/" {
		t.Error("Fails on root")
	}

	if getDirectory("https://eyskens.me/test/", 1) != "https://eyskens.me/" {
		t.Error("Fails on only the directory with offset 1")
	}

	if getDirectory("https://eyskens.me/test/test2/test3/", 3) != "https://eyskens.me/" {
		t.Error("Fails on offset 3")
	}

	if getDirectory("https://eyskens.me/test/go/test.html", 1) != "https://eyskens.me/test/" {
		t.Error("Fails on file in subdirectory with offset 1")
	}
}

func TestIsInSameDomain(t *testing.T) {
	mainURL = "https://eyskens.me/"
	if isInSameDomain("https://gocardless.com/404") {
		t.Error("Fails on external domain")
	}
	if !isInSameDomain("https://eyskens.me/404") {
		t.Error("Fails on same domain")
	}
}

func TestAddURLToScan(t *testing.T) {
	urlsAdded = map[string]bool{}
	urlsAdded["eyskens.me/test"] = true
	addURLToScan("https://eyskens.me/test/")
	if len(urlsAdded) > 1 {
		t.Error("Doesn't ignore tailing /")
	}

	urlsAdded = map[string]bool{}
	urlsAdded["eyskens.me/test"] = true
	addURLToScan("http://eyskens.me/test")
	if len(urlsAdded) > 1 {
		t.Error("Doesn't ignore tailing scheme")
	}
}
