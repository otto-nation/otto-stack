package core

import (
	"reflect"

	"github.com/spf13/cobra"
)

// parseFlags uses reflection to parse cobra flags into a struct
func parseFlags(cmd *cobra.Command, flagStruct any) {
	v := reflect.ValueOf(flagStruct).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		flagName := fieldType.Tag.Get("flag")

		if flagName == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if val, err := cmd.Flags().GetString(flagName); err == nil {
				field.SetString(val)
			}
		case reflect.Bool:
			if val, err := cmd.Flags().GetBool(flagName); err == nil {
				field.SetBool(val)
			}
		case reflect.Int:
			if val, err := cmd.Flags().GetInt(flagName); err == nil {
				field.SetInt(int64(val))
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				if val, err := cmd.Flags().GetStringSlice(flagName); err == nil {
					field.Set(reflect.ValueOf(val))
				}
			}
		}
	}
}
