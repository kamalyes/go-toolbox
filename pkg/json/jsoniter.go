//go:build jsoniter

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 18:15:36
 * @FilePath: \go-middleware\json\jsoniter.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package json

import jsoniter "github.com/json-iterator/go"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by go-toolbox/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by go-toolbox/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by go-toolbox/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by go-toolbox/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by go-toolbox/json package.
	NewEncoder = json.NewEncoder
)
