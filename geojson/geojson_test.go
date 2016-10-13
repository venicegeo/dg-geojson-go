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
	"fmt"
	"log"
	"testing"
)

var inputFiles = [...]string{
	"test/point.geojson",
	"test/linestring.geojson",
	"test/polygon.geojson",
	"test/polygon-dateline.geojson",
	"test/polygon-hole.geojson",
	"test/multipoint.geojson",
	"test/multilinestring.geojson",
	"test/multipolygon.geojson",
	"test/geometrycollection.geojson",
	"test/sample.geojson",
	"test/boundingbox.geojson"}

func testProcess(filename string) {
	var (
		gj    interface{}
		err   error
		bytes []byte
	)
	if gj, err = ParseFile(filename); err != nil {
		log.Panicf("Parse error: %v\n", err)
	}
	fmt.Printf("%T: %#v\n", gj, gj)

	if bytes, err = Write(gj); err != nil {
		log.Panicf("Write error: %v\n", err)
	}
	fmt.Printf("%v\n", string(bytes))
}

// TestGeoJSON tests GeoJSON readers
func TestGeoJSON(t *testing.T) {
	for _, fileName := range inputFiles {
		testProcess(fileName)
	}
}

func TestNullInputs(t *testing.T) {
	bb, _ := NewBoundingBox(nil)
	if "" != bb.String() {
		fmt.Printf(bb.String())
		t.Error("Couldn't handle nil bounding box")
	}
	point := bb.Centroid()
	if point != nil {
		t.Error("Expected a nil Centroid for an empty bounding box")
	}
	fc := NewFeatureCollection(nil)
	if fc.String() != `{"type":"FeatureCollection","features":[]}` {
		t.Errorf("Received %v for empty Feature Collection.", fc.String())
	}
	f := NewFeature(nil, nil, nil)
	if f.String() != `{"type":"Feature","geometry":null}` {
		t.Errorf("Received %v for an empty Feature.", f.String())
	}
	fc.Features = append(fc.Features, f)
	if fc.String() != `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":null}]}` {
		t.Errorf("Received %v for a feature collection with a single empty feature", fc.String())
	}
}
