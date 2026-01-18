package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Input struct {
	Model struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
}

type Output struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func mutedColor(text string) string {
	return fmt.Sprintf("\033[90m%s\033[0m", text)
}

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Println("{}")
		return
	}

	var value string
	if input.Model.DisplayName != "" {
		value = input.Model.DisplayName
	} else if input.Model.ID != "" {
		value = input.Model.ID
	} else {
		value = mutedColor("N/A")
	}

	output := Output{
		Label: "Model",
		Value: value,
	}

	json.NewEncoder(os.Stdout).Encode(output)
}
