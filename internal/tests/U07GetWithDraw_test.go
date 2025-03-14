package tests

// Basic imports
import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Repinoid/diploma56/internal/securitate"
)

func (suite *TstHandlers) Test07GetDraws() {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		testName string
		username string
		want     want
	}{
		// {
		// 	testName: "Right case1",
		// 	username: "user01",
		// 	want: want{
		// 		code:        http.StatusOK,
		// 		contentType: "application/json",
		// 	},
		// },
		{
			testName: "Right case2",
			username: "user02",
			want: want{
				code:        http.StatusNoContent,
				contentType: "application/json",
			},
		},
	}

	
	for _, tt := range tests {
		suite.Run(tt.testName, func() {
			var token string
			Interbase.GetToken(suite.ctx, tt.username, &token)
			request := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer <"+token+">")
			Interbase.GetWithDrawals(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			suite.Require().NoError(err)
			suite.Assert().Equal(tt.want.code, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				orda := []securitate.WithStruct{}
				err = json.Unmarshal([]byte(resBody), &orda)
				suite.Assert().NoError(err)
				log.Printf("%+v\n", orda)

				//	assert.JSONEq(t, tt.want.response, string(resBody))
				//suite.Assert().Equal(tt.want.contentType, res.Header.Get("Content-Type"))

			}
		})
	}
}
