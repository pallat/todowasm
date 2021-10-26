package todo

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/augustoroman/promise"
	dom "honnef.co/go/js/dom/v2"
)

type Todo struct {
	ID   uint
	Text string
}

type TodoApp struct {
	token    string
	todos    []Todo
	d        dom.Document
	body     *dom.HTMLBodyElement
	credDiv  *dom.HTMLDivElement
	todo     *dom.BasicHTMLElement
	loginBtn *dom.HTMLButtonElement
	newtodo  *dom.HTMLInputElement
	todolist *dom.HTMLUListElement
	count    *dom.HTMLSpanElement

	todoLI []*dom.HTMLLIElement
}

func New(d dom.Document) *TodoApp {

	body := d.GetElementByID("bodypage").(*dom.HTMLBodyElement)
	credDiv := d.GetElementByID("credential").(*dom.HTMLDivElement)
	todo := d.GetElementByID("todoapp").(*dom.BasicHTMLElement)

	login := d.GetElementByID("login").(*dom.HTMLButtonElement)
	newtodo := d.GetElementByID("new-todo").(*dom.HTMLInputElement)
	todolist := d.GetElementByID("todo-list").(*dom.HTMLUListElement)
	count := d.GetElementByID("todo-count").(*dom.HTMLSpanElement)

	todoapp := &TodoApp{
		todos:    []Todo{},
		d:        d,
		body:     body,
		credDiv:  credDiv,
		todo:     todo,
		loginBtn: login,
		newtodo:  newtodo,
		todolist: todolist,
		count:    count,
		todoLI:   []*dom.HTMLLIElement{},
	}

	login.AddEventListener("click", false, todoapp.LoginBtnClickEvent)
	newtodo.AddEventListener("keyup", false, todoapp.AddTodoEvent)

	return todoapp
}

func (t *TodoApp) LoginBtnClickEvent(event dom.Event) {
	var p = &promise.Promise{}

	p.Then(
		func(value interface{}) interface{} {
			t.token = value.(string)
			t.todo.Class().Remove("invisible")
			t.body.RemoveChild(t.credDiv)
			t.FetchTodoList()
			return p
		}, func(value interface{}) interface{} {
			println("error", value)
			js.Global().Call("alert", value)
			return p
		},
	)

	PromiseToken(p)
}

func (t *TodoApp) FetchTodoList() {
	var p = &promise.Promise{}

	p.Then(
		func(value interface{}) interface{} {
			t.todos = value.([]Todo)
			t.refreshTodoList()
			return p
		}, func(value interface{}) interface{} {
			println("error", value)
			return p
		},
	)

	PromiseTodoList(t.token, p)
}

func (t *TodoApp) AddTodo(val string) {
	var p = &promise.Promise{}

	p.Then(
		func(value interface{}) interface{} {
			newToDo := Todo{ID: 0, Text: val}
			t.todos = append(t.todos, newToDo)
			t.refreshTodoList()

			return p
		}, func(value interface{}) interface{} {
			println("error", value)
			js.Global().Call("alert", value)
			return p
		},
	)
	PromiseAddTodo(t.token, val, p)

}

func (t *TodoApp) AddTodoEvent(event dom.Event) {
	ke := event.(*dom.KeyboardEvent)
	if ke.KeyCode() == 13 {
		input := event.Target().(*dom.HTMLInputElement)
		if input.Value() == "" {
			return
		}
		t.AddTodo(input.Value())
		input.SetValue("")
	}
}

func (t *TodoApp) RemoveTodoEvent(event dom.Event) {
	var p = &promise.Promise{}

	input := event.Target().(*dom.HTMLButtonElement)
	println("remove at", input.Value())

	id, err := strconv.Atoi(input.Value())
	if err != nil {
		println("error", err.Error())
		return
	}

	p.Then(
		func(value interface{}) interface{} {
			t.FetchTodoList()
			return p
		}, func(value interface{}) interface{} {
			println("error", value)
			js.Global().Call("alert", value)
			return p
		},
	)

	PromiseRemoveTodo(t.token, uint(id), p)
}

func (t *TodoApp) refreshTodoList() {
	for _, li := range t.todoLI {
		t.todolist.RemoveChild(li)
	}

	t.todoLI = []*dom.HTMLLIElement{}

	for _, todo := range t.todos {
		li := t.d.CreateElement("li").(*dom.HTMLLIElement)

		div := t.d.CreateElement("div").(*dom.HTMLDivElement)
		div.SetClass("view")
		cb := t.d.CreateElement("input").(*dom.HTMLInputElement)
		cb.SetType("checkbox")
		cb.SetClass("toggle")
		lb := t.d.CreateElement("label").(*dom.HTMLLabelElement)
		lb.SetInnerHTML(todo.Text)
		btn := t.d.CreateElement("button").(*dom.HTMLButtonElement)
		btn.SetClass("destroy")
		btn.SetValue(strconv.Itoa(int(todo.ID)))
		btn.AddEventListener("click", false, t.RemoveTodoEvent)

		div.AppendChild(cb)
		div.AppendChild(lb)
		div.AppendChild(btn)

		inp := t.d.CreateElement("input").(*dom.HTMLInputElement)
		inp.SetClass("edit")

		li.AppendChild(div)
		li.AppendChild(inp)
		li.SetClass("todo-list")
		li.SetAttribute("data-id", strconv.Itoa(int(todo.ID)))

		t.todolist.AppendChild(li)
		t.todoLI = append(t.todoLI, li)

		t.refreshFooter()
	}
}

func (t *TodoApp) refreshFooter() {
	t.count.SetInnerHTML(fmt.Sprintf("<strong>%d</strong> left", len(t.todos)))
}
