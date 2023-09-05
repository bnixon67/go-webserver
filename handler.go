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
	"html/template"
)

// Handler encapsulates the behavior for processing HTTP requests.
type Handler struct {
	AppName string             // AppName is the name of the application using this handler.
	Tmpl    *template.Template // Tmpl holds the parsed templates to be rendered.
}

// NewHandler returns a new Handler instance with the given application name and template.
func NewHandler(appName string, tmpl *template.Template) *Handler {
	return &Handler{AppName: appName, Tmpl: tmpl}
}
