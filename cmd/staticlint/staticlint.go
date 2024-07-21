// Package main implements a static code analysis tool that checks
// for the use of os.Exit in the main function. It combines standard
// static analysis checks provided by golang.org/x/tools/go/analysis/passes
// and honnef.co/go/tools/staticcheck with a custom analyzer that forbids
// the use of os.Exit in the main package's main function.
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"honnef.co/go/tools/staticcheck"
)

// main initializes and runs the multichecker with the specified analyzers,
// including the custom noosexit analyzer which checks for os.Exit calls in the main function.
func main() {
	// Standard analyzers from golang.org/x/tools/go/analysis/passes
	analyzers := []*analysis.Analyzer{
		atomicalign.Analyzer,
		bools.Analyzer,
		ctrlflow.Analyzer,
	}

	// Append all SA class analyzers from staticcheck.io
	for _, a := range staticcheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// myAnalyzer is a custom analyzer added to discourage the use of os.Exit in the main function.
	analyzers = append(analyzers, myAnalyzer)

	// multichecker.Main runs the analysis.
	multichecker.Main(analyzers...)
}

// myAnalyzer defines the custom analyzer looking for os.Exit calls within the main function.
// It is part of the suite of analyzers run by the multichecker.
var myAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "forbids the use of os.Exit in the main function",
	Run:  run,
}

// run is the entry point for the analyzer. It iterates over files in the package
// and delegates the inspection of the main function to checkForOsExitInMain.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkForOsExitInMain(pass, file)
	}
	return nil, nil
}

// checkForOsExitInMain inspects the provided file for a main function and,
// if found, examines it for calls to os.Exit.
func checkForOsExitInMain(pass *analysis.Pass, file *ast.File) {
	// Implementation omitted for brevity.
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true // Not a function, continuing traversal

		}
		if funcDecl.Name.Name != "main" {
			return true // Not the main function, continuing traversal
		}
		inspectMainFuncForOsExit(pass, funcDecl)
		return false // No need to traverse inner nodes of the main function, as we have already checked its body
	})
}

// inspectMainFuncForOsExit checks if the provided function declaration (assumed to be 'main')
// contains calls to os.Exit, reporting them if found.
func inspectMainFuncForOsExit(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Implementation omitted for brevity.
	for _, stmt := range funcDecl.Body.List {
		callExpr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue // Not a function call, moving to the next statement
		}
		call, ok := callExpr.X.(*ast.CallExpr)
		if !ok {
			continue // Not a call expression, moving to the next statement
		}
		if isOsExitCall(call) {
			pass.Reportf(call.Pos(), "call to os.Exit found in main function")
		}
	}
}

// isOsExitCall determines whether the given call expression is a call to os.Exit.
func isOsExitCall(call *ast.CallExpr) bool {
	// Implementation omitted for brevity.
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false // Not a method call, therefore not os.Exit
	}
	ident, ok := selExpr.X.(*ast.Ident)
	return ok && ident.Name == "os" && selExpr.Sel.Name == "Exit"
}
