package http

import "net/http"

func (a *Adapter) Register(w http.ResponseWriter, r *http.Request) {
	// credentials := user.User{}
	// if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
	// 	http.Error(w, "", http.StatusBadRequest)
	// 	return
	// }
	// if err := credentials.IsValid(); err != nil {
	// 	http.Error(w, "", http.StatusBadRequest)
	// 	return
	// }
	// signedToken, err := a.auth.Register(r.Context(), credentials)
	// if err != nil {
	// 	if errors.Is(err, user.ErrDuplicateLogin) {
	// 		w.WriteHeader(http.StatusConflict)
	// 	} else {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	}
	// 	return
	// }
	// w.Header().Set("Authorization", signedToken)
	w.WriteHeader(http.StatusOK)
}
