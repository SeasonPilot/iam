// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import "github.com/marmotedu/iam/internal/apiserver/config"

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg) // 创建 HTTP 和 GRPC 服务实例
	if err != nil {
		return err
	}

	return server.PrepareRun().Run() // 启动 HTTP 和 GRPC Web 服务
}
