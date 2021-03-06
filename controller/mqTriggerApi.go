/*
Copyright 2017 The Fission Authors.

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

package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/client-go/1.5/pkg/api"

	"github.com/fission/fission"
	"github.com/fission/fission/tpr"
)

func (a *API) MessageQueueTriggerApiList(w http.ResponseWriter, r *http.Request) {
	//mqType := r.FormValue("mqtype") // ignored for now
	triggers, err := a.fissionClient.Messagequeuetriggers(api.NamespaceAll).List(api.ListOptions{})
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	resp, err := json.Marshal(triggers.Items)
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	a.respondWithSuccess(w, resp)
}

func (a *API) MessageQueueTriggerApiCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	var mqTrigger tpr.Messagequeuetrigger
	err = json.Unmarshal(body, &mqTrigger)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	err = validateResourceName(mqTrigger.Metadata.Name)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	tnew, err := a.fissionClient.Messagequeuetriggers(mqTrigger.Metadata.Namespace).Create(&mqTrigger)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	resp, err := json.Marshal(tnew.Metadata)
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	a.respondWithSuccess(w, resp)
}

func (a *API) MessageQueueTriggerApiGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["mqTrigger"]
	ns := vars["namespace"]
	if len(ns) == 0 {
		ns = api.NamespaceDefault
	}

	mqTrigger, err := a.fissionClient.Messagequeuetriggers(ns).Get(name)
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	resp, err := json.Marshal(mqTrigger)
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	a.respondWithSuccess(w, resp)
}

func (a *API) MessageQueueTriggerApiUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["mqTrigger"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	var mqTrigger tpr.Messagequeuetrigger
	err = json.Unmarshal(body, &mqTrigger)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	if name != mqTrigger.Metadata.Name {
		err = fission.MakeError(fission.ErrorInvalidArgument, "Message queue trigger name doesn't match URL")
		a.respondWithError(w, err)
		return
	}

	tnew, err := a.fissionClient.Messagequeuetriggers(mqTrigger.Metadata.Namespace).Update(&mqTrigger)
	if err != nil {
		a.respondWithError(w, err)
		return
	}

	resp, err := json.Marshal(tnew.Metadata)
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	a.respondWithSuccess(w, resp)
}

func (a *API) MessageQueueTriggerApiDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["mqTrigger"]
	ns := vars["namespace"]
	if len(ns) == 0 {
		ns = api.NamespaceDefault
	}

	err := a.fissionClient.Messagequeuetriggers(ns).Delete(name, &api.DeleteOptions{})
	if err != nil {
		a.respondWithError(w, err)
		return
	}
	a.respondWithSuccess(w, []byte(""))
}
