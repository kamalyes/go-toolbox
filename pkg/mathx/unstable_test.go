/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\mathx\unstable_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnstableAroundDuration(t *testing.T) {
	unstable := NewUnstable(0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundDuration(time.Second)
		assert.True(t, float64(time.Second)*0.95 <= float64(val))
		assert.True(t, float64(val) <= float64(time.Second)*1.05)
	}
}

func TestUnstableAroundInt(t *testing.T) {
	const target = 10000
	unstable := NewUnstable(0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.True(t, float64(target)*0.95 <= float64(val))
		assert.True(t, float64(val) <= float64(target)*1.05)
	}
}

func TestUnstableAroundIntLarge(t *testing.T) {
	const target int64 = 10000
	unstable := NewUnstable(5)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.True(t, 0 <= val)
		assert.True(t, val <= 2*target)
	}
}

func TestUnstableAroundIntNegative(t *testing.T) {
	const target int64 = 10000
	unstable := NewUnstable(-0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.Equal(t, target, val)
	}
}

func TestUnstableDistribution(t *testing.T) {
	const (
		seconds = 10000
		total   = 10000
	)

	m := make(map[int]int)
	expiry := NewUnstable(0.05)
	for i := 0; i < total; i++ {
		val := int(expiry.AroundInt(seconds))
		m[val]++
	}

	_, ok := m[0]
	assert.False(t, ok)

	mi := make(map[any]int, len(m))
	for k, v := range m {
		mi[k] = v
	}
	entropy := CalcEntropy(mi)
	assert.True(t, len(m) > 1)
	assert.True(t, entropy > 0.95)
}

// Address 结构体
type AddressClone struct {
	City  string
	Zip   string
	Extra map[string]string
}

// Person 结构体，包含递归引用
type PersonClone struct {
	Name    string
	Age     int
	Friends []*PersonClone
	Company *CompanyClone
}

// Company 结构体
type CompanyClone struct {
	Name      string
	Address   AddressClone
	Employees []*PersonClone
}

func TestCloneComplex(t *testing.T) {

	// 创建复杂的结构
	alice := &PersonClone{
		Name: "Alice",
		Age:  30,
		Company: &CompanyClone{
			Name: "Wonderland Inc.",
			Address: AddressClone{
				City: "Wonderland",
				Zip:  "12345",
				Extra: map[string]string{
					"Country": "Wonderland",
				},
			},
			Employees: []*PersonClone{},
		},
		Friends: []*PersonClone{},
	}

	bob := &PersonClone{
		Name: "Bob",
		Age:  25,
		Company: &CompanyClone{
			Name: "Builder Co.",
			Address: AddressClone{
				City: "Builderland",
				Zip:  "54321",
				Extra: map[string]string{
					"Country": "Builderland",
				},
			},
			Employees: []*PersonClone{},
		},
		Friends: []*PersonClone{alice},
	}

	alice.Friends = append(alice.Friends, bob)

	// 深拷贝
	clonedAlice := Clone(alice, nil).(*PersonClone)

	// 修改原始结构以验证克隆是否成功
	alice.Name = "Alice Updated"
	alice.Company.Name = "Wonderland Enterprises"
	alice.Company.Address.City = "New Wonderland"
	alice.Company.Address.Extra["Country"] = "New Wonderland"
	alice.Friends[0].Name = "Bob Updated"

	// 使用 assert 验证克隆的内容
	assert.Equal(t, "Alice", clonedAlice.Name, "Cloned Alice's name should remain 'Alice'")
	assert.Equal(t, 30, clonedAlice.Age, "Cloned Alice's age should remain 30")
	assert.Equal(t, "Wonderland Inc.", clonedAlice.Company.Name, "Cloned Alice's company name should remain 'Wonderland Inc.'")
	assert.Equal(t, "Wonderland", clonedAlice.Company.Address.City, "Cloned Alice's company address city should remain 'Wonderland'")
	assert.Equal(t, "Wonderland", clonedAlice.Company.Address.Extra["Country"], "Cloned Alice's company address extra country should remain 'Wonderland'")
	assert.Equal(t, 1, len(clonedAlice.Friends), "Cloned Alice should have 1 friend")
	assert.Equal(t, "Bob", clonedAlice.Friends[0].Name, "Cloned Alice's friend's name should remain 'Bob'")

	// 验证 Bob 的朋友是否未被修改
	assert.Equal(t, "Bob", clonedAlice.Friends[0].Name, "Cloned Bob's name should remain 'Bob'")
	assert.Equal(t, "Builder Co.", clonedAlice.Friends[0].Company.Name, "Cloned Bob's company name should remain 'Builder Co.'")

	// 验证修改原始结构不影响克隆
	assert.Equal(t, "Alice Updated", alice.Name, "Original Alice's name should be updated")
	assert.Equal(t, "Bob", clonedAlice.Friends[0].Name, "Cloned Bob's name should remain 'Bob' after original is updated")
}
