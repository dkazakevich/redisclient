package main

import "fmt"

func main() {

	client := newClient("http://localhost:8080/api/v1/")

	client.put("list", [6]int{2, 3, 5, 7, 11, 13})
	fmt.Println(client.get("list"))
	fmt.Println(client.getListElement("list", 0))

	client.put("planets", map[string]string{"planet1": "Mercury", "planet2": "Venus", "planet3": "Earth"})
	fmt.Println(client.get("planets"))
	fmt.Println(client.getDictValue("planets", "planet3"))
	fmt.Println(client.getDictValue("planets", "planet4"))

	fmt.Println(client.keys())
}
