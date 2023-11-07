// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2022, SCANOSS
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package purlutils

import (
	"github.com/package-url/packageurl-go"
	"reflect"
	"testing"
)

// Help with test details can be found here: https://go.dev/doc/code

func TestPurlFromString(t *testing.T) {

	w, _ := packageurl.FromString("pkg:maven/io.prestosql/presto-main@v1.0")
	w2, _ := packageurl.FromString("pkg:npm/%40babel/core")
	tests := []struct {
		name    string
		input   string
		want    packageurl.PackageURL
		wantErr bool
	}{
		{
			name:  "Purl from String",
			input: "pkg:maven/io.prestosql/presto-main@v1.0",
			want:  w,
		},
		{
			name:  "Purl from String",
			input: "pkg:npm/%40babel/core",
			want:  w2,
		},
		{
			name:    "Empty String",
			input:   "",
			want:    w,
			wantErr: true,
		},
		{
			name:    "Rubbish String",
			input:   "rubbish.string",
			want:    w,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PurlFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.PurlFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Got: %v: '%v' '%v' '%v' '%v' '%v' '%v'", got, got.Type, got.Namespace, got.Name, got.Version, got.Qualifiers, got.Subpath)
			t.Logf("Exp: %v: '%v' '%v' '%v' '%v' '%v' '%v'", tt.want, tt.want.Type, tt.want.Namespace, tt.want.Name, tt.want.Version, tt.want.Qualifiers, tt.want.Subpath)
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.PurlFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPurlNameFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "Maven",
			input: "pkg:maven/io.prestosql/Presto-main@v1.0",
			want:  "io.prestosql/presto-main",
		},
		{
			name:  "NPM1",
			input: "pkg:npm/%40babel/Core",
			want:  "%40babel/Core",
		},
		{
			name:  "NPM2",
			input: "pkg:npm/%40babel/Core@7.0.0",
			want:  "%40babel/Core",
		},
		{
			name:  "NPM3",
			input: "pkg:npm/Core@0.0.1",
			want:  "Core",
		},
		{
			name:    "Empty String",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Rubbish String",
			input:   "rubbish.string",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PurlNameFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.PurlFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.PurlFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertPurlString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Maven",
			input: "pkg:maven/io.prestosql/presto-main@v1.0",
			want:  "pkg:maven/io.prestosql/presto-main@v1.0",
		},
		{
			name:  "Golang1",
			input: "pkg:golang/github.com/scanoss/papi",
			want:  "pkg:github/scanoss/papi",
		},
		{
			name:  "Golang2",
			input: "pkg:golang/github.com/scanoss",
			want:  "pkg:github/scanoss",
		},
		{
			name:  "Golang3",
			input: "pkg:golang/github.com/scanoss/papi/v2",
			want:  "pkg:github/scanoss/papi",
		},
		{
			name:  "Golang4",
			input: "pkg:golang/github.com/scanoss/papi/v3",
			want:  "pkg:github/scanoss/papi",
		},
		{
			name:  "Empty String",
			input: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertGoPurlStringToGithub(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.PurlFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPurlUrl(t *testing.T) {
	tests := []struct {
		name    string
		pname   string
		ptype   string
		want    string
		wantErr bool
	}{
		{
			name:  "GitHub",
			pname: "scanoss/scanoss.py",
			ptype: "github",
			want:  "https://github.com/scanoss/scanoss.py",
		},
		{
			name:  "Maven",
			pname: "io.prestosql/presto-main",
			ptype: "maven",
			want:  "https://mvnrepository.com/artifact/io.prestosql/presto-main",
		},
		{
			name:  "NPM",
			pname: "%40babel/core",
			ptype: "npm",
			want:  "https://www.npmjs.com/package/%40babel/core",
		},
		{
			name:  "PyPI",
			pname: "scanoss",
			ptype: "pypi",
			want:  "https://pypi.org/project/scanoss",
		},
		{
			name:  "Gem",
			pname: "tablestyle",
			ptype: "gem",
			want:  "https://rubygems.org/gems/tablestyle",
		},
		{
			name:  "Golang",
			pname: "github.com/scanoss/papi",
			ptype: "golang",
			want:  "https://pkg.go.dev/github.com/scanoss/papi",
		},
		{
			name:  "NuGet",
			pname: "System.Buffers",
			ptype: "nuget",
			want:  "https://www.nuget.org/packages/System.Buffers",
		},
		{
			name:    "Empty String1",
			pname:   "",
			ptype:   "gem",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty String2",
			pname:   "io.prestosql/presto-main",
			ptype:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Rubbish String",
			pname:   "rubbish.string",
			ptype:   "rubbish.string",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProjectUrl(tt.pname, tt.ptype)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.ProjectUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.ProjectUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionFromReq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Req1",
			input: "v1.0.0",
			want:  "v1.0.0",
		},
		{
			name:  "Req2",
			input: ">1.0.0",
			want:  "",
		},
		{
			name:  "Empty String",
			input: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetVersionFromReq(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.GetVersionFromReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
