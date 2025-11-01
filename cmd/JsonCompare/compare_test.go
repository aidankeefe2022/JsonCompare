package JsonCompare

import (
	"testing"
)

func TestJsonIntCompareEqual(t *testing.T) {
	int1 := float64(1)
	int2 := float64(1)
	var currentPath []string
	check(int1, int2, currentPath)
	if len(mismatches) > 0 {
		t.Error("Expected no mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestJsonIntCompareNotEqual(t *testing.T) {
	int1 := float64(2)
	int2 := float64(1)
	var currentPath []string
	check(int1, int2, currentPath)
	if len(mismatches) != 1 {
		t.Error("Expected mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestStringsAreEqual(t *testing.T) {
	string1 := "Hello World"
	string2 := "Hello World"
	check(string1, string2, []string{})
	if len(mismatches) > 0 {
		t.Error("Expected no mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestStringsAreNotEqual(t *testing.T) {
	string1 := "Hello World"
	string2 := "Hello World2"
	check(string1, string2, []string{})
	if len(mismatches) != 1 {
		t.Error("Expected mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestMapAreEqual(t *testing.T) {
	map1 := map[string]interface{}{
		"foo": "bar",
	}
	map2 := map[string]interface{}{
		"foo": "bar",
	}
	check(map1, map2, []string{})
	if len(mismatches) > 0 {
		t.Error("Expected no mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestMapAreNotEqual(t *testing.T) {
	map1 := map[string]interface{}{
		"hello": "bar",
	}
	map2 := map[string]interface{}{
		"foo": "bar",
	}
	check(map1, map2, []string{})
	if len(mismatches) != 1 {
		t.Error("Expected mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestNestedMapsAreEqual(t *testing.T) {
	map1 := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}
	map2 := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}
	check(map1, map2, []string{})
	if len(mismatches) > 0 {
		t.Error("Expected no mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestNestedMapsAreNotEqual(t *testing.T) {
	map1 := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}
	map2 := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "bart",
		},
	}
	check(map1, map2, []string{})
	if len(mismatches) != 1 {
		t.Error("Expected mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestSlicesAreEqual(t *testing.T) {
	slice1 := []interface{}{"foo", "bar"}
	slice2 := []interface{}{"foo", "bar"}
	check(slice1, slice2, []string{})
	if len(mismatches) > 0 {
		t.Error("Expected no mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestSlicesAreNotEqual(t *testing.T) {
	slice1 := []interface{}{"foo", "barf"}
	slice2 := []interface{}{"foo", "bar"}
	check(slice1, slice2, []string{})
	if len(mismatches) != 1 {
		t.Error("Expected mismatches")
	}
	mismatches = make([]SnapShot, 0)
}

func TestCompareFiles(t *testing.T) {
	output = Output{
		Score: 1.0,
	}
	out := CompareFiles("/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/test1.json", "/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/test2.json")
	if out.Score != 0 {
		t.Error("Expected 0.0 score")
	}
	mismatches = make([]SnapShot, 0)

}

func TestCompareFilesSlices(t *testing.T) {
	output = Output{
		Score: 1.0,
	}
	out := CompareFiles("/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/test1Slice.json", "/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/test2Slice.json")
	if out.Score != 1.0 {
		t.Error("Expected 1.0 score")
	}
	mismatches = make([]SnapShot, 0)
}

func TestCompareFilesHard(t *testing.T) {
	output = Output{
		Score: 1.0,
	}
	out := CompareFiles("/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_2.json", "/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_3.json")
	if out.Score == 1.0 {
		t.Error("Did not expect 1.0 score")
	}
	mismatches = make([]SnapShot, 0)
}

func TestCompareFilesHardEqual(t *testing.T) {
	output = Output{
		Score: 1.0,
	}
	out := CompareFiles("/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_2.json", "/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_1.json")
	if out.Score != 1.0 {
		t.Error("Did not expect 1.0 score")
	}
	mismatches = make([]SnapShot, 0)
}

func TestCompareFiles34(t *testing.T) {
	output = Output{
		Score: 1.0,
	}
	out := CompareFiles("/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_4.json", "/home/aidankeefe/Work_Techlink/JsonCompare/testFiles/example_3.json")
	if out.Score == 1.0 {
		t.Error("Did not expect 1.0 score")
	}
	mismatches = make([]SnapShot, 0)
}
