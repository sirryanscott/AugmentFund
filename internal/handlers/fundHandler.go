package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/services"
	"github.com/go-chi/chi"
)

type FundHandler struct {
	FundService *services.FundService
}

func NewFundHandler(fundService *services.FundService) *FundHandler {
	return &FundHandler{FundService: fundService}
}

func (h *FundHandler) GetCapTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	funds, err := h.FundService.GetCapTables(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(funds)
}

func (h *FundHandler) GetCapTableByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	fund, err := h.FundService.GetCapTableByID(ctx, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(fund)
}

func (h *FundHandler) CreateFund(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fundBytes, err := ioReaderToBytes(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var fund data.Fund
	err = json.Unmarshal(fundBytes, &fund)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	funds, err := h.FundService.CreateFund(ctx, fund)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(funds)
}

func (h *FundHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transferDataBytes, err := ioReaderToBytes(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var transferData data.Transfer
	err = json.Unmarshal(transferDataBytes, &transferData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	funds, err := h.FundService.CreateTransfer(ctx, transferData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(funds)
}

func (h *FundHandler) GetTransferHistoryForFund(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	transferHistoryRecords, err := h.FundService.GetTransferHistoryForFund(ctx, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(transferHistoryRecords)
}
