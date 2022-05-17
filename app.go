package main

import (
	"context"
	"fmt"

	backend "github.com/i-mes/imes-app/backend"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// wails 通过这里将上下文回吐给 user
	// user 后续用该上下文(a.ctx) 与 wails runtime 交互
	a.ctx = ctx
	backend.ImesBind(&ctx)
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

type Person struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

// Greet returns a greeting for the given name
func (a *App) Greet(p Person) string {
	fmt.Println("receive frontend greet message")
	return fmt.Sprintf("Hello %s, age: %d!", p.Name, p.Age)
}
