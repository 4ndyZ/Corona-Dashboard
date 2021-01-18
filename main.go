package main


func main() {
	a := App{}

	a.Initialize("http://server:8086", "user", "pw")

	a.Run()
}
