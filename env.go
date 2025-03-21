package jenv

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/oarkflow/date"
	"gopkg.in/yaml.v3"

	"github.com/oarkflow/jenv/utils"
)

func UnmarshalJSON(jsonData []byte, cfg any) error {
	var rawMap map[string]any
	if err := json.Unmarshal(jsonData, &rawMap); err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}
	return populateFields(cfg, rawMap)
}

func UnmarshalYAML(yamlData []byte, cfg any) error {
	var rawMap map[string]any
	if err := yaml.Unmarshal(yamlData, &rawMap); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	return populateFields(cfg, rawMap)
}

func populateFields(cfg any, rawMap map[string]any) error {
	val := reflect.ValueOf(cfg).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if key == "" {
			key = strings.Split(field.Tag.Get("yaml"), ",")[0]
		}
		rawValue, exists := rawMap[key]
		if !exists {
			continue
		}
		if err := setFieldValue(val.Field(i), rawValue); err != nil {
			return fmt.Errorf("error setting field '%s': %v", field.Name, err)
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, rawValue any) error {
	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		val, err := getEnvValueInt(rawValue)
		if err != nil {
			return err
		}
		field.SetInt(int64(val))
	case reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			val, err := getEnvValueDuration(rawValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(val))
		} else {
			val, err := getEnvValueInt64(rawValue)
			if err != nil {
				return err
			}
			field.SetInt(val)
		}
	case reflect.Float32, reflect.Float64:
		val, err := getEnvValueFloat(rawValue)
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.String:
		field.SetString(getEnv(rawValue))
	case reflect.Bool:
		val, err := getEnvValueBool(rawValue)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Slice:
		if field.Type() == reflect.TypeOf([]byte{}) || field.Type() == reflect.TypeOf(json.RawMessage{}) {
			if rawValue != nil {
				bt, err := json.Marshal(rawValue)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(bt))
			}
		} else {
			rawSlice, ok := rawValue.([]any)
			if !ok {
				return fmt.Errorf("expected slice for field, got %T", rawValue)
			}
			slice := reflect.MakeSlice(field.Type(), len(rawSlice), len(rawSlice))
			for i := 0; i < len(rawSlice); i++ {
				if err := setFieldValue(slice.Index(i), rawSlice[i]); err != nil {
					return err
				}
			}
			field.Set(slice)
		}
	case reflect.Map:
		rawMap, ok := rawValue.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map for field, got %T", rawValue)
		}
		newMap := reflect.MakeMap(field.Type())
		for k, v := range rawMap {
			elem := reflect.New(field.Type().Elem()).Elem()
			if err := setFieldValue(elem, v); err != nil {
				return err
			}
			newMap.SetMapIndex(reflect.ValueOf(k), elem)
		}
		field.Set(newMap)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			val, err := getEnvValueTime(rawValue)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(val))
		} else {
			rawStructMap, ok := rawValue.(map[string]any)
			if !ok {
				return fmt.Errorf("expected struct map for field, got %T", rawValue)
			}
			if err := populateFields(field.Addr().Interface(), rawStructMap); err != nil {
				return err
			}
		}
	case reflect.Interface:
		if rawValue != nil {
			field.Set(reflect.ValueOf(rawValue))
		}
	default:
		return fmt.Errorf("unsupported field type: %s %v %v", field.Kind(), field, rawValue)
	}
	return nil
}

func getEnv(rawValue any) string {
	strValue := fmt.Sprintf("%v", rawValue)
	if strings.HasPrefix(strValue, "${") && strings.HasSuffix(strValue, "}") {
		envVar := strings.TrimSpace(strValue[2 : len(strValue)-1])
		parts := strings.SplitN(envVar, ":", 2)
		envValue := Getenv(parts[0])
		if envValue == "" && len(parts) > 1 {
			envValue = parts[1]
		}
		return strings.ReplaceAll(envValue, "'", "")
	}
	return strValue
}

func getEnvValueInt(rawValue any) (int, error) {
	val := getEnv(rawValue)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func getEnvValueInt64(rawValue any) (int64, error) {
	val := getEnv(rawValue)
	if val == "" {
		return 0, nil
	}
	return strconv.ParseInt(getEnv(rawValue), 10, 64)
}

func getEnvValueFloat(rawValue any) (float64, error) {
	val := getEnv(rawValue)
	if val == "" {
		return 0, nil
	}
	return strconv.ParseFloat(getEnv(rawValue), 64)
}

func getEnvValueBool(rawValue any) (bool, error) {
	val := getEnv(rawValue)
	if val == "" {
		return false, nil
	}
	return strconv.ParseBool(getEnv(rawValue))
}

func getEnvValueDuration(rawValue any) (time.Duration, error) {
	val := getEnv(rawValue)
	if val == "" {
		return 0, nil
	}
	return time.ParseDuration(getEnv(rawValue))
}

func getEnvValueTime(rawValue any) (time.Time, error) {
	val := getEnv(rawValue)
	if val == "" {
		return time.Time{}, nil // Return zero time if empty
	}
	switch rawValue := rawValue.(type) {
	case string:
		return date.Parse(getEnv(rawValue))
	case time.Time:
		return rawValue, nil
	}
	return time.Parse("2006-01-02T15:04:05Z07:00", getEnv(rawValue))
}

type GetEnvFn func(v string, defaultVal ...any) string

var Getenv GetEnvFn

func getenv(v string, defaultVal ...any) string {
	val := os.Getenv(v)
	if val != "" {
		return val
	}
	if len(defaultVal) > 0 && defaultVal[0] != nil {
		val, _ := utils.ToString(defaultVal[0])
		return val
	}
	return ""
}

func init() {
	Getenv = getenv
}
