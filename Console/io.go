package Console

import "fmt"

func Write(value string) (int, error)             { return fmt.Print(value) }
func Writeln(a ...any) (int, error)               { return fmt.Println(a...) }
func Writef(format string, a ...any) (int, error) { return fmt.Printf(format, a...) }
