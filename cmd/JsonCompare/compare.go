package JsonCompare

import (
	"encoding/json"
	"log"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ErrorType int

const (
	KeyError ErrorType = iota
	ValueError
	TypeError
)

var errors = map[ErrorType]string{
	0: "KeyError",
	1: "ValueError",
	2: "TypeError",
}

type Output struct {
	TotalBytes    float64
	MismatchBytes float64
	Score         float32
	File1Mismatch []SnapShot
	File2Mismatch []SnapShot
}

type SnapShot struct {
	Path        []string
	MisMatch    string
	Description string
	Error       string
}

func (o *Output) incrScore(val int) {
	o.MismatchBytes += float64(val)
	o.Score = float32(1 - (o.MismatchBytes / o.TotalBytes))
}

func getJsonSize(object interface{}) int {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`"[^"]*"|\d+(?:\.\d+)?|\w+`)

	str := re.FindAllString(string(jsonBytes), -1)

	return len(strings.Join(str, ""))
}

var mismatches []SnapShot

var output = Output{
	Score: 1.0,
}

func CompareFiles(path1 string, path2 string) Output {
	abs, err := filepath.Abs(path1)
	if err != nil {
		log.Fatal(err)
	}
	file1, err := os.Open(abs)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	abs2, err := filepath.Abs(path2)
	if err != nil {
		log.Fatal(err)
	}
	file2, err := os.Open(abs2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	err, jsonObject1 := parseJson(*file1)
	if err != nil {
		log.Println("File1 from path ", path1, " was not a valid json object")
		log.Fatal(err)
	}
	err, jsonObject2 := parseJson(*file2)
	if err != nil {
		log.Println("File2 from path ", path2, " was not a valid json object")
		log.Fatal(err)
	}

	mismatches = make([]SnapShot, 0)

	output.TotalBytes = float64(getJsonSize(jsonObject1) + getJsonSize(jsonObject2))

	currentPath := []string{"root"}

	check(jsonObject1, jsonObject2, currentPath)

	output.File1Mismatch = mismatches
	mismatches = make([]SnapShot, 0)

	check(jsonObject2, jsonObject1, currentPath)

	output.File2Mismatch = mismatches

	return output

}

/*
Object1 is compared against object2
*/
func check(object1 any, object2 any, currentPath []string) {
	if object1 == nil {
		return
	}
	if object2 == nil {
		typeError(object1, object2, currentPath)
		return
	}
	if reflect.TypeOf(object1).Kind() == reflect.Map {
		mapCheck(object1, object2, currentPath)
	} else if reflect.TypeOf(object1).Kind() == reflect.Slice {
		sliceCheck(object1, object2, currentPath)
	} else if reflect.TypeOf(object1).Kind() == reflect.Float64 {
		floatCheck(object1, object2, currentPath)
	} else if reflect.TypeOf(object1).Kind() == reflect.String {
		stringCheck(object1, object2, currentPath)
	} else if reflect.TypeOf(object2).Kind() == reflect.Bool {
		boolCheck(object1, object2, currentPath)
	}
}

func floatCheck(object1 any, object2 any, currentPath []string) {
	if reflect.TypeOf(object2).Kind() != reflect.Float64 {
		typeError(object1, object2, currentPath)
		return
	}
	intAreEqual := object1.(float64) == object2.(float64)
	if !intAreEqual {
		mismatches = append(mismatches, SnapShot{
			Path:        append([]string(nil), currentPath...),
			MisMatch:    shorten(strconv.FormatFloat(object1.(float64), 'f', -1, 64), 15) + ":" + shorten(strconv.FormatFloat(object2.(float64), 'f', -1, 64), 15),
			Description: "Floats are not equal",
			Error:       errors[ValueError],
		})
		output.incrScore(len(strconv.FormatFloat(object1.(float64), 'f', -1, 64)))
	}
}

func stringCheck(object1 any, object2 any, currentPath []string) {
	if reflect.TypeOf(object2).Kind() != reflect.String {
		typeError(object1, object2, currentPath)
		return
	}
	stringsAreEqual := object1.(string) == object2.(string)
	if !stringsAreEqual {
		mismatches = append(mismatches, SnapShot{
			Path:        append([]string(nil), currentPath...),
			MisMatch:    shorten(object1.(string), 15) + ":" + shorten(object2.(string), 15),
			Description: "String Values are not equal",
			Error:       errors[ValueError],
		})
		output.incrScore(len(object1.(string)))
	}
}

func sliceCheck(object1 any, object2 any, currentPath []string) {
	if reflect.TypeOf(object2).Kind() != reflect.Slice {
		typeError(object1, object2, currentPath)
		return
	}
	for index, item1 := range object1.([]interface{}) {
		if index >= len(object2.([]interface{})) {
			if item1Bytes, err := json.Marshal(item1); err != nil {
				log.Fatal(err)
			} else {
				mismatches = append(mismatches, SnapShot{
					Path:        append([]string(nil), currentPath...),
					MisMatch:    shorten(string(item1Bytes), 15),
					Description: "Items in list is after last entry in other file",
					Error:       errors[ValueError],
				})
				output.incrScore(len(item1Bytes))
			}
			continue
		}
		check(item1, object2.([]any)[index], append(currentPath, strconv.Itoa(index)))
	}
}

func boolCheck(object1 any, object2 any, currentPath []string) {
	if reflect.TypeOf(object2).Kind() != reflect.Bool {
		typeError(object1, object2, currentPath)
		return
	}
	if object1.(bool) != object2.(bool) {
		mismatches = append(mismatches, SnapShot{
			Path:        append([]string(nil), currentPath...),
			MisMatch:    boolToString(object1.(bool)) + "  :  " + boolToString(object2.(bool)),
			Description: "Bools are not equal",
			Error:       errors[ValueError],
		})
		output.incrScore(len(boolToString(object1.(bool))))
	}
}

func mapCheck(object1 any, object2 any, currentPath []string) {
	if reflect.TypeOf(object2).Kind() != reflect.Map {
		typeError(object1, object2, currentPath)
		return
	}
	for key := range maps.Keys(object1.(map[string]interface{})) {
		_, keyInMap := object2.(map[string]any)[key]
		if !keyInMap {
			mismatches = append(mismatches, SnapShot{
				Path:        append([]string(nil), currentPath...),
				MisMatch:    key,
				Description: "Key Not in Other File",
				Error:       errors[KeyError],
			})
			output.incrScore(getJsonSize(object1.(map[string]interface{})[key]) + getJsonSize(key))
		} else {
			check(object1.(map[string]any)[key], object2.(map[string]any)[key], append(currentPath, key))
		}
	}
}

func typeError(object1 any, object2 any, currentPath []string) {

	jsonString1, err := json.Marshal(object1)
	if err != nil {
		log.Fatal(err)
	}
	jsonString2, err := json.Marshal(object2)
	if err != nil {
		log.Fatal(err)
	}
	mismatches = append(mismatches, SnapShot{
		Path:        append([]string(nil), currentPath...),
		MisMatch:    shorten(string(jsonString1), 15) + ":" + shorten(string(jsonString2), 15),
		Description: "Value Type Mismatch",
		Error:       errors[TypeError],
	})
}
func shorten(s string, max int) string {
	if len(s) > max {
		return s[:max] + "...  "
	}
	return s + "   "
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func parseJson(file os.File) (error, any) {
	var jsonRes interface{}
	err := json.NewDecoder(&file).Decode(&jsonRes)
	return err, jsonRes
}
