package verify_params

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

// 校验
func Verify(stru interface{}) (string, error) {
	validate := validator.New()
	uniTrans := ut.New(zh.New())
	trans, _ := uniTrans.GetTranslator("zh")
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(fmt.Sprintf("注册翻译器错误: %s\n", err.Error()))
	}
	// 验证
	err = validate.Struct(stru)
	var errString string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errString = errorMsgTrans(stru, err)
			if errString == "" {
				errString = err.Translate(trans)
			}
			return errString, nil
		}
	}
	return "", nil
}

// 错误信息转换
func errorMsgTrans(stru interface{}, err validator.FieldError) string {
	field := err.Field()
	typeOf := reflect.TypeOf(stru)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	fieldErrMsg, ok := typeOf.FieldByName(field)
	if ok {
		return fieldErrMsg.Tag.Get("valid_msg")
	}
	return ""
}
