package json

import (
	"encoding/json"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/pkg/errors"
)

type intermediateJsonPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

type PatchRequest struct {
	JsonDoc   string
	JsonPath  string
	Operation string
	Value     string
}

func Patch(req PatchRequest) (string, error) {
	jObj, err := oj.ParseString(req.JsonDoc)
	if err != nil {
		return "", errors.Wrap(err, "Error during parsing json")
	}

	jPath, err := jp.ParseString(req.JsonPath)
	if err != nil {
		return "", errors.Wrap(err, "Error during parsing jpath")
	}

	locations := jPath.Locate(jObj, 0)

	jsonDoc := req.JsonDoc
	for _, l := range locations {

		resolvedJsonPath := l.String()
		jsonPointer := convertToJsonPointer(resolvedJsonPath)

		patch, err := createJsonPatch(req.Operation, jsonPointer, req.Value)
		if err != nil {
			return "", err
		}

		modified, err := patch.Apply([]byte(jsonDoc))
		if err != nil {
			return "", errors.Wrap(err, "Error during patch apply")
		}

		jsonDoc = string(modified)
	}

	return jsonDoc, nil
}

func createJsonPatch(operation string, path string, value string) (jsonpatch.Patch, error) {
	var valueObj any
	err := json.Unmarshal([]byte(value), &valueObj)
	if err != nil {
		return nil, errors.Wrap(err, "Error during unmarshaling value")
	}

	iPatch := intermediateJsonPatch{
		Op:    operation,
		Path:  path,
		Value: valueObj,
	}

	var iPatches [1]intermediateJsonPatch
	iPatches[0] = iPatch

	patchAsJson, err := json.Marshal(iPatches)
	if err != nil {
		return nil, errors.Wrap(err, "Error during creating jsonPatch (marshal)")
	}

	patch, err := jsonpatch.DecodePatch(patchAsJson)
	if err != nil {
		return nil, errors.Wrap(err, "Error during creating jsonPatch (decode)")
	}

	return patch, nil
}

func convertToJsonPointer(resolvedJsonPath string) string {
	jsonPointer := resolvedJsonPath
	//fmt.Println("JsonPath:" + jsonPointer)

	jsonPointer = strings.Replace(jsonPointer, "].", "/", -1)
	jsonPointer = strings.Replace(jsonPointer, "[", "/", -1)
	jsonPointer = strings.Replace(jsonPointer, "]", "/", -1)
	jsonPointer = strings.Replace(jsonPointer, ".", "/", -1)
	jsonPointer = strings.Replace(jsonPointer, "$", "", -1)

	//fmt.Println("JsonPointer:" + jsonPointer)
	return jsonPointer
}
