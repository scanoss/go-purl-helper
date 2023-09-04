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
	"errors"
	"github.com/package-url/packageurl-go"
	"regexp"
)

import (
	"fmt"
	"strings"
)

var pkgRegex = regexp.MustCompile(`^pkg:(?P<type>\w+)/(?P<name>.+)$`) // regex to parse purl name from purl string
var typeRegex = regexp.MustCompile(`^(npm|nuget)$`)                   // regex to parse purl types that should not be lower cased
var vRegex = regexp.MustCompile(`^(=|==|)(?P<name>\w+\S+)$`)          // regex to parse purl name from purl string

// PurlFromString takes an input Purl string and returns a decomposed structure of all the elements
func PurlFromString(purlString string) (packageurl.PackageURL, error) {
	if len(purlString) == 0 {
		return packageurl.PackageURL{}, errors.New("no Purl string specified to parse")
	}
	purl, err := packageurl.FromString(purlString)
	if err != nil {
		return packageurl.PackageURL{}, err
	}
	return purl, nil
}

// PurlNameFromString take an input Purl string and returns the Purl Name only
func PurlNameFromString(purlString string) (string, error) {
	if len(purlString) == 0 {
		return "", fmt.Errorf("no purl string supplied to parse")
	}
	matches := pkgRegex.FindStringSubmatch(purlString)
	if len(matches) > 0 {
		ti := pkgRegex.SubexpIndex("type")
		ni := pkgRegex.SubexpIndex("name")
		if ni >= 0 {
			// Remove any version@/subpath?/qualifiers# info from the PURL
			pn := strings.Split(strings.Split(strings.Split(matches[ni], "@")[0], "?")[0], "#")[0]
			// Lowercase the purl name if it's not on the exclusion list (defined in the regex)
			if ti >= 0 && !typeRegex.MatchString(matches[ti]) {
				pn = strings.ToLower(pn)
			}
			return pn, nil
		}
	}
	return "", fmt.Errorf("no purl name found in '%v'", purlString)
}

// ConvertGoPurlStringToGithub takes an input PURL string and converts it to its GitHub equivalent if possible
func ConvertGoPurlStringToGithub(purlString string) string {
	// Replace Golang GitHub package reference with just GitHub
	if len(purlString) > 0 && strings.HasPrefix(purlString, "pkg:golang/github.com/") {
		s := strings.Replace(purlString, "pkg:golang/github.com/", "pkg:github/", -1)
		p := strings.Split(s, "/")
		if len(p) >= 3 {
			return fmt.Sprintf("%s/%s/%s", p[0], p[1], p[2]) // Only return the GitHub part of the url
		}
		return s
	}
	return purlString
}

// GetVersionFromReq parses a requirement string looking for an exact version specifier
func GetVersionFromReq(purlReq string) string {
	matches := vRegex.FindStringSubmatch(purlReq)
	if len(matches) > 0 {
		ni := vRegex.SubexpIndex("name")
		if ni >= 0 {
			return matches[ni]
		}
	}
	return ""
}

// ProjectUrl returns a browsable URL for the given purl type and name
func ProjectUrl(purlName, purlType string) (string, error) {
	if len(purlName) == 0 {
		return "", fmt.Errorf("no purl name supplied")
	}
	if len(purlType) == 0 {
		return "", fmt.Errorf("no purl type supplied")
	}
	switch purlType {
	case "github":
		return fmt.Sprintf("https://github.com/%v", purlName), nil
	case "npm":
		return fmt.Sprintf("https://www.npmjs.com/package/%v", purlName), nil
	case "maven":
		return fmt.Sprintf("https://mvnrepository.com/artifact/%v", purlName), nil
	case "gem":
		return fmt.Sprintf("https://rubygems.org/gems/%v", purlName), nil
	case "pypi":
		return fmt.Sprintf("https://pypi.org/project/%v", purlName), nil
	case "golang":
		return fmt.Sprintf("https://pkg.go.dev/%v", purlName), nil
	}
	return "", fmt.Errorf("no url prefix found for '%v': %v", purlType, purlName)
}
