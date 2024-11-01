/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:37:19
 * @FilePath: \go-middleware\tests\validator_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("TestVerify", TestVerify)
}

func TestVerify(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	rules := validator.Rules{
		"Name":  {validator.NotEmpty()},
		"Age":   {validator.Ge("18")},
		"Email": {validator.RegexpMatch(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)},
	}

	user := User{Name: "John Doe", Age: 25, Email: "john.doe@example.com"}
	err := validator.Verify(user, rules)
	assert.NoError(t, err)

	// 测试错误情况
	user = User{Name: "", Age: 16, Email: "invalid_email"}
	err = validator.Verify(user, rules)
	assert.Error(t, err)
}
