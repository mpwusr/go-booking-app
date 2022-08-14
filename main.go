package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const conferenceTickets uint = 50
const timeDelay time.Duration = 10

var conferenceName string = "Go Conference"
var remainingTickets uint = conferenceTickets

type UserData struct {
	firstName              string
	lastName               string
	email                  string
	numberOfTickets        uint
	isOptedInForNewsletter bool
}

var bookings = make([]UserData, 0)
var wg = sync.WaitGroup{}

func main() {
	ticketsRemaining := remainingTickets > 0
	for ticketsRemaining && len(bookings) < 50 {
		greetUsers()
		firstName, lastName, email, userTickets, city := getUserInput()
		isValidName, isValidTicketNumber, isValidEmail, isValidCity := validateUserInput(firstName, lastName, email, userTickets, remainingTickets, city)
		if isValidName && isValidTicketNumber && isValidEmail && isValidCity {
			bookTicket(firstName, lastName, userTickets, email)
			var fTicket = createFormattedTicket(userTickets, firstName, lastName)
			wg.Add(1)
			go emailTicket(fTicket, email)
			printBookings(firstName, lastName, userTickets, email)
		} else {
			if !isValidTicketNumber {
				fmt.Printf("We only have %v remaining tickets; so you cannot book %v tickets\n", remainingTickets, userTickets)
				ticketsRemaining = remainingTickets > 0
				continue
			}
			if !isValidName {
				fmt.Printf("Please enter valid name; %v %v is not a valid name\n", firstName, lastName)
				continue
			}
			if !isValidEmail {
				fmt.Printf("Please enter valid email; %v is not a valid email\n", email)
				continue
			}
			if !isValidCity {
				fmt.Printf("please enter valid city; %v is not a valid city\n", city)
				continue
			}
		}
	}
	// end program
	wg.Wait()
	fmt.Println("Sorry our conference is booked up for this year. Come back next year.")
}
func greetUsers() {
	//clearScreen()
	fmt.Printf("Welcome to **%v** booking application\n", conferenceName)
	fmt.Printf("conferenceTickets is %v, remainTickets is %v, conferenceName is %v\n", conferenceTickets, remainingTickets, conferenceName)
}
func validateUserInput(fName string, lName string, email string, uTickets uint, remTickets uint, city string) (bool, bool, bool, bool) {
	isVName := len(fName) >= 2 && len(lName) >= 2
	isVEmail := strings.Contains(email, "@")
	isVTicketNumber := uTickets > 0 && uTickets <= remTickets
	isVCity := city == "Singapore" || city == "London"
	return isVName, isVEmail, isVTicketNumber, isVCity
}

func getFirstNames() []string {
	var firstNames []string
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}
func printBookings(fName string, lName string, uTickets uint, email string) {
	fmt.Printf("Thank You %v %v for booking %v tickets. You will receive a confirmation email at %v \n", fName, lName, uTickets, email)
}
func printTicketInfo() {
	fmt.Printf("%v tickets remaining for %v.\n", remainingTickets, conferenceName)
}
func getUserInput() (string, string, string, uint, string) {
	var fName string
	var lName string
	var email string
	var uTickets uint
	var city string
	city = "London"
	fmt.Println("Enter your first name: ")
	fmt.Scan(&fName)
	fmt.Println("Enter your last name: ")
	fmt.Scan(&lName)
	fmt.Println("Enter your email address: ")
	fmt.Scan(&email)
	fmt.Println("Enter number of tickets: ")
	fmt.Scan(&uTickets)
	return fName, lName, email, uTickets, city
}
func bookTicket(fName string, lName string, uTickets uint, email string) {
	const isOptedInForNewsletter = false
	var userData = UserData{
		fName,
		lName,
		email,
		uTickets,
		isOptedInForNewsletter,
	}
	bookings = append(bookings, userData)
	printBookings(fName, lName, uTickets, email)
	remainingTickets = remainingTickets - uTickets
	printTicketInfo()
	firstNames := getFirstNames()
	fmt.Printf("the first names of bookings %v are: %v\n", bookings, firstNames)
}
func createFormattedTicket(userTickets uint, firstName string, lastName string) string {
	var formattedTicket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	return formattedTicket
}
func emailTicket(fTicket string, email string) {
	fmt.Println("***********************************************")
	fmt.Printf("Sending formatted ticket to email address %v\n", email)
	fmt.Println("***********************************************")
	time.Sleep(timeDelay * time.Second)
	wg.Done()
}
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
