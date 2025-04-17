package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	envJSON := getEnvAsJSON()

	// Output the JSON string
	//nolint:forbidigo
	fmt.Println(envJSON)
}

// getEnvAsJSON retrieves environment variables as a JSON string.
func getEnvAsJSON() string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		pair := splitEnv(env)
		envs[pair[0]] = pair[1]
	}

	//nolint:errchkjson
	envJSON, _ := json.Marshal(envs)

	return string(envJSON)
}

// splitEnv splits an environment variable into key and value parts.
func splitEnv(env string) []string {
	for i := range env {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}

	return []string{env, ""}
}
