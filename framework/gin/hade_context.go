// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gin

import (
	"context"
)

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}

// context 实现container的几个封装

// 实现make的封装
func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.container.Make(key)
}

// 实现mustGetInstance的封装
func (ctx *Context) MustGetInstance(key string) interface{} {
	return ctx.container.MustGetInstance(key)
}

// 实现makenew的封装
func (ctx *Context) MakeNew(key string) (interface{}, error) {
	return ctx.container.MakeNew(key)
}
