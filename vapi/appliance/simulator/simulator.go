/*
Copyright (c) 2022 VMware, Inc. All Rights Reserved.

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

package simulator

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/appliance/access/consolecli"
	"github.com/vmware/govmomi/vapi/appliance/access/dcui"
	"github.com/vmware/govmomi/vapi/appliance/access/shell"
	"github.com/vmware/govmomi/vapi/appliance/access/ssh"
	vapi "github.com/vmware/govmomi/vapi/simulator"
)

func init() {
	simulator.RegisterEndpoint(func(s *simulator.Service, r *simulator.Registry) {
		New(s.Listen).Register(s, r)
	})
}

// Handler implements the Appliance API simulator
type Handler struct {
	URL        *url.URL
	consolecli consolecli.Access
	dcui       dcui.Access
	ssh        ssh.Access
	shell      shell.Access
}

// New creates a Handler instance
func New(u *url.URL) *Handler {
	return &Handler{
		URL:        nil,
		consolecli: consolecli.Access{Enabled: false},
		dcui:       dcui.Access{Enabled: false},
		ssh:        ssh.Access{Enabled: false},
		shell:      shell.Access{Enabled: false, Timeout: 0},
	}
}

// Register Appliance Management API paths with the vapi simulator's http.ServeMux
func (h *Handler) Register(s *simulator.Service, r *simulator.Registry) {
	s.HandleFunc(consolecli.Path, h.consoleCLIAccess)
	s.HandleFunc(dcui.Path, h.dcuiAccess)
	s.HandleFunc(ssh.Path, h.sshAccess)
	s.HandleFunc(shell.Path, h.shellAccess)
}

func (h *Handler) decode(r *http.Request, w http.ResponseWriter, val interface{}) bool {
	return Decode(r, w, val)
}

// Decode decodes the request Body into val, returns true on success, otherwise false.
func Decode(request *http.Request, writer http.ResponseWriter, val interface{}) bool {
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(val)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.RequestURI, err)
		return false
	}
	return true
}

func (h *Handler) consoleCLIAccess(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		vapi.StatusOK(writer, h.consolecli.Enabled)
	case http.MethodPut:
		var input consolecli.Access
		if h.decode(request, writer, &input) {
			h.consolecli.Enabled = input.Enabled
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.NotFound(writer, request)
	}
}

func (h *Handler) dcuiAccess(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		vapi.StatusOK(writer, h.dcui.Enabled)
	case http.MethodPut:
		var input dcui.Access
		if h.decode(request, writer, &input) {
			h.dcui.Enabled = input.Enabled
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.NotFound(writer, request)
	}
}

func (h *Handler) sshAccess(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		vapi.StatusOK(writer, h.ssh.Enabled)
	case http.MethodPut:
		var input ssh.Access
		if h.decode(request, writer, &input) {
			h.ssh.Enabled = input.Enabled
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.NotFound(writer, request)
	}
}

func (h *Handler) shellAccess(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		vapi.StatusOK(writer, h.shell)
	case http.MethodPut:
		var input shell.Access
		if h.decode(request, writer, &input) {
			h.shell.Enabled = input.Enabled
			h.shell.Timeout = input.Timeout
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.NotFound(writer, request)
	}
}
