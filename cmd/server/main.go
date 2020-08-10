package main

func main() {
	a := App{}
	a.SetUp("./exercise.db")
	a.Run(":8010")
}
