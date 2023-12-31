package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchHandler struct {
	SearchUsecase domain.SearchUsecase
}

func NewSearchHandler(router *mux.Router, su domain.SearchUsecase) {
	handler := &SearchHandler{
		SearchUsecase: su,
	}

	router.HandleFunc("/v1/search/{searchStr}", handler.GetSearchData).Methods(http.MethodGet, http.MethodOptions)
}

// GetSearchData godoc
// @Summary      Search data
// @Description  Get search data by incoming string
// @Tags         Search
// @Produce 	 json
// @Param 		 searchStr path string true "The string to be searched for"
// @Success      200  {object} object{body=object{films=[]domain.Video, cast=[]domain.Cast}}
// @Failure      404  {object} object{err=string}
// @Failure      500  {object} object{err=string}
// @Router       /api/v1/search/{searchStr} [get]
func (h *SearchHandler) GetSearchData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchStr := vars["searchStr"]

	data, err := h.SearchUsecase.GetSearchData(searchStr)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetSearchData", err, "Failed to get search data")
		return
	}

	logs.Logger.Debug("data:", data)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"films": data.Films,
			"cast":  data.Cast,
		},
		http.StatusOK,
	)
}
