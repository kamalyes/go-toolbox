/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:37:19
 * @FilePath: \go-middleware\validator\validator_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

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

	rules := Rules{
		"Name":  {NotEmpty()},
		"Age":   {Ge("18")},
		"Email": {RegexpMatch(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)},
	}

	user := User{Name: "John Doe", Age: 25, Email: "john.doe@example.com"}
	err := Verify(user, rules)
	assert.NoError(t, err)

	// 测试错误情况
	user = User{Name: "", Age: 16, Email: "invalid_email"}
	err = Verify(user, rules)
	assert.Error(t, err)
}
