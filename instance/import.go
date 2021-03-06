package instance

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	errs "github.com/ONSdigital/dp-dataset-api/apierrors"
	"github.com/ONSdigital/dp-dataset-api/models"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// UpdateObservations increments the count of inserted_observations against
// an instance
func (s *Store) UpdateObservations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	instanceID := vars["instance_id"]
	insert := vars["inserted_observations"]
	logData := log.Data{"instance_id": instanceID, "inserted_observations": insert}

	if err := func() error {
		observations, err := strconv.ParseInt(insert, 10, 64)
		if err != nil {
			log.Event(ctx, "update imported observations: failed to parse inserted_observations string to int", log.ERROR, log.Error(err), logData)
			return errs.ErrInsertedObservationsInvalidSyntax
		}

		if err = s.UpdateObservationInserted(instanceID, observations); err != nil {
			log.Event(ctx, "update imported observations: store.UpdateObservationInserted returned an error", log.ERROR, log.Error(err), logData)
			return err
		}

		return nil
	}(); err != nil {
		handleInstanceErr(ctx, err, w, logData)
		return
	}

	log.Event(ctx, "update imported observations: request successful", log.INFO, logData)
}

// UpdateImportTask updates any task in the request body against an instance
func (s *Store) UpdateImportTask(w http.ResponseWriter, r *http.Request) {

	defer dphttp.DrainBody(r)

	ctx := r.Context()
	vars := mux.Vars(r)
	instanceID := vars["instance_id"]
	logData := log.Data{"instance_id": instanceID}
	defer r.Body.Close()

	updateErr := func() *taskError {
		tasks, err := unmarshalImportTasks(r.Body)
		if err != nil {
			log.Event(ctx, "failed to unmarshal request body to UpdateImportTasks model", log.ERROR, log.Error(err), logData)
			return &taskError{err, http.StatusBadRequest}
		}

		validationErrs := make([]error, 0)
		var hasImportTasks bool

		if tasks.ImportObservations != nil {
			hasImportTasks = true
			if tasks.ImportObservations.State != "" {
				if tasks.ImportObservations.State != models.CompletedState {
					validationErrs = append(validationErrs, fmt.Errorf("bad request - invalid task state value for import observations: %v", tasks.ImportObservations.State))
				} else {
					if err := s.UpdateImportObservationsTaskState(instanceID, tasks.ImportObservations.State); err != nil {
						log.Event(ctx, "failed to update import observations task state", log.ERROR, log.Error(err), logData)
						return &taskError{err, http.StatusInternalServerError}
					}
				}
			} else {
				validationErrs = append(validationErrs, errors.New("bad request - invalid import observation task, must include state"))
			}
		}

		if tasks.BuildHierarchyTasks != nil {
			hasImportTasks = true
			var hasHierarchyImportTask bool
			for _, task := range tasks.BuildHierarchyTasks {
				hasHierarchyImportTask = true
				if err := models.ValidateImportTask(task.GenericTaskDetails); err != nil {
					validationErrs = append(validationErrs, err)
				} else {
					if err := s.UpdateBuildHierarchyTaskState(instanceID, task.DimensionName, task.State); err != nil {
						if err.Error() == errs.ErrNotFound.Error() {
							notFoundErr := task.DimensionName + " hierarchy import task does not exist"
							log.Event(ctx, notFoundErr, log.ERROR, log.Error(err), logData)
							return &taskError{errors.New(notFoundErr), http.StatusNotFound}
						}
						log.Event(ctx, "failed to update build hierarchy task state", log.ERROR, log.Error(err), logData)
						return &taskError{err, http.StatusInternalServerError}
					}
				}
			}
			if !hasHierarchyImportTask {
				validationErrs = append(validationErrs, errors.New("bad request - missing hierarchy task"))
			}
		}

		if tasks.BuildSearchIndexTasks != nil {
			hasImportTasks = true
			var hasSearchIndexImportTask bool
			for _, task := range tasks.BuildSearchIndexTasks {
				hasSearchIndexImportTask = true
				if err := models.ValidateImportTask(task.GenericTaskDetails); err != nil {
					validationErrs = append(validationErrs, err)
				} else {
					if err := s.UpdateBuildSearchTaskState(instanceID, task.DimensionName, task.State); err != nil {
						if err.Error() == "not found" {
							notFoundErr := task.DimensionName + " search index import task does not exist"
							log.Event(ctx, notFoundErr, log.ERROR, log.Error(err), logData)
							return &taskError{errors.New(notFoundErr), http.StatusNotFound}
						}
						log.Event(ctx, "failed to update build hierarchy task state", log.ERROR, log.Error(err), logData)
						return &taskError{err, http.StatusInternalServerError}
					}
				}
			}
			if !hasSearchIndexImportTask {
				validationErrs = append(validationErrs, errors.New("bad request - missing search index task"))
			}
		}

		if !hasImportTasks {
			validationErrs = append(validationErrs, errors.New("bad request - request body does not contain any import tasks"))
		}

		if len(validationErrs) > 0 {
			for _, err := range validationErrs {
				log.Event(ctx, "validation error", log.ERROR, log.Error(err), logData)
			}
			// todo: add all validation errors to the response
			return &taskError{validationErrs[0], http.StatusBadRequest}
		}
		return nil
	}()

	if updateErr != nil {
		log.Event(ctx, "updateImportTask endpoint: request unsuccessful", log.ERROR, log.Error(updateErr), logData)
		http.Error(w, updateErr.Error(), updateErr.status)
		return
	}

	log.Event(ctx, "updateImportTask endpoint: request successful", log.INFO, logData)
}

func unmarshalImportTasks(reader io.Reader) (*models.InstanceImportTasks, error) {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errs.ErrUnableToReadMessage
	}

	var tasks models.InstanceImportTasks
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}

	return &tasks, nil
}
