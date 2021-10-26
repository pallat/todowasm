package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/augustoroman/promise"
	"github.com/google/uuid"
)

const baseURL = "http://localhost:8081"

func PromiseToken(p *promise.Promise) {
	go func() {
		t, err := TokenAPI()
		if err != nil {
			p.Reject(err)
		} else {
			p.Resolve(t)
		}
	}()
}

func TokenAPI() (string, error) {
	type accessToken struct {
		Token string `json:"token"`
	}
	resp, err := http.Get(baseURL + "/tokenz")
	if err != nil {
		js.Global().Call("alert", err.Error())
		return "", err
	}

	var t accessToken
	json.NewDecoder(resp.Body).Decode(&t)

	return t.Token, nil
}

func PromiseTodoList(token string, p *promise.Promise) {
	go func() {
		t, err := TodoListAPI(token)
		if err != nil {
			p.Reject(err.Error())
		} else {
			p.Resolve(t)
		}
	}()
}

func TodoListAPI(token string) ([]Todo, error) {
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/todos", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	var t []Todo
	json.NewDecoder(resp.Body).Decode(&t)

	return t, nil
}

func PromiseAddTodo(token, text string, p *promise.Promise) {
	go func() {
		t, err := AddTodoAPI(token, text)
		if err != nil {
			p.Reject(err.Error())
			println(err.Error())
		} else {
			p.Resolve(t)
		}
	}()
}

func AddTodoAPI(token string, text string) (Todo, error) {
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(&Todo{Text: text})

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/todos", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("TransactionID", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return Todo{}, err
	}
	if resp.StatusCode >= 400 {
		return Todo{}, errors.New(resp.Status)
	}

	var t Todo
	json.NewDecoder(resp.Body).Decode(&t)

	return t, nil
}

func PromiseRemoveTodo(token string, id uint, p *promise.Promise) {
	go func() {
		err := RemoveTodoAPI(token, id)
		println("remove", err.Error())
		if err != nil {
			println("error:", err.Error())
			p.Reject(err.Error())
		} else {
			p.Resolve("")
		}
	}()
}

func RemoveTodoAPI(token string, id uint) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos/%d", baseURL, id), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("TransactionID", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}

	return nil
}
