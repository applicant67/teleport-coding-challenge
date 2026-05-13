package server

import (
	"encoding/json"
	"fmt"
	"job-worker/lib"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var jobWorker *lib.JobWorker

func init() {
	jobWorker = lib.NewJobWorker()
}

// CreateJob handles the API request to create a new job.
// It writes the job details in response body.
func CreateJob(w http.ResponseWriter, r *http.Request) {
	request := lib.CreateJobRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Errorln(err)
		}
	}()

	if job, err := jobWorker.CreateJob(request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, job)
	}
}

// StopJob handles the API request to stop a job.
// It writes the confirmation in response body.
func StopJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := jobWorker.StopJob(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondMessage(w, http.StatusOK, fmt.Sprintf("Termination requested for job %v, please check job status", id))
	}
}

// GetJobStatus handles the API request to fetch the status of a job.
// It writes the job details in response body.
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if job, err := jobWorker.GetJobStatus(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, job)
	}
}
