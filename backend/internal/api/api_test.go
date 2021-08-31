package api

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/kyma-incubator/Kyma-Showcase/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const fixedID = "FEA98D88-0669-4FFD-B17A-8F80BB97C381"

func TestDBGetHandler(t *testing.T) {
	const key = fixedID
	t.Run("should return error with status code 404 when url is wrong", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/{id}/wrong", nil)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("DBGETHANDLER: 404 not found")

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 0)
		assert.Contains(t, recorder.Body.String(), err.Error())
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("should return error with status code 500 when database does not respond", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/"+key, nil)
		vars := map[string]string{
			"id": key,
		}
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("GETFROMDB: database not respond error")
		dbManagerMock.On("GetFromDB", key).Return(nil, err)

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Contains(t, recorder.Body.String(), err.Error())
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("should return error with status code 500 when key does not exist in database", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/"+key, nil)
		vars := map[string]string{
			"id": key,
		}
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("GETFROMDB:key " + key + " does not exist")
		dbManagerMock.On("GetFromDB", key).Return(nil, err)

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Contains(t, recorder.Body.String(), err.Error())
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("should return error with status code 500 when key has no value assigned", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/"+key, nil)
		vars := map[string]string{
			"id": key,
		}
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("GETFROMDB:for key " + key + " value is empty")
		dbManagerMock.On("GetFromDB", key).Return("", err)

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), err.Error())
	})

	t.Run("should return error with status code 500 when key exist but value is not json", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/"+key, nil)
		value := "not json"
		vars := map[string]string{
			"id": key,
		}
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetFromDB", key).Return(value, nil)

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Contains(t, recorder.Body.String(), "DBGETHANDLER: failed to convert marshal to json:")
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("should return status code 200 when key exists in database and value is correct", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/"+key, nil)
		value := `
			{
				"url":"raccoon.com",
				"gcp":"image.png",
				"img":"image.png"
			}`
		vars := map[string]string{
			"id": key,
		}
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetFromDB", key).Return(value, nil)

		//when
		testSubject.DBGetHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Equal(t, value, recorder.Body.String())
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestDBGetAllHandler(t *testing.T) {
	t.Run("should return error with status code 404 when url is wrong", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images/wrong", nil)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("DBGETALLHANDLER: 404 not found")

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 0)
		assert.Contains(t, recorder.Body.String(), err.Error())
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("should return empty value when database is empty", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)

		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return(nil, nil)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		assert.Equal(t, http.StatusOK, recorder.Code)

	})

	t.Run("should return 500 code when all values are empty", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)

		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1", "2", "3"}, nil)
		dbManagerMock.On("GetFromDB", "1").Return("", errors.New("value is empty"))
		dbManagerMock.On("GetFromDB", "2").Return("", errors.New("value is empty"))
		dbManagerMock.On("GetFromDB", "3").Return("", errors.New("value is empty"))

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertNotCalled(t, "GetFromDB", "2")
		dbManagerMock.AssertNotCalled(t, "GetFromDB", "3")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1) // 1 poniewaz jest return po pierwszej pustej wartosci
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	})

	t.Run("should return 500 code when one of the values is empty", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)
		firstReturn := `
			{
				"id":"abcd1234",
				"content":"base64",
				"gcp":"JSON1",
				"status":false
			}`

		secondReturn := `
			{
				"id":"zaqwsx",
				"content":"base64_2",
				"gcp":"JSON2",
				"status":false
			}
			`

		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1", "2", "3"}, nil)
		dbManagerMock.On("GetFromDB", "1").Return(firstReturn, nil)
		dbManagerMock.On("GetFromDB", "2").Return(secondReturn, nil)
		dbManagerMock.On("GetFromDB", "3").Return("", errors.New("value is empty"))
		//handler := http.HandleFunc(DBGetAllHandler)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertCalled(t, "GetFromDB", "2")
		dbManagerMock.AssertCalled(t, "GetFromDB", "3")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 3)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	})

	t.Run("should return 500 code when error occurred during type assertion", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1"}, nil)
		dbManagerMock.On("GetFromDB", "1").Return(100, nil)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("should return JSON array compatible with given data ", func(t *testing.T) {

		//given
		expected := `[{` +
			`"id":"abcd1234",` +
			`"content":"base64",` +
			`"gcp":"JSON1",` +
			`"status":false` +
			`},` +
			`{` +
			`"id":"zaqwsx",` +
			`"content":"base64_2",` +
			`"gcp":"JSON2",` +
			`"status":false` +
			`}]`

		req, err := http.NewRequest("GET", "/v1/images", nil)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1", "2"}, nil)

		firstReturn := `
			{
				"id":"abcd1234",
				"content":"base64",
				"gcp":"JSON1",
				"status":false
			}`

		secondReturn := `
			{
				"id":"zaqwsx",
				"content":"base64_2",
				"gcp":"JSON2",
				"status":false
			}
			`

		dbManagerMock.On("GetFromDB", "1").Return(firstReturn, nil)
		dbManagerMock.On("GetFromDB", "2").Return(secondReturn, nil)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertCalled(t, "GetFromDB", "2")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 2)
		assert.Equal(t, expected, recorder.Body.String())
	})

	t.Run("should return error while is error during unmarshal", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1", "2"}, nil)

		firstReturn := `
			{
				id":"abcd1234",
				"content":"base64",
				"gcp":"JSON1",
				"status":false
			}`

		secondReturn := `
			{
				"id":"zaqwsx"
				"content":"base64_2",
				"gcp":"JSON2"
				"status":false
			}
			`

		dbManagerMock.On("GetFromDB", "1").Return(firstReturn, nil)
		dbManagerMock.On("GetFromDB", "2").Return(secondReturn, nil)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertNotCalled(t, "GetFromDB", "2")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1) //return po pierwszym blednym odczycie
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("should not return error while function can create valid JSON array", func(t *testing.T) {

		//given
		req, err := http.NewRequest("GET", "/v1/images", nil)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetAllKeys").Return([]string{"1", "2"}, nil)

		firstReturn := `
			{
				"url":"raccoon.com",
				"gcp":"image.png",
				"img":"image.png"
			}`

		secondReturn := `
			{
				"url":"raccoon.com",
				"gcp":"image2.png",
				"img":"image2.png"
			}
			`

		dbManagerMock.On("GetFromDB", "1").Return(firstReturn, nil)
		dbManagerMock.On("GetFromDB", "2").Return(secondReturn, nil)

		//when
		testSubject.DBGetAllHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetAllKeys", 1)
		dbManagerMock.AssertCalled(t, "GetFromDB", "1")
		dbManagerMock.AssertCalled(t, "GetFromDB", "2")
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 2)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestDBPostHandler(t *testing.T) {
	t.Run("should return error with status code 404 when url is wrong", func(t *testing.T) {

		//given
		req, err := http.NewRequest("POST", "/v1/images/wrong", nil)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("POST: 404 not found")

		//when
		testSubject.DBPostHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 0)
		assert.Contains(t, recorder.Body.String(), err.Error())
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
	t.Run("should return 400 when Content-Type is incorrect", func(t *testing.T) {

		//given
		var jsonStr = []byte(`{"content":"raccoon"}`)
		req, err := http.NewRequest("POST", "/v1/images", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/golang")
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)

		//when
		testSubject.DBPostHandler(recorder, req)

		//then
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should return 400 error when request body is not json", func(t *testing.T) {

		//given
		var jsonStr = "string"
		req, err := http.NewRequest("POST", "/v1/images", bytes.NewBuffer([]byte(jsonStr)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)

		//when
		testSubject.DBPostHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "InsertToDB", 0)
		assert.Contains(t, recorder.Body.String(), "POST: invalid input:")
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should return 500 error when unable to insert json to db", func(t *testing.T) {

		//given
		var jsonStr = `{` +
			`"id":"` + fixedID + `",` +
			`"content":"base64",` +
			`"gcp":"JSON1",` +
			`"status":false` +
			`}`
		req, err := http.NewRequest("POST", "/v1/images", bytes.NewBuffer([]byte(jsonStr)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		err = errors.New("failed to insert json to db")
		dbManagerMock.On("InsertToDB", fixedID, jsonStr).Return(err)

		//when
		testSubject.DBPostHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "InsertToDB", 1)
		assert.Equal(t, "POST: failed to insert values to database: "+err.Error()+"\n", recorder.Body.String())
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	})

	t.Run("should return 200 code when request, data and connection with database are correct", func(t *testing.T) {

		//given
		value := `{` +
			`"id":"` + fixedID + `",` +
			`"content":"base64",` +
			`"gcp":"JSON1",` +
			`"status":false` +
			`}`
		req, err := http.NewRequest("POST", "/v1/images", bytes.NewBuffer([]byte(value)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("InsertToDB", fixedID, value).Return(nil)

		//when
		testSubject.DBPostHandler(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "InsertToDB", 1)
		assert.Contains(t, recorder.Body.String(), fixedID)
		assert.Equal(t, http.StatusOK, recorder.Code)

	})
}
func TestUpdate(t *testing.T) {
	t.Run("UPDATE TESTY", func(t *testing.T) {

		//given
		value := `{` +
			`labels:[aaa,bbb]`+
			`}`
		returnedValue := `{` +
			`"id":"` + fixedID + `",` +
			`"content":"base64",` +
			//`"gcp": "[{labels:[aaa,bbb]}, {mood:[ccc,ddd]}}",` +
			`"status":false` +
			`}`

		req, err := http.NewRequest("PUT", "/v1/images/"+fixedID, bytes.NewBuffer([]byte(value)))
		vars := map[string]string{
			"id": fixedID,
		}
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, vars)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		dbManagerMock := mocks.DBManager{}
		idMock := mocks.IdGenerator{}
		idMock.On("NewID").Return(fixedID, nil)
		testSubject := NewHandler(&dbManagerMock, &idMock)
		dbManagerMock.On("GetFromDB", fixedID).Return(returnedValue, nil)

		//when
		testSubject.Update(recorder, req)

		//then
		dbManagerMock.AssertNumberOfCalls(t, "GetFromDB", 1)
		//assert.Equal(t, value, recorder.Body.String())
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

}
