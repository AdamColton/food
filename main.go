package main

import (
	"bufio"
	"fmt"
	"github.com/adamcolton/err"
	"github.com/adamcolton/food/food"
	"os"
	"strconv"
	"strings"
)

var commands = map[string]func(args string){}

func init() {
	commands["help"] = func(string) {
		for cmd, _ := range commands {
			fmt.Println(cmd)
		}
	}
	commands["find"] = func(args string) {
		foods := food.Search(args).Foods()
		for _, f := range foods {
			fmt.Println(f)
		}
	}
	commands["exit"] = func(string) {
		//placeholder
	}
	commands["build"] = func(string) {
		fmt.Print("Building database...")
		e := food.PopulateDB()
		if e == nil {
			fmt.Println("Done")
		} else {
			fmt.Println("Error")
			fmt.Println(e)
		}
	}
	commands["show"] = func(args string) {
		id, e := strconv.Atoi(args)
		if e != nil {
			fmt.Println("Bad ID")
			return
		}

		foodDes := &food.FoodDes{
			Id: uint32(id),
		}
		foodDes.Get()
		if foodDes.LongDesc == "" {
			fmt.Println("Bad ID")
			return
		}
		fmt.Println(foodDes.Detailed())
	}

	commands["recipe"] = func(args string) {
		activeRecipe = &food.Recipe{
			Name: args,
		}
		activeRecipe.Get()
		activeRecipe.Save()
	}

	commands["print-table"] = func(string) {
		if activeRecipe == nil {
			fmt.Println("No recipe selected")
			return
		}
		fmt.Println(activeRecipe.Detailed())
	}

	commands["save"] = func(string) {
		if activeRecipe == nil {
			fmt.Println("No recipe selected")
			return
		}
		filename := strings.Replace(activeRecipe.Name, " ", "_", -1) + ".txt"
		f, e := os.Create(filename)
		err.Panic(e)
		defer f.Close()

		fmt.Fprint(f, activeRecipe.Detailed())
	}

	commands["all-recipes"] = func(string) {
		for _, r := range food.AllRecipes() {
			fmt.Println(r)
		}
	}

	commands["nutrients"] = func(string) {
		var disp = ""
		for _, n := range food.AllNutrients() {
			if n.Display {
				disp = "[X] "
			} else {
				disp = "[ ] "
			}
			fmt.Println(disp, n.Id, ": ", n.Name)
		}
	}

	commands["nutrient"] = func(args string) {
		lst := strings.Split(args, " ")
		if len(lst) < 2 {
			fmt.Print("Not enough args; expect 'nutrientId show/hide'")
			return
		}
		nutrientId, e := strconv.Atoi(lst[0])
		if e != nil {
			fmt.Println(e)
			return
		}
		display := true
		if lst[1] == "hide" {
			display = false
		} else if lst[1] != "show" {
			fmt.Println("Second arg must be either 'show' or 'hide'")
			return
		}
		nutrDef := food.NutrDef{
			Id: uint16(nutrientId),
		}
		nutrDef.Get()
		if nutrDef.Name == "" {
			fmt.Printf("%d is not a valid nutrient ID", nutrientId)
			return
		}
		nutrDef.Display = display
		nutrDef.Save()
	}

	commands["add"] = func(args string) {
		if activeRecipe == nil {
			fmt.Println("No recipe selected")
			return
		}
		lst := strings.Split(args, " ")
		if len(lst) < 2 {
			fmt.Print("Not enough args; expect 'foodId amount'")
			return
		}
		foodId, e := strconv.Atoi(lst[0])
		if e != nil {
			fmt.Println(e)
			return
		}
		amount, e := strconv.ParseFloat(lst[1], 32)
		if e != nil {
			fmt.Println(e)
			return
		}
		activeRecipe.Add(uint32(foodId), float32(amount))
	}

	commands["print"] = func(args string) {
		if activeRecipe == nil {
			fmt.Println("No recipe selected")
			return
		}
		fmt.Println(activeRecipe)
	}
}

var activeRecipe *food.Recipe

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		if activeRecipe != nil {
			fmt.Print(activeRecipe.Name)
		}
		fmt.Print("> ")
		lineBytes, _, _ := reader.ReadLine()
		input := strings.SplitN(string(lineBytes), " ", 2)
		cmdStr := input[0]
		if cmdStr == "exit" {
			break
		}
		args := ""
		if len(input) == 2 {
			args = input[1]
		}
		cmd, ok := commands[cmdStr]
		if !ok {
			cmd = commands["help"]
		}
		cmd(args)
	}
}
