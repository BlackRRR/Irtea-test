package security

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var (
	// sensitiveValue Используется для замены чувствительных данных.
	sensitiveValue = "Sensitive content *****"

	// sensitiveFieldsNames содержит список полей, которые считаются чувствительными.
	sensitiveFieldsNames = []string{
		"password",
		"auth2fa_code",
		"any2fa_code",
		"sessionId",
		"email",
		"login",
		"username",
		"authorization",
		"cookie",
		"newPassword",
		"keyPath",
		"mnemonics",
		"treasurerEmail",
		"userEmail",
		"cashierEmail",
		"nickname",
		"accessToken",
		"refreshToken",
		"_secret",
		"csrf-token",
		"user-agent",
	}
)

// HideSensitiveData Скрывает чувствительную информацию в предоставленных данных.
func HideSensitiveData(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		// Попробуем распарсить строку как JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(v), &jsonData); err == nil {
			return hideFields(jsonData)
		}
		// Если строка не JSON, проверяем на наличие чувствительных данных
		return hideStringSensitive(v)
	case error:
		// Обрабатываем интерфейс error, скрываем чувствительные данные в сообщении
		return errors.New(hideStringSensitive(v.Error()))
	case map[string]interface{}:
		return hideFields(v)
	default:
		return data
	}
}

// HideSensitiveField скрывает значение, если ключ является чувствительным.
func HideSensitiveField(key string, value interface{}) interface{} {
	if isSensitiveField(key) {
		return sensitiveValue
	}

	return HideSensitiveData(value)
}

// hideFields Скрывает чувствительные поля в объекте.
func hideFields(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		if isSensitiveField(key) {
			data[key] = sensitiveValue
			continue
		}

		subMap, ok := value.(map[string]interface{})
		if ok {
			data[key] = hideFields(subMap)

			continue
		}

		subStr, okString := value.(string)
		if !okString {
			continue
		}

		data[key] = hideStringSensitive(subStr)
	}

	return data
}

// hideStringSensitive Скрывает чувствительные данные в строке.
func hideStringSensitive(data string) string {
	for _, field := range sensitiveFieldsNames {
		regex := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(field))
		if !regex.MatchString(data) {
			continue
		}

		return sensitiveValue
	}
	return data
}

// isSensitiveField Проверяет, является ли поле чувствительным.
func isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	for _, sensitive := range sensitiveFieldsNames {
		if strings.ToLower(sensitive) != field {
			continue
		}

		return true
	}
	return false
}
