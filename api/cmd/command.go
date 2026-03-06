package main

import (
	"fmt"
	"os"
	"strings"
)

func HandleCommand(args []string) {
	if len(args) == 0 {
		showHelp()
		return
	}

	cmd := args[0]
	params := args[1:]

	switch cmd {
	case "hello":
		helloCommand(params)
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("For available commands, use the 'help' command.")
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Flowize API - Komut Satırı Aracı")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  make cmd <command> [arguments]")
	fmt.Println("  go run cmd/*.go <command> [arguments]")
	fmt.Println()
	fmt.Println("Mevcut Komutlar:")
	fmt.Println("  hello [arguments]       - Hello command run")
	fmt.Println("  help                    - Show this help message")
	fmt.Println()
	fmt.Println("Alternatif:")
	fmt.Println("  go run cmd/*.go vapi-call-tenant <assistant_id> <vapi_phone_number_id> <to_phone> [to_name]")
	fmt.Println("  go run cmd/*.go vapi-call-lead <landing_page_id> <to_phone> [to_name]")
}

func helloCommand(params []string) {

	if len(params) > 0 {
		fmt.Println("Hello", strings.Join(params, " "))
		return
	}

}
