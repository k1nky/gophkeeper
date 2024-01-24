package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/k1nky/gophkeeper/internal/protocol/rest"
)

func (a *Adapter) GetSecretMeta(w http.ResponseWriter, r *http.Request) {
	uk := chi.URLParam(r, "id")
	secret, err := a.keeper.GetSecretMeta(r.Context(), vault.UniqueKey(uk))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if secret == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := a.writeJSON(w, secret); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) GetSecretData(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	uk := chi.URLParam(r, "id")
	data, err := a.keeper.GetSecretData(r.Context(), vault.UniqueKey(uk), claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer data.Close()
	_, err = io.Copy(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) PutSecret(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	request := rest.NewSecretRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	meta := vault.Meta{
		UserID: claims.ID,
		Extra:  request.Extra,
	}
	data := vault.NewDataReader(vault.NewBytesBuffer([]byte(request.Secret)))
	newMeta, err := a.keeper.PutSecret(r.Context(), meta, data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := a.writeJSON(w, newMeta); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) PutSecretFile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	extra := r.FormValue("extra")
	meta := vault.Meta{
		UserID: claims.ID,
		Extra:  extra,
	}
	f, _, err := r.FormFile("secret")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data := vault.NewDataReader(f)
	defer data.Close()
	newMeta, err := a.keeper.PutSecret(r.Context(), meta, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := a.writeJSON(w, newMeta); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
