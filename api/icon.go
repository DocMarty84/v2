// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api // import "miniflux.app/api"

import (
	"net/http"
	"time"

	"miniflux.app/http/request"
	"miniflux.app/http/response"
	"miniflux.app/http/response/json"
)

func (h *handler) feedIcon(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")

	if !h.store.HasIcon(feedID) {
		json.NotFound(w, r)
		return
	}

	icon, err := h.store.IconByFeedID(request.UserID(r), feedID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if icon == nil {
		json.NotFound(w, r)
		return
	}

	json.OK(w, r, &feedIconResponse{
		ID:       icon.ID,
		MimeType: icon.MimeType,
		Data:     icon.DataURL(),
	})
}

func (h *handler) showIcon(w http.ResponseWriter, r *http.Request) {
	iconID := request.RouteInt64Param(r, "iconID")
	icon, err := h.store.IconByID(iconID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if icon == nil {
		json.NotFound(w, r)
		return
	}

	response.New(w, r).WithCaching(icon.Hash, 72*time.Hour, func(b *response.Builder) {
		b.WithHeader("Content-Type", icon.MimeType)
		b.WithBody(icon.Content)
		b.WithoutCompression()
		b.Write()
	})
}
