package securitate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Repinoid/diploma56/internal/models"
)

// api/user/register
func (dataBase *DBstruct) RegisterUser(rwr http.ResponseWriter, req *http.Request) {

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debug("not application/json\n")
		return
	}

	rwr.Header().Set("Content-Type", "application/json")

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()

	logos := struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}{}
	err = json.Unmarshal([]byte(telo), &logos)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)
		return
	}

	err = dataBase.IfUserExists(req.Context(), logos.UserName)
	if err == nil {
		fmt.Printf("User exists %v\n", err)
		rwr.WriteHeader(http.StatusConflict) // 409 — логин уже занят;
		fmt.Fprintf(rwr, `{"status":"StatusConflict"}`)
		return
	}

	Token, err := BuildJWTString(logos.UserName, []byte(SecretKey))
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("BuildJWTString %+v\n", err)
		return
	}

	err = dataBase.AddUser(req.Context(), logos.UserName, logos.Password, Token)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debugf("addUser %+v %+v\n", logos, err)
		return
	}
	rwr.Header().Add("Authorization", "Bearer <"+Token+">")
	tok := struct {
		Token string
		Until time.Time
	}{Token: Token, Until: time.Now().Add(TokenExp)}
	rwr.WriteHeader(http.StatusOK) // 200 — пользователь успешно зарегистрирован и аутентифицирован;
	json.NewEncoder(rwr).Encode(tok)
}
