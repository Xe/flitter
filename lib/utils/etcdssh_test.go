package utils

import "testing"

func TestSSHFingerPrintGen(t *testing.T) {
	key := "AAAAB3NzaC1yc2EAAAADAQABAAABAQDnT2qR9ETxfMTEV71SKstcluusH66Kf1+D987e177x1Ku+ejPShzw5lRN6LYe+gv1zlcyQNErt+zhZKyqiJhPffW1Tn47ub70sOtITcGkWc3xme/bFlR+1r4DrlAiiMmeikeCFFyR/F8BSDR4o3XCr8S6kvwukDnNcfepL+uWFD5Mrdl2FZ2TpmQz/XwLpioPb6BBZsw6SYiWRXO5UbpCleP4CGvqKZM5gJIB4JYYIKmSxGN2UNTdrWiHnhPMYd+0+1a6MhmlY/MeenRzW+m6Okyc3jkFW0U2nVtPSvw+LcIcXj1Np0ZfF9hd/3kHngFurqTrKqOS9ByaKQVKGfFG1"

	fp := GetFingerprint(key)

	if fp == "AAAAB3NzaC1yc2EAAAADAQABAAABAQDnT2qR9ETxfMTEV71SKstcluusH66Kf1+D987e177x1Ku+ejPShzw5lRN6LYe+gv1zlcyQNErt+zhZKyqiJhPffW1Tn47ub70sOtITcGkWc3xme/bFlR+1r4DrlAiiMmeikeCFFyR/F8BSDR4o3XCr8S6kvwukDnNcfepL+uWFD5Mrdl2FZ2TpmQz/XwLpioPb6BBZsw6SYiWRXO5UbpCleP4CGvqKZM5gJIB4JYYIKmSxGN2UNTdrWiHnhPMYd+0+1a6MhmlY/MeenRzW+m6Okyc3jkFW0U2nVtPSvw+LcIcXj1Np0ZfF9hd/3kHngFurqTrKqOS9ByaKQVKGfFG1" {
		t.Fatalf("GetFingerprint is a no-op: %s", fp)
	}
}

func TestAddColons(t *testing.T) {
	sum := "d3b07384d113edec49eaa6238ad5ff00"
	colons := FPAddColons(sum)

	if sum == colons {
		t.Fatalf("AddColons is a no-op: %s %s", colons, sum)
	}

	t.Log(colons)
}
