package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

const form = `<html><body><form action="#" method="post" name="bar">
                    <input type="text" name="in"/>
                    <input type="text" name="in"/>
					<input type="submit" value="Submit"/>
             </form></html></body>`

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
	// 注意:如果没有调用 ParseForm 方法，下面无法获取表单的数据
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		//io.WriteString(w, form)
		t1, err := template.ParseFiles("/Users/carlos/go/src/awesomeProject/webtest/login.html")
		if err != nil {
			fmt.Println(err)
		}
		log.Println(t1.Execute(w, nil))
	} else {
		err := r.ParseForm()   // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
		if err != nil {
			// handle error http.Error() for example
			log.Fatal("ParseForm: ", err)
		}
		// 请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		name, psw, err := selInfo()
		if err != nil {
			panic(err)
		}
		if name == r.Form["username"][0] && psw == r.Form["password"][0] {
			io.WriteString(w, "login success!")
		}
	}
}

func selInfo() (name, psw string, err error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:23306)/test?parseTime=true&charset=utf8")
	checkErr(err)
	_selectInfoSQL := "SELECT `name`, psd from login WHERE id=1;"
	// 查询数据
	row := db.QueryRow(_selectInfoSQL)
	err = row.Scan(&name, &psw)
	checkErr(err)
	defer db.Close()
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)       // 设置访问的路由
	http.HandleFunc("/login", login)         // 设置访问的路由
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
