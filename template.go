/*
Copyright 2023 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// InitTemplates parses the templates.
func InitTemplates(pattern string) (*template.Template, error) {
	tmpls, err := template.New("html").ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("InitTemplates: %w", err)
	}
	return tmpls, nil
}

const MsgTemplateError = "Sorry, the server was unable to display this page. Please contact the administrator."

// RenderTemplate executes the named template with the given data and writes the result to the provided HTTP response writer.
// If an error occurs during template execution, the HTTP response status is set to Internal Server Error (HTTP 500), and the function returns the error.
// The caller must ensure no further writes are done for a non-nil error.
func RenderTemplate(t *template.Template, w http.ResponseWriter, name string, data interface{}) error {
	// handle nil template
	if t == nil {
		return errors.New("RenderTemplate: nil template")
	}

	// Create a buffer to store the template output since if an error occurs executing the template or writing its output, execution stops, but partial results may already have been written to the output writer.
	var tmplBuffer bytes.Buffer

	// Execute the template with the provided data.
	err := t.ExecuteTemplate(&tmplBuffer, name, data)
	if err != nil {
		// If an error occurs, set the HTTP response status to Internal Server Error (HTTP 500).
		http.Error(w, MsgTemplateError, http.StatusInternalServerError)
		return err
	}

	// Write the template output to the response writer and check for errors.
	_, writeErr := tmplBuffer.WriteTo(w)
	if writeErr != nil {
		// If an error occurs while writing to the response writer, return it.
		http.Error(w, MsgTemplateError, http.StatusInternalServerError)
		return writeErr
	}

	return nil
}
