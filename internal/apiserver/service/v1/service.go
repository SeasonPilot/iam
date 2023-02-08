// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

//go:generate mockgen -self_package=github.com/marmotedu/iam/internal/apiserver/service/v1 -destination mock_service.go -package v1 github.com/marmotedu/iam/internal/apiserver/service/v1 Service,UserSrv,SecretSrv,PolicySrv

import "github.com/marmotedu/iam/internal/apiserver/store"

// 抽象工厂
// Service defines functions used to return resource interface.
type Service interface {
	Users() UserSrv //返回值为 产品接口
	Secrets() SecretSrv
	Policies() PolicySrv
}

// 具体产品工厂
// 实现抽象工厂接口，负责创建产品对象。
type service struct {
	store store.Factory
}

// 方法签名的返回值类型为接口。
// 相当于 Java 代码中 client 使用工厂时，创建具体工厂后赋值给抽象工厂,然后通过抽象工厂调用创建对象方法。 例： AnimalFactory f = new DogFactory();  Animal a = f.createAnimal(); https://www.zhihu.com/question/27125796#:~:text=AnimalFactory%20f%20%3D%20new%20DogFactory()%3B
// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

// 具体产品工厂,返回值类型的是 产品接口.
func (s *service) Users() UserSrv {
	return newUsers(s)
}

func (s *service) Secrets() SecretSrv {
	return newSecrets(s) // 目的是通过newSecrets传入不同的参数，可以创建不同的实例。  可以参考工厂模式的介绍
}

func (s *service) Policies() PolicySrv {
	return newPolicies(s)
}
