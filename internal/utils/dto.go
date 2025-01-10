package utils

// Вспомогательная функция для копирования значений через указатели (разыменовывает ссылку и копирует значение)
func CopyIfNotNil[T any](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}