package main

//
//import "fmt"
//
//type IFather interface {
//	getName() string
//	setName(string)
//}
//
//type Person struct {
//	name string
//}
//
//func (p Person) getName() string {
//	return p.name
//}
//
//func (p *Person) setName(name string) {
//	p.name = name
//}
//
//func main() {
//	per := Person{}
//	fmt.Printf("未初始化默认值：per :%s\n", per.getName())
//	per.setName("test")
//	fmt.Printf("设置值后：per :%s\n", per.getName())
//
//	per2 := &Person{name: "lisi"}
//	fmt.Printf("初始化值后：per2 :%s\n", per2.getName())
//	per2.setName("wanger")
//	fmt.Printf("设置值后：per2 :%s\n", per2.getName())
//}

//package main
//
//import "fmt"
//
//type IFather interface {
//	getName() string
//	setName(string)
//}
//
//type Person struct {
//	name string
//}
//
//func (p Person) getName() string {
//	return p.name
//}
//
//func (p *Person) setName(name string) {
//	p.name = name
//}
//
//func main() {
//	var IPer IFather = &Person{}
//	var IPer2 IFather = Person{name: "lisi"}
//	// IPer := &Person{}
//	// IPer2 := Person{name: "lisi"}
//	fmt.Printf("未初始化默认值：s1:%s\n", IPer.getName())
//	IPer.setName("test")
//	fmt.Printf("设置值后：s1:%s\n", IPer.getName())
//}
