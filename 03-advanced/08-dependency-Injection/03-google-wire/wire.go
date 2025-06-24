//go:build wireinject
// +build wireinject

package main

import (
	"wire-di/controller"
	"wire-di/repository"
	"wire-di/service"

	"github.com/google/wire"
)

// InitializeUserController ⚡ Wire จะ generate โค้ดจริงให้
//
//go:generate wire
func InitializeUserController() *controller.UserController {
	wire.Build(
		repository.NewUserRepository, // สร้าง UserRepository
		service.NewUserService,       // สร้าง UserService (รับ UserRepository)
		controller.NewUserController, // สร้าง UserController (รับ UserService)
	)
	return nil // Wire จะแทนที่บรรทัดนี้ด้วยโค้ดจริง
}
