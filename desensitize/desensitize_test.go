/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 19:32:52
 * @FilePath: \go-toolbox\desensitize\desensitize_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

import (
	"fmt"
	"testing"

	"github.com/kamalyes/go-toolbox/stringx"
)

func TestDesensitize(t *testing.T) {
	fmt.Printf("%#v\n", Desensitize("17603393007", MobilePhone))
	fmt.Printf("%#v\n", Desensitize("石晓浩", ChineseName))
	fmt.Printf("%#v\n", Desensitize("李赵一特", ChineseName))

	fmt.Printf("%#v\n", Desensitize("131182199702251631", IDCard))
	fmt.Printf("%#v\n", Desensitize("北京市海淀区清华东路西口29#1#601", Address))
	fmt.Printf("%#v\n", stringx.Length("北京市海淀区清华东路西口29#1#601"))
	fmt.Printf("%#v\n", Desensitize("北京市海淀区清华东路西口29#1#601", Email))
	fmt.Printf("%#v\n", Desensitize("12313123123", Password))
	fmt.Printf("%#v\n", Desensitize("苏D40000", CarLicense))
	fmt.Printf("%#v\n", Desensitize("陕A12345D", CarLicense))
	fmt.Printf("%#v\n", Desensitize("1234    2222 3333 4444 6789 9", BankCard))
	fmt.Printf("%#v\n", Desensitize("1234 2222 3333 4444 6789 91", BankCard))
	fmt.Printf("%#v\n", Desensitize("1234 2222 3333 4444 678", BankCard))
	fmt.Printf("%#v\n", Desensitize("1234 2222 3333 4444 6789", BankCard))
	fmt.Printf("%#v\n", Desensitize("127.0.0.1", IPV4))
	fmt.Printf("%#v\n", Desensitize("2001:0db8:86a3:08d3:1319:8a2e:0370:7344", IPV6))
}
