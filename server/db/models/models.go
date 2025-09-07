package models

var ModelsToMigrate = []any{}

func RegisterModel[T any](model *T) {
	ModelsToMigrate = append(ModelsToMigrate, model)
}
