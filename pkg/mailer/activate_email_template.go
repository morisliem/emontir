package mailer

import "fmt"

func ActivationEmailLinkTemplate(link string) string {
	greet := "Hi There\n\n"
	instruction := "Please activate your email by clicking the link below\n"
	regards := "\n\nCheers\ne-Montir team"
	return fmt.Sprintf(greet+instruction+"%s"+regards, link)
}
