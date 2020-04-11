package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

const (
	rightAnswer = "RIGHT_ANSWER"
	wrongAnswer = "WRONG_ANSWER"
)

func main() {

	printInstructions()
	var filename *string
	var timeLimit *int
	timeLimit = flag.Int("timeLimit", 30, "This flag is used to set the test time limit")
	filename = flag.String("file", "problems.csv", "This flag is used to pass the csv filename path")
	flag.Parse()
	absPath, _ := filepath.Abs(*filename)
	csvFile, err := os.Open(absPath)
	r := csv.NewReader(csvFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	var questions []string
	var answers []string
	for {
		result, _ := r.Read()
		if result == nil {
			break
		}
		questions = append(questions, result[0])
		answers = append(answers, result[1])
	}
	reader := bufio.NewScanner(os.Stdin)
	color.Yellow("Do you want to start the quiz")
	color.Yellow("press y/Y to start or n/N to stop")
	for reader.Scan() {
		readString := reader.Text()
		if readString == "Y" || readString == "y" {
			break
		} else {
			color.Red("OOh you decided to cancel the quiz bye....")
			os.Exit(1)
		}
	}
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correctAnswers := make(map[int]string)
	wrongAnswers := make(map[int]string)
	userAnswerChannel := make(chan string)
	//startTime := time.Now()
Loop:
	for i := range questions {
		fmt.Println("Question ", i+1, " : ", questions[i])
		go takeAnswer(reader, userAnswerChannel)
		select {
		case answerGiven := <-userAnswerChannel:
			if answerGiven == answers[i] {
				correctAnswers[i] = answerGiven
			} else {
				wrongAnswers[i] = answerGiven
			}
		case <-timer.C:
			var c *color.Color
			c = color.New(color.FgRed)
			c.Println(fmt.Sprintf("Ooops you took more than %d seconds", *timeLimit))
			color.Red("Better Luck Next Time")
			color.Blue("Here is your progress")
			break Loop
		}
	}

	// if time.Now().Sub(startTime).Seconds() >= 30 {
	// 	color.Red("OOOOOOps you took more than 30 seconds to answer questions")
	// 	fmt.Println("Here is a list of your progress")
	// 	break
	// }
	fmt.Println("List of correct answers")
	printAnswers(correctAnswers, answers, questions, rightAnswer)
	fmt.Println("List of wrong answers")
	printAnswers(wrongAnswers, answers, questions, wrongAnswer)
}

func printAnswers(answersMap map[int]string, answers []string, questions []string, answerType string) {
	var c *color.Color
	if len(answersMap) == 0 {
		fmt.Println("None Found")
		return
	}
	for key, value := range answersMap {
		if answerType == wrongAnswer {
			c = color.New(color.FgRed)
			c.Println("Question ", questions[key], " Your answer", value, " Correct Answer ", answers[key])
		} else {
			c = color.New(color.FgGreen)
			c.Println("Question ", questions[key], " Your answer", value, " Correct Answer ", answers[key])
		}
	}
}

func printInstructions() {
	color.Cyan("Here is a list of instructions for the quiz:")
	color.Cyan("The quiz contains a list of questions you will be asked")
	color.Cyan("You will have 30 seconds to answer all the questions")
}
func takeAnswer(reader *bufio.Scanner, ch chan string) {
	if reader.Scan() {
		ch <- reader.Text()
	}
}

// func startTimer(ch chan string) {

// }
