/*
Copyright 2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package geojson

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// FEATURE is a GeoJSON Feature
const FEATURE = "Feature"

// The Feature object represents an array of features
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   interface{}            `json:"geometry"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	ID         string                 `json:"id,omitempty"`
	Bbox       BoundingBox            `json:"bbox,omitempty"`
}

// FeatureFromBytes constructs a Feature from a GeoJSON byte array
// and returns its pointer
func FeatureFromBytes(bytes []byte) (*Feature, error) {
	var result Feature
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	result.ResolveGeometry()
	return &result, nil
}

// ForceBbox returns a bounding box, creating one by brute force if needed
func (feature *Feature) ForceBbox() BoundingBox {
	if len(feature.Bbox) > 0 {
		return feature.Bbox
	}
	if bboxIfc, ok := feature.Geometry.(BoundingBoxIfc); ok {
		return bboxIfc.ForceBbox()
	}

	return BoundingBox{}
}

// String returns the string representation
func (feature *Feature) String() string {
	var result string
	if bytes, err := json.Marshal(feature); err == nil {
		result = string(bytes)
	} else {
		result = err.Error()
	}
	return result
}

// NewFeature is the normal factory method for a feature
func NewFeature(geometry interface{}, id string, properties map[string]interface{}) *Feature {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	return &Feature{Type: FEATURE, Geometry: geometry, Properties: properties, ID: id}
}

// ResolveGeometry reconstructs a Feature's geometries
// since unmarshaled objects come back as interfaces, not real geometries,
func (feature *Feature) ResolveGeometry() {
	if feature.Geometry != nil {
		if geometry, ok := feature.Geometry.(map[string]interface{}); ok {
			feature.Geometry = NewGeometry(geometry)
		}
	}
}

func stringify(input interface{}) string {
	var result string
	switch itype := input.(type) {
	case string:
		result = itype
	case fmt.Stringer:
		result = itype.String()
	case int:
		result = stringify(int64(itype))
	case float32:
		result = stringify(float64(itype))
	case int64:
		result = strconv.FormatInt(itype, 10)
	case float64:
		result = strconv.FormatFloat(itype, 'f', -1, 64)
	}
	return result
}

// PropertyString returns the string value of the property if it exists
// and is a string, or the empty string otherwise
func (feature *Feature) PropertyString(propertyName string) string {
	var result string
	if property, ok := feature.Properties[propertyName]; ok {
		return stringify(property)
	}
	return result
}

func intify(input interface{}) int {
	var result int

	switch itype := input.(type) {
	case int:
		result = itype
	case string:
		tempInt64, _ := strconv.ParseInt(itype, 10, 64)
		result = int(tempInt64)
	case fmt.Stringer:
		result = intify(itype.String())
	case int64:
		result = intify(int(itype))
	case float32:
		result = intify(int(itype))
	case float64:
		result = intify(int(itype))
	}
	return result
}

// PropertyInt returns the integer value of the property if it exists
// or 0 otherwise
func (feature *Feature) PropertyInt(propertyName string) int {
	var result int
	if property, ok := feature.Properties[propertyName]; ok {
		result = intify(property)
	}
	return result
}

// PropertyStringSlice returns the string slice value of the property if it exists
// or an empty slice otherwise
func (feature *Feature) PropertyStringSlice(propertyName string) []string {
	var result []string
	if property, ok := feature.Properties[propertyName]; ok {
		switch ptype := property.(type) {
		case []string:
			result = ptype
		case []interface{}:
			for _, curr := range ptype {
				if currString, ok := curr.(string); ok {
					result = append(result, currString)
				}
			}
		}
	}
	return result
}

func floatify(input interface{}) float64 {
	var result = math.NaN()

	switch itype := input.(type) {
	case float64:
		result = itype
	case string:
		if parsedFloat, err := strconv.ParseFloat(itype, 64); err == nil {
			result = parsedFloat
		}
	case fmt.Stringer:
		result = floatify(itype.String())
	case int64:
		result = floatify(float64(itype))
	case float32:
		result = floatify(float64(itype))
	case int:
		result = floatify(float64(itype))
	}
	return result
}

// PropertyFloat returns the floating point value of the property if it exists
// and is a float or is a parseable string, or math.NaN() otherwise
func (feature *Feature) PropertyFloat(propertyName string) float64 {
	var result = math.NaN()
	if property, ok := feature.Properties[propertyName]; ok {
		result = floatify(property)
	}
	return result
}

// FeatureFromMap constructs a Feature from a map
// and returns its pointer
func FeatureFromMap(input map[string]interface{}) *Feature {
	var result Feature
	result.Type = input["type"].(string)
	result.Properties = input["properties"].(map[string]interface{})
	result.Geometry = input["geometry"]
	result.ResolveGeometry()
	if bboxIfc, ok := input["bbox"]; ok {
		result.Bbox, _ = NewBoundingBox(bboxIfc)
	}
	if id, ok := input["id"]; ok {
		switch idtype := id.(type) {
		case string:
			result.ID = idtype
		case int:
			result.ID = strconv.FormatInt(int64(idtype), 10)
		}
	}
	return &result
}
