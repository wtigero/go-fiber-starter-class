package container

import (
	"dig-di/controller"
	"dig-di/repository"
	"dig-di/service"

	"go.uber.org/dig"
)

// BuildContainer สร้างและ configure Dig container
func BuildContainer() *dig.Container {
	container := dig.New()

	// 🏗️ ลงทะเบียน Constructors ทุกตัว
	// Dig จะหาลำดับการสร้างเอง

	// 1. Repository layer (ไม่มี dependencies)
	if err := container.Provide(repository.NewUserRepository); err != nil {
		panic("Failed to provide UserRepository: " + err.Error())
	}

	// 2. Service layer (ต้องการ UserRepository)
	if err := container.Provide(service.NewUserService); err != nil {
		panic("Failed to provide UserService: " + err.Error())
	}

	// 3. Controller layer (ต้องการ UserService)
	if err := container.Provide(controller.NewUserController); err != nil {
		panic("Failed to provide UserController: " + err.Error())
	}

	return container
}

// จำเป็นต้องเพิ่ม struct สำหรับ Fiber app ถ้าต้องการ inject
type FiberApp struct {
	UserController *controller.UserController
}
