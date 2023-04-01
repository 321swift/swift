package utils

import "fmt"

func WelcomeScreen() int64 {
	var role int64

	for i := false; !i; {
		fmt.Println("")
		fmt.Println("Welcome to the swift application")
		fmt.Println("Select the role you wish to assume: \n\t 1. Sending \n\t 2. Receiving")

		var number int64
		_, err := fmt.Scanf("%d", &number)

		if err != nil {
			fmt.Println(err)
		} else if number < 0 {
			fmt.Println("Please enter one number only acording to the given options")
		} else if number > 2 {
			fmt.Println("Please enter one number only acording to the given options")
		} else {
			i = true
			role = number
		}
	}
	return role
}
