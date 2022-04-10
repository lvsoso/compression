package main

import (
	"errors"
	"go-zip/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestPack(t *testing.T) {
	err := utils.Pack("../models", "../models.zip")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnPack(t *testing.T) {
	dst, _ := filepath.Abs("/tmp/models")
	_, err := os.Stat(dst)
	if err != nil && errors.Is(err, &os.PathError{}) {
		t.Fatal(err)
	}
	src, _ := filepath.Abs("../models.zip")

	err = utils.Unpack(dst, src)
	if err != nil {
		t.Fatal(err)
	}
}

// i5-9600KF@3.70GHz
// BestSpeed 32.707s
// BestCompression ?
// DefaultCompression 158.680s
func TestKPack(t *testing.T) {
	src, _ := filepath.Abs("../models")
	dst, _ := filepath.Abs("../models.zip")

	t.Log(src)
	t.Log(dst)

	err := utils.KPack(src, dst)
	if err != nil {
		t.Fatal(err)
	}
}

// i5-9600KF@3.70GHz
// DefaultCompression 29.641s
func TestKUnPack(t *testing.T) {
	dst, _ := filepath.Abs("/tmp/models")
	_, err := os.Stat(dst)
	if err != nil && errors.Is(err, &os.PathError{}) {
		t.Fatal(err)
	}
	src, _ := filepath.Abs("../models.zip")

	err = utils.KUnpack(dst, src)
	if err != nil {
		t.Fatal(err)
	}
}

// 1.367s
func TestKZipToZip(t *testing.T) {
	src, _ := filepath.Abs("../models.zip.example")
	dst, _ := filepath.Abs("../models.zip.z2z")

	t.Log(src)
	t.Log(dst)
	err := utils.KZipToZip(src, dst, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestKZipToZipIgnore(t *testing.T) {
	src, _ := filepath.Abs("../models.zip.example")
	dst, _ := filepath.Abs("../models.zip.z2z.ignore")

	t.Log(src)
	t.Log(dst)

	ignoreMap := map[string]bool{
		"zfnet512/init_net.pb": true,
		"vgg19/init_net.pb":    true,
	}
	err := utils.KZipToZip(src, dst, ignoreMap)
	if err != nil {
		t.Fatal(err)
	}
}
