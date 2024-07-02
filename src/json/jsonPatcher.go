package json

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
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
	pointers, err := ConvertPathToPointers([]byte(req.JsonDoc), req.JsonPath)
	if err != nil {
		return "", errors.Wrap(err, "Error during parsing json")
	}

	jsonDoc := req.JsonDoc
	for _, pointer := range pointers {
		patch, err := createJsonPatch(req.Operation, pointer, req.Value)
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
