/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 13:15:40
 * @FilePath: \go-toolbox\pkg\safe\nil_panic_detector.go
 * @Description: Nil Panic检测工具，帮助发现项目中可能的nil指针访问
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// NilPanicDetector Nil Panic检测器
type NilPanicDetector struct {
	fileSet      *token.FileSet
	issues       []NilPanicIssue
	riskPatterns []RiskPattern
}

// NilPanicIssue Nil Panic问题
type NilPanicIssue struct {
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Code        string `json:"code"`
}

// RiskPattern 风险模式
type RiskPattern struct {
	Name        string
	Pattern     string
	Description string
	Severity    string
}

// NewNilPanicDetector 创建新的检测器
func NewNilPanicDetector() *NilPanicDetector {
	return &NilPanicDetector{
		fileSet: token.NewFileSet(),
		issues:  make([]NilPanicIssue, 0),
		riskPatterns: []RiskPattern{
			{
				Name:        "NestedFieldAccess",
				Pattern:     "x.y.z",
				Description: "嵌套字段访问可能导致nil panic",
				Severity:    "HIGH",
			},
			{
				Name:        "PointerDereference",
				Pattern:     "*ptr",
				Description: "指针解引用可能导致nil panic",
				Severity:    "HIGH",
			},
			{
				Name:        "SliceIndexAccess",
				Pattern:     "slice[i]",
				Description: "切片索引访问可能导致越界",
				Severity:    "MEDIUM",
			},
			{
				Name:        "MapAccess",
				Pattern:     "map[key]",
				Description: "Map访问需要检查ok值",
				Severity:    "MEDIUM",
			},
			{
				Name:        "InterfaceAssertion",
				Pattern:     "v.(Type)",
				Description: "类型断言可能失败",
				Severity:    "MEDIUM",
			},
		},
	}
}

// ScanDirectory 扫描目录中的Go文件
func (d *NilPanicDetector) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理.go文件
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "_test.go") {
			if err := d.ScanFile(path); err != nil {
				fmt.Printf("扫描文件 %s 时出错: %v\n", path, err)
			}
		}
		return nil
	})
}

// ScanFile 扫描单个文件
func (d *NilPanicDetector) ScanFile(filePath string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(d.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	// 访问AST节点
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			d.checkSelectorExpr(filePath, node)
		case *ast.StarExpr:
			d.checkStarExpr(filePath, node)
		case *ast.IndexExpr:
			d.checkIndexExpr(filePath, node)
		case *ast.TypeAssertExpr:
			d.checkTypeAssertExpr(filePath, node)
		}
		return true
	})

	return nil
}

// checkSelectorExpr 检查选择器表达式 (x.y.z)
func (d *NilPanicDetector) checkSelectorExpr(filePath string, expr *ast.SelectorExpr) {
	// 检查嵌套的选择器表达式
	if sel, ok := expr.X.(*ast.SelectorExpr); ok {
		pos := d.fileSet.Position(expr.Pos())

		// 检查是否是多层嵌套
		depth := d.getSelectorDepth(expr)
		if depth >= 3 {
			d.addIssue(NilPanicIssue{
				File:        filePath,
				Line:        pos.Line,
				Column:      pos.Column,
				Type:        "NestedFieldAccess",
				Description: fmt.Sprintf("深度为%d的嵌套字段访问，建议使用安全访问模式", depth),
				Severity:    "HIGH",
				Code:        d.getCodeSnippet(filePath, pos.Line),
			})
		}

		// 递归检查内层选择器
		d.checkSelectorExpr(filePath, sel)
	}
}

// checkStarExpr 检查指针解引用
func (d *NilPanicDetector) checkStarExpr(filePath string, expr *ast.StarExpr) {
	pos := d.fileSet.Position(expr.Pos())

	// 检查是否有nil检查
	if !d.hasNilCheck(expr) {
		d.addIssue(NilPanicIssue{
			File:        filePath,
			Line:        pos.Line,
			Column:      pos.Column,
			Type:        "PointerDereference",
			Description: "指针解引用没有nil检查,可能导致panic",
			Severity:    "HIGH",
			Code:        d.getCodeSnippet(filePath, pos.Line),
		})
	}
}

// checkIndexExpr 检查索引表达式
func (d *NilPanicDetector) checkIndexExpr(filePath string, expr *ast.IndexExpr) {
	pos := d.fileSet.Position(expr.Pos())

	d.addIssue(NilPanicIssue{
		File:        filePath,
		Line:        pos.Line,
		Column:      pos.Column,
		Type:        "IndexAccess",
		Description: "索引访问可能导致越界，建议添加边界检查",
		Severity:    "MEDIUM",
		Code:        d.getCodeSnippet(filePath, pos.Line),
	})
}

// checkTypeAssertExpr 检查类型断言
func (d *NilPanicDetector) checkTypeAssertExpr(filePath string, expr *ast.TypeAssertExpr) {
	pos := d.fileSet.Position(expr.Pos())

	d.addIssue(NilPanicIssue{
		File:        filePath,
		Line:        pos.Line,
		Column:      pos.Column,
		Type:        "TypeAssertion",
		Description: "类型断言可能失败，建议使用 v, ok := x.(Type) 形式",
		Severity:    "MEDIUM",
		Code:        d.getCodeSnippet(filePath, pos.Line),
	})
}

// getSelectorDepth 获取选择器深度
func (d *NilPanicDetector) getSelectorDepth(expr *ast.SelectorExpr) int {
	depth := 1
	if sel, ok := expr.X.(*ast.SelectorExpr); ok {
		depth += d.getSelectorDepth(sel)
	}
	return depth
}

// hasNilCheck 检查是否有nil检查（简单的启发式检查）
func (d *NilPanicDetector) hasNilCheck(expr ast.Expr) bool {
	// 这是一个简化的检查，实际实现会更复杂
	return false
}

// getCodeSnippet 获取代码片段
func (d *NilPanicDetector) getCodeSnippet(filePath string, line int) string {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(src), "\n")
	if line <= 0 || line > len(lines) {
		return ""
	}

	return strings.TrimSpace(lines[line-1])
}

// addIssue 添加问题
func (d *NilPanicDetector) addIssue(issue NilPanicIssue) {
	d.issues = append(d.issues, issue)
}

// GetIssues 获取所有检测到的问题
func (d *NilPanicDetector) GetIssues() []NilPanicIssue {
	return d.issues
}

// GenerateReport 生成报告
func (d *NilPanicDetector) GenerateReport() string {
	var report strings.Builder

	report.WriteString("🔍 Nil Panic 检测报告\n")
	report.WriteString("========================\n\n")

	highCount := 0
	mediumCount := 0

	for _, issue := range d.issues {
		switch issue.Severity {
		case "HIGH":
			highCount++
		case "MEDIUM":
			mediumCount++
		}

		report.WriteString(fmt.Sprintf("📍 %s:%d:%d\n", issue.File, issue.Line, issue.Column))
		report.WriteString(fmt.Sprintf("   类型: %s (%s)\n", issue.Type, issue.Severity))
		report.WriteString(fmt.Sprintf("   描述: %s\n", issue.Description))
		report.WriteString(fmt.Sprintf("   代码: %s\n\n", issue.Code))
	}

	report.WriteString(fmt.Sprintf("总计: %d 个问题 (高风险: %d, 中风险: %d)\n",
		len(d.issues), highCount, mediumCount))

	return report.String()
}

// GetFixSuggestions 获取修复建议
func (d *NilPanicDetector) GetFixSuggestions() []string {
	suggestions := []string{
		"1. 使用 goconfig.New() 进行安全访问",
		"2. 使用 goconfig.SafeConfig() 进行配置专用安全访问",
		"3. 在指针解引用前添加 nil 检查",
		"4. 使用 v, ok := m[key] 形式进行map访问",
		"5. 使用 v, ok := x.(Type) 形式进行类型断言",
		"6. 添加边界检查后再进行切片/数组访问",
		"7. 考虑使用 Optional 类型或者 Maybe 模式",
	}
	return suggestions
}
