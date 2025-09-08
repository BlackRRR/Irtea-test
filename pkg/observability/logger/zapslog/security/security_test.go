package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHideSensitiveData_AllSensitiveFields(t *testing.T) {
	t.Run("Должен скрывать все чувствительные поля", func(t *testing.T) {
		input := map[string]interface{}{
			"password":        "mySecretPassword",
			"auth2fa_code":    "123456",
			"any2fa_code":     "654321",
			"sessionId":       "session123",
			"email":           "user@example.com",
			"login":           "userLogin",
			"username":        "userName",
			"authorization":   "Bearer token123",
			"cookie":          "someCookie",
			"newPassword":     "newSecretPassword",
			"keyPath":         "/path/to/key",
			"mnemonics":       "some mnemonics",
			"treasurerEmail":  "treasurer@example.com",
			"userEmail":       "user@example.com",
			"cashierEmail":    "cashier@example.com",
			"nickname":        "userNickname",
			"accessToken":     "accessToken123",
			"refreshToken":    "refreshToken123",
			"_secret":         "superSecret",
			"csrf-token":      "csrfToken123",
			"user-agent":      "Mozilla/5.0",
			"nonSensitiveKey": "nonSensitiveValue",
		}
		expected := map[string]interface{}{
			"password":        sensitiveValue,
			"auth2fa_code":    sensitiveValue,
			"any2fa_code":     sensitiveValue,
			"sessionId":       sensitiveValue,
			"email":           sensitiveValue,
			"login":           sensitiveValue,
			"username":        sensitiveValue,
			"authorization":   sensitiveValue,
			"cookie":          sensitiveValue,
			"newPassword":     sensitiveValue,
			"keyPath":         sensitiveValue,
			"mnemonics":       sensitiveValue,
			"treasurerEmail":  sensitiveValue,
			"userEmail":       sensitiveValue,
			"cashierEmail":    sensitiveValue,
			"nickname":        sensitiveValue,
			"accessToken":     sensitiveValue,
			"refreshToken":    sensitiveValue,
			"_secret":         sensitiveValue,
			"csrf-token":      sensitiveValue,
			"user-agent":      sensitiveValue,
			"nonSensitiveKey": "nonSensitiveValue",
		}

		result := HideSensitiveData(input).(map[string]interface{})
		assert.Equal(t, expected, result)
	})
}

func TestHideSensitiveData_SensitiveStrings(t *testing.T) {
	t.Run("Должен скрывать чувствительные строки, если они содержат ключевые слова", func(t *testing.T) {
		for _, field := range sensitiveFieldsNames {
			input := "Это строка содержит " + field + ": значение123"
			expected := sensitiveValue

			result := HideSensitiveData(input).(string)
			assert.Equal(t, expected, result, "Поле: "+field)
		}
	})

	t.Run("Не должен изменять строки без чувствительных данных", func(t *testing.T) {
		input := "Это строка без чувствительных данных"
		expected := "Это строка без чувствительных данных"

		result := HideSensitiveData(input).(string)
		assert.Equal(t, expected, result)
	})
}

func TestHideSensitiveData_NestedObjects(t *testing.T) {
	t.Run("Должен скрывать чувствительные данные во вложенных объектах", func(t *testing.T) {
		input := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"password": "mySecretPassword",
					"email":    "user@example.com",
				},
				"nonSensitiveKey": "nonSensitiveValue",
			},
		}
		expected := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"password": sensitiveValue,
					"email":    sensitiveValue,
				},
				"nonSensitiveKey": "nonSensitiveValue",
			},
		}

		result := HideSensitiveData(input).(map[string]interface{})
		assert.Equal(t, expected, result)
	})
}

func TestHideSensitiveData_UserAgentHeader(t *testing.T) {
	t.Run("Должен скрывать user-agent, если он присутствует", func(t *testing.T) {
		input := map[string]interface{}{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		}
		expected := map[string]interface{}{
			"user-agent": sensitiveValue,
		}

		result := HideSensitiveData(input).(map[string]interface{})
		assert.Equal(t, expected, result)
	})
}

func TestIsSensitiveField(t *testing.T) {
	t.Run("Должен возвращать true для всех чувствительных полей", func(t *testing.T) {
		for _, field := range sensitiveFieldsNames {
			assert.True(t, isSensitiveField(field), "Поле: "+field)
		}
	})

	t.Run("Должен возвращать false для нечувствительных полей", func(t *testing.T) {
		assert.False(t, isSensitiveField("nonSensitiveKey"))
		assert.False(t, isSensitiveField("randomField"))
	})
}
