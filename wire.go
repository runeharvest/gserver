//go:build wireinject

package gserver

// func NewLoginService(storager storage.Storager) *login.LoginService {
// 	return &login.LoginService{Storager: storager}
// }

// var LoginServiceSet = wire.NewSet(NewLoginService, wire.Bind(new(storage.Storager), new(*memory.MemoryStorage)), memory.NewMemoryStorage)

// func InitializeLoginService() (*login.LoginService, error) {
// 	wire.Build(LoginServiceSet)
// 	return &login.LoginService{}, nil
// }

// func NewMemoryStorage() *memory.MemoryStorage {
// 	return memory.NewMemoryStorage()
// }
