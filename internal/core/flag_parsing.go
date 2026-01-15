package core

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

// parseFlags uses reflection to parse cobra flags into a struct.
// The flagStruct parameter must be a pointer to a struct with `flag` tags.
func parseFlags(cmd *cobra.Command, flagStruct any) error {
	v := reflect.ValueOf(flagStruct)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("flagStruct must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		flagName := fieldType.Tag.Get("flag")

		if flagName == "" {
			continue
		}

		if !field.CanSet() {
			continue
		}

		if err := setFieldFromFlag(cmd, field, flagName); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func setFieldFromFlag(cmd *cobra.Command, field reflect.Value, flagName string) error {
	switch field.Kind() {
	case reflect.String:
		val, err := cmd.Flags().GetString(flagName)
		if err != nil {
			return err
		}
		field.SetString(val)

	case reflect.Bool:
		val, err := cmd.Flags().GetBool(flagName)
		if err != nil {
			return err
		}
		field.SetBool(val)

	case reflect.Int:
		val, err := cmd.Flags().GetInt(flagName)
		if err != nil {
			return err
		}
		field.SetInt(int64(val))

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			val, err := cmd.Flags().GetStringSlice(flagName)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(val))
		}
	}

	return nil
}
