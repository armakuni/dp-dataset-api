package instance_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errs "github.com/ONSdigital/dp-dataset-api/apierrors"
	"github.com/ONSdigital/dp-dataset-api/instance"
	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/dp-dataset-api/store/datastoretest"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const secretKey = "coffee"
const host = "http://localhost:8080"

var internalError = errors.New("internal error")

func createRequestWithToken(method, url string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, url, body)
	r.Header.Add("internal-token", secretKey)
	return r
}

func TestGetInstancesReturnsOK(t *testing.T) {
	t.Parallel()
	Convey("Get instances returns a ok status code", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstancesFunc: func([]string) (*models.InstanceResults, error) {
				return &models.InstanceResults{}, nil
			},
		}

		instance := &instance.Store{Host: "http://lochost://8080", Storer: mockedDataStore}
		instance.GetList(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstancesCalls()), ShouldEqual, 1)
	})
}

func TestGetInstancesFiltersOnState(t *testing.T) {
	t.Parallel()
	Convey("Get instances filtered by a single state value returns only instances with that value", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances?state=completed", nil)
		w := httptest.NewRecorder()
		var result []string

		mockedDataStore := &storetest.StorerMock{
			GetInstancesFunc: func(filterString []string) (*models.InstanceResults, error) {
				result = filterString
				return &models.InstanceResults{}, nil
			},
		}

		instance := &instance.Store{Host: "http://lochost://8080", Storer: mockedDataStore}
		instance.GetList(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstancesCalls()), ShouldEqual, 1)
		So(result, ShouldResemble, []string{models.CompletedState})
	})

	Convey("Get instances filtered by multiple state values returns only instances with those values", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances?state=completed,edition-confirmed", nil)
		w := httptest.NewRecorder()
		var result []string

		mockedDataStore := &storetest.StorerMock{
			GetInstancesFunc: func(filterString []string) (*models.InstanceResults, error) {
				result = filterString
				return &models.InstanceResults{}, nil
			},
		}

		instance := &instance.Store{Host: "http://lochost://8080", Storer: mockedDataStore}
		instance.GetList(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstancesCalls()), ShouldEqual, 1)
		So(result, ShouldResemble, []string{models.CompletedState, models.EditionConfirmedState})
	})
}

func TestGetInstancesReturnsError(t *testing.T) {
	t.Parallel()
	Convey("Get instances returns an internal error", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstancesFunc: func([]string) (*models.InstanceResults, error) {
				return nil, internalError
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.GetList(w, r)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(len(mockedDataStore.GetInstancesCalls()), ShouldEqual, 1)
	})

	Convey("Get instances returns bad request error", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances?state=foo", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.GetList(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(len(mockedDataStore.GetInstancesCalls()), ShouldEqual, 0)
	})
}

func TestGetInstanceReturnsOK(t *testing.T) {
	t.Parallel()
	Convey("Get instance returns a ok status code", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances/123", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(ID string) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Get(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
	})
}

func TestGetInstanceReturnsInternalError(t *testing.T) {
	t.Parallel()
	Convey("Get instance returns an internal error", t, func() {
		r := createRequestWithToken("GET", "http://localhost:21800/instances/123", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(ID string) (*models.Instance, error) {
				return nil, internalError
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Get(w, r)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
	})
}

func TestAddInstancesReturnsCreated(t *testing.T) {
	t.Parallel()
	Convey("Add instance returns a created code", t, func() {
		body := strings.NewReader(`{"links": { "job": { "id":"123-456", "href":"http://localhost:2200/jobs/123-456" } } }`)
		r := createRequestWithToken("POST", "http://localhost:21800/instances", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			AddInstanceFunc: func(*models.Instance) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Add(w, r)

		So(w.Code, ShouldEqual, http.StatusCreated)
		So(len(mockedDataStore.AddInstanceCalls()), ShouldEqual, 1)
	})
}

func TestAddInstancesReturnsBadRequest(t *testing.T) {
	t.Parallel()
	Convey("Add instance returns a bad request with invalid json", t, func() {
		body := strings.NewReader(`{`)
		r := createRequestWithToken("POST", "http://localhost:21800/instances", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			AddInstanceFunc: func(*models.Instance) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Add(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("Add instance returns a bad request with a empty json", t, func() {
		body := strings.NewReader(`{}`)
		r := createRequestWithToken("POST", "http://localhost:21800/instances", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			AddInstanceFunc: func(*models.Instance) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Add(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestAddInstancesReturnsInternalError(t *testing.T) {
	t.Parallel()
	Convey("Add instance returns an internal error", t, func() {
		body := strings.NewReader(`{"links": {"job": { "id":"123-456", "href":"http://localhost:2200/jobs/123-456" } } }`)
		r := createRequestWithToken("POST", "http://localhost:21800/instances", body)
		w := httptest.NewRecorder()
		mockedDataStore := &storetest.StorerMock{
			AddInstanceFunc: func(*models.Instance) (*models.Instance, error) {
				return nil, internalError
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Add(w, r)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(len(mockedDataStore.AddInstanceCalls()), ShouldEqual, 1)
	})
}

func TestUpdateInstanceReturnsOk(t *testing.T) {
	t.Parallel()
	Convey("when an instance has a state of created", t, func() {
		body := strings.NewReader(`{"state":"created"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 1)
	})

	Convey("when an instance changes its state to edition-confirmed", t, func() {
		body := strings.NewReader(`{"state":"edition-confirmed", "edition": "2017"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()

		currentInstanceTestData := &models.Instance{
			Edition: "2017",
			Links: &models.InstanceLinks{
				Job: &models.IDLink{
					ID:   "7654",
					HRef: "job-link",
				},
				Dataset: &models.IDLink{
					ID:   "4567",
					HRef: "dataset-link",
				},
				Self: &models.IDLink{
					HRef: "self-link",
				},
			},
			State: models.CompletedState,
		}

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return currentInstanceTestData, nil
			},
			GetEditionFunc: func(datasetID string, edition string, state string) (*models.EditionUpdate, error) {
				return nil, errs.ErrEditionNotFound
			},
			UpsertEditionFunc: func(datasetID, edition string, editionDoc *models.EditionUpdate) error {
				return nil
			},
			GetNextVersionFunc: func(string, string) (int, error) {
				return 1, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.GetEditionCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpsertEditionCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.GetNextVersionCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 1)
	})
}

func TestUpdateInstanceFailure(t *testing.T) {
	t.Parallel()
	Convey("when the json body is in the incorrect structure return a bad request error", t, func() {
		body := strings.NewReader(`{"state":`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()
		mockedDataStore := &storetest.StorerMock{}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 0)
	})

	Convey("when the instance does not exist return status not found", t, func() {
		body := strings.NewReader(`{"edition": "2017"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()
		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return nil, errs.ErrInstanceNotFound
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestUpdatePublishedInstanceToCompletedReturnsForbidden(t *testing.T) {
	t.Parallel()
	Convey("Given a 'published' instance, when we update to 'completed' then we get a bad-request error", t, func() {
		body := strings.NewReader(`{"state":"completed"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/1235", body)
		w := httptest.NewRecorder()

		currentInstanceTestData := &models.Instance{
			Edition: "2017",
			Links: &models.InstanceLinks{
				Dataset: &models.IDLink{
					ID: "4567",
				},
			},
			State: models.PublishedState,
		}

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return currentInstanceTestData, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusForbidden)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestUpdateEditionConfirmedInstanceToCompletedReturnsForbidden(t *testing.T) {
	t.Parallel()
	Convey("update to an instance returns an internal error", t, func() {
		body := strings.NewReader(`{"state":"completed"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()

		currentInstanceTestData := &models.Instance{
			Edition: "2017",
			Links: &models.InstanceLinks{
				Dataset: &models.IDLink{
					ID: "4567",
				},
			},
			State: models.EditionConfirmedState,
		}

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return currentInstanceTestData, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusForbidden)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestInsertedObservationsReturnsOk(t *testing.T) {
	t.Parallel()
	Convey("Updateding the inserted observations returns ok", t, func() {
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/inserted_observations/200", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateObservationInsertedFunc: func(id string, ob int64) error {
				return nil
			},
		}
		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		router := mux.NewRouter()
		router.HandleFunc("/instances/{id}/inserted_observations/{inserted_observations}", instance.UpdateObservations).Methods("PUT")
		router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.UpdateObservationInsertedCalls()), ShouldEqual, 1)
	})
}

func TestInsertedObservationsReturnsBadRequest(t *testing.T) {
	t.Parallel()
	Convey("Updateding the inserted observations returns bad request", t, func() {
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/inserted_observations/aa12a", nil)
		w := httptest.NewRecorder()
		mockedDataStore := &storetest.StorerMock{}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.UpdateObservations(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestInsertedObservationsReturnsNotFound(t *testing.T) {
	t.Parallel()
	Convey("Updating the inserted observations returns not found", t, func() {
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/inserted_observations/200", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateObservationInsertedFunc: func(id string, ob int64) error {
				return errs.ErrInstanceNotFound
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		router := mux.NewRouter()
		router.HandleFunc("/instances/{id}/inserted_observations/{inserted_observations}", instance.UpdateObservations).Methods("PUT")
		router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(len(mockedDataStore.UpdateObservationInsertedCalls()), ShouldEqual, 1)
	})
}

func TestStore_UpdateImportTask_UpdateImportObservations(t *testing.T) {

	t.Parallel()
	Convey("update to an import task returns http 200 response if no errors occur", t, func() {
		body := strings.NewReader(`{"import_observations":{"state":"completed"}}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/import_tasks", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateImportObservationsTaskStateFunc: func(id string, state string) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		instance.UpdateImportTask(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.UpdateImportObservationsTaskStateCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateBuildHierarchyTaskStateCalls()), ShouldEqual, 0)
	})
}

func TestStore_UpdateImportTask_UpdateImportObservations_InvalidState(t *testing.T) {

	t.Parallel()
	Convey("update to an import task with an invalid state returns http 400 response", t, func() {
		body := strings.NewReader(`{"import_observations":{"state":"notvalid"}}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/import_tasks", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateImportObservationsTaskStateFunc: func(id string, state string) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		instance.UpdateImportTask(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(len(mockedDataStore.UpdateImportObservationsTaskStateCalls()), ShouldEqual, 0)
		So(len(mockedDataStore.UpdateBuildHierarchyTaskStateCalls()), ShouldEqual, 0)
	})
}

func TestStore_UpdateImportTask_UpdateBuildHierarchyTask_InvalidState(t *testing.T) {

	t.Parallel()
	Convey("update to an import task with an invalid state returns http 400 response", t, func() {
		body := strings.NewReader(`{"build_hierarchies":[{"state":"notvalid"}]}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/import_tasks", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateBuildHierarchyTaskStateFunc: func(id string, dimension string, state string) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		instance.UpdateImportTask(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(len(mockedDataStore.UpdateImportObservationsTaskStateCalls()), ShouldEqual, 0)
		So(len(mockedDataStore.UpdateBuildHierarchyTaskStateCalls()), ShouldEqual, 0)
	})
}

func TestStore_UpdateImportTask_UpdateBuildHierarchyTask(t *testing.T) {

	t.Parallel()
	Convey("update to an import task returns http 200 response if no errors occur", t, func() {
		body := strings.NewReader(`{"build_hierarchies":[{"state":"completed"}]}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/import_tasks", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateBuildHierarchyTaskStateFunc: func(id string, dimension string, state string) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		instance.UpdateImportTask(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.UpdateImportObservationsTaskStateCalls()), ShouldEqual, 0)
		So(len(mockedDataStore.UpdateBuildHierarchyTaskStateCalls()), ShouldEqual, 1)
	})
}

func TestStore_UpdateImportTask_ReturnsInternalError(t *testing.T) {

	t.Parallel()
	Convey("update to an import task returns an internal error", t, func() {
		body := strings.NewReader(`{"import_observations":{"state":"completed"}}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123/import_tasks", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			UpdateImportObservationsTaskStateFunc: func(id string, state string) error {
				return internalError
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		instance.UpdateImportTask(w, r)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(len(mockedDataStore.UpdateImportObservationsTaskStateCalls()), ShouldEqual, 1)
	})
}


func TestUpdateInstanceReturnsErrorWhenStateIsPublished(t *testing.T) {
	t.Parallel()
	Convey("when an instance has a state of published, then put request to change to it to completed ", t, func() {
		body := strings.NewReader(`{"state":"completed"}`)
		r := createRequestWithToken("PUT", "http://localhost:21800/instances/123", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return &models.Instance{State: models.PublishedState}, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.Update(w, r)

		So(w.Code, ShouldEqual, http.StatusForbidden)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
	})
}


func TestUpdateDimensionReturnsNotFound(t *testing.T) {
	t.Parallel()
	Convey("When update dimension return status not found", t, func() {
		r := createRequestWithToken("PUT", "http://localhost:22000/instances/dimensions/age", nil)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return nil, errs.ErrInstanceNotFound
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.UpdateDimension(w, r)

		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestUpdateDimensionReturnsForbidden(t *testing.T) {
	t.Parallel()
	Convey("When update dimension returns forbidden (for already published) ", t, func() {
		r := createRequestWithToken("PUT", "http://localhost:22000/instances/123/dimensions/age", nil)
		w := httptest.NewRecorder()

		currentInstanceTestData := &models.Instance{
			State: models.PublishedState,
		}

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return currentInstanceTestData, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.UpdateDimension(w, r)

		So(w.Code, ShouldEqual, http.StatusForbidden)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestUpdateDimensionReturnsBadRequest(t *testing.T) {
	t.Parallel()
	Convey("When update dimension returns bad request", t, func() {
		body := strings.NewReader("{")
		r := createRequestWithToken("PUT", "http://localhost:22000/instances/123/dimensions/age", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		instance.UpdateDimension(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestUpdateDimensionReturnsNotFoundWithWrongName(t *testing.T) {
	t.Parallel()
	Convey("When update dimension fails to update an instance", t, func() {
		body := strings.NewReader(`{"label":"notages"}`)
		r := createRequestWithToken("PUT", "http://localhost:22000/instances/123/dimensions/notage", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return &models.Instance{State: models.EditionConfirmedState,
					InstanceID: "123",
					Dimensions: []models.CodeList{{Name: "age", ID: "age"}}}, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}

		router := mux.NewRouter()
		router.HandleFunc("/instances/{id}/dimensions/{dimension}", instance.UpdateDimension).Methods("PUT")

		router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 0)
	})
}

func TestUpdateDimensionReturnsOk(t *testing.T) {
	t.Parallel()
	Convey("When update dimension fails to update an instance", t, func() {
		body := strings.NewReader(`{"label":"ages"}`)
		r := createRequestWithToken("PUT", "http://localhost:22000/instances/123/dimensions/age", body)
		w := httptest.NewRecorder()

		mockedDataStore := &storetest.StorerMock{
			GetInstanceFunc: func(id string) (*models.Instance, error) {
				return &models.Instance{State: models.EditionConfirmedState,
					InstanceID: "123",
					Dimensions: []models.CodeList{{Name: "age", ID: "age"}}}, nil
			},
			UpdateInstanceFunc: func(id string, i *models.Instance) error {
				return nil
			},
		}

		instance := &instance.Store{Host: host, Storer: mockedDataStore}
		router := mux.NewRouter()
		router.HandleFunc("/instances/{id}/dimensions/{dimension}", instance.UpdateDimension).Methods("PUT")

		router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(len(mockedDataStore.GetInstanceCalls()), ShouldEqual, 1)
		So(len(mockedDataStore.UpdateInstanceCalls()), ShouldEqual, 1)
	})
}
