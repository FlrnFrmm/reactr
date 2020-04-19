package hive

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/suborbital/gust/gapi"
)

// Server is a hive server
type Server struct {
	*gapi.Server
	h        *Hive
	inFlight map[string]*Result
	sync.Mutex
}

func newServer(h *Hive, opts ...gapi.OptionsModifier) *Server {
	s := gapi.New(opts...)

	server := &Server{
		Server:   s,
		Mutex:    sync.Mutex{},
		h:        h,
		inFlight: make(map[string]*Result),
	}

	server.POST("/do/:jobtype", server.scheduleHandler())
	server.GET("/then/:id", server.thenHandler())

	return server
}

type doResponse struct {
	ResultID string `json:"resultId"`
}

func (s *Server) scheduleHandler() gapi.HandlerFunc {
	return func(r *http.Request, ctx *gapi.Ctx) (interface{}, error) {
		jobType := ctx.Params.ByName("jobtype")
		if jobType == "" {
			return nil, gapi.E(http.StatusBadRequest, "missing jobtype")
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, gapi.E(http.StatusInternalServerError, "failed to read request body")
		}
		defer r.Body.Close()

		res := s.h.Do(NewJob(jobType, data))

		then := r.URL.Query().Get("then")
		if then == "true" {
			result, err := res.Then()
			if err != nil {
				return nil, gapi.E(http.StatusInternalServerError, errors.Wrap(err, "job resulted in error").Error())
			}

			return result, nil
		}

		s.addInFlight(res)

		resp := doResponse{
			ResultID: res.ID,
		}

		return resp, nil
	}
}

func (s *Server) thenHandler() gapi.HandlerFunc {
	return func(r *http.Request, ctx *gapi.Ctx) (interface{}, error) {
		id := ctx.Params.ByName("id")
		if len(id) != 24 {
			return nil, gapi.E(http.StatusBadRequest, "invalid result ID")
		}

		res := s.getInFlight(id)
		if res == nil {
			return nil, gapi.E(http.StatusNotFound, fmt.Sprintf("result with ID %s not found", id))
		}

		defer s.removeInFlight(id)

		result, err := res.Then()
		if err != nil {
			return nil, gapi.E(http.StatusInternalServerError, errors.Wrap(err, "job resulted in error").Error())
		}

		return result, nil
	}
}

func (s *Server) addInFlight(r *Result) {
	s.Lock()
	defer s.Unlock()

	s.inFlight[r.ID] = r
}

func (s *Server) getInFlight(id string) *Result {
	s.Lock()
	defer s.Unlock()

	r, ok := s.inFlight[id]
	if !ok {
		return nil
	}

	return r
}

func (s *Server) removeInFlight(id string) {
	s.Lock()
	defer s.Unlock()

	delete(s.inFlight, id)
}