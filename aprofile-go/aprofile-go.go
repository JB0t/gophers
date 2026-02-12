package main

import (
	"regexp"
	"github.com/eiannone/keyboard"
	"os/exec"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Profile struct {
	name string
	sso_start_url string
	sso_region string
	sso_account_id string
	sso_role_name string
	region string
}

var profiles = make(map[string]Profile)
var a_out_path = fmt.Sprintf("%s/.a_out",getEnvVar("GOPATH",getEnvVar("HOME","tmp")))
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var debug = false
var test = false

// Unit tests
func testUserInput() {
	testStringList := []string{"abcdefg","1234567","!@#$%^&*()", "abc123!@#", "1234567890", "!!@#$%^&*()_+"}
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Failed to open keyboard: %s\n", err)
		return
	}
	defer func() {
		_ = keyboard.Close()
	}()
	userSelect := make([]int32, 0)
	userSelectStr := ""
	totalStringList := make([]string, 0)
	currentStringList := make([]string, 0)
	// init profile list
	for _,string := range testStringList {
			totalStringList = append(totalStringList, string)
	}
	currentStringList = totalStringList
	for{
		fmt.Printf("\nlen(userSelect): %d\nlen(userSelectStr): %d\n",len(userSelect), len(userSelectStr))
		fmt.Println("Test strings:")
		for _,string := range currentStringList {
			fmt.Printf("%s\n", string)
		}
		fmt.Println("Input under here:")
		fmt.Print(userSelectStr)
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Printf("Failed to read keyboard input: %s\n", err)
			return
		}
		if key == keyboard.KeyEnter {
			fmt.Printf("Final input: '%s'\n", userSelectStr)
			return
		} else if key == keyboard.KeyTab {
			if len(currentStringList) == 0 {
				fmt.Println("No users match current input, ignoring tab")
			} else if len(currentStringList) == 1 {
				clearScreen()
				fmt.Printf("Final input: '%s'\n", currentStringList[0])
				return
			} else {
				largestListSharedPrefix := getNextSharedPrefix(currentStringList)
				fmt.Printf("Next largest prefix chunk: '%s'\n", largestListSharedPrefix)
				// set the user selection to the total up to the next chunk
				if len(largestListSharedPrefix) > len(userSelectStr) {
					userSelectStr = largestListSharedPrefix
					userSelect = []rune(userSelectStr)
					currentStringList = selectPrefixMatches(userSelectStr, currentStringList)
				}
			}
		} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if len(userSelect) > 0 {
				userSelect = userSelect[:len(userSelect)-1]
				userSelectStr = convertToString(userSelect)
				// todo, the previx matches can be a linked list, even a prefix tree
				currentStringList = selectPrefixMatches(userSelectStr, totalStringList)
			}
		} else {
			userSelect = append(userSelect, char)
			userSelectStr = convertToString(userSelect)
			currentStringList = selectPrefixMatches(userSelectStr, currentStringList)
		}
		fmt.Printf("Current input: '%s'\n", userSelectStr)
		if len(userSelectStr) >= 10 {
			fmt.Println("Test input limit reached, exiting test\n")
			fmt.Printf("Final input: '%s'\n", userSelectStr)
			return
		}
	}
}

// Utils

func checkExists(path string) bool {
	_, error := os.Stat(path)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

func removeFileIfExists(path string) {
	if checkExists(path) {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Could not remove file %s: %s\n",path,err)
			return
		}
	}
}

func writeProfile(profile string) {
	removeFileIfExists(a_out_path)
	a_out, err := os.Create(a_out_path)
	if err != nil {
		fmt.Printf("Could not create file %s: %s\n",a_out_path,err)
		return
	}
	_, w_err := a_out.WriteString(profile)
	if w_err != nil {
		fmt.Printf("Could not write to file %s: %s\n",a_out_path,err)
		return
	}
	if debug {
		fmt.Printf("%s written to %s\n",profile,a_out_path)
	}
}

func getEnvVar(key string, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func maximumProfileNameLength(currentProfileList []string) int {
	currentMax := 0
	for _,profileName := range currentProfileList {
		if len(profileName) > currentMax {
			currentMax = len(profileName)
		}
	}
	return currentMax
}

func maximumProfileName(currentProfileList []string) string {
	currentMax := currentProfileList[0]
	for _,profileName := range currentProfileList {
		if len(profileName) > len(currentMax) {
			currentMax = profileName
		}
	}
	return currentMax
}

func minimumProfileName(currentProfileList []string) string {
	currentMin := currentProfileList[0]
	for _,profileName := range currentProfileList {
		if len(profileName) < len(currentMin) {
			currentMin = profileName
		}
	}
	return currentMin
}

func convertToString(runeList []int32) string {
	if len(runeList) == 0 {
		return ""
	}
	count := 0
	for _, r := range runeList {
		if r != 0 {
			count++
		}
	}
	return string(runeList[:count])
}

func getNextStringChunk(fullstring string, compareString string) string {
	cleanString := nonAlphanumericRegex.ReplaceAllString(fullstring, " ")
	alphanumericDelims := nonAlphanumericRegex.FindAllString(fullstring, -1)
	chunks := strings.Fields(nonAlphanumericRegex.ReplaceAllString(fullstring, " "))
	if debug {
		fmt.Printf("DEBUG: Clean string: '%s'\n", cleanString)
		fmt.Printf("DEBUG: alphanumericDelims: '%v'\n", alphanumericDelims)
		fmt.Printf("DEBUG: chunks: '%v'\n", chunks)
	}
	for i, chunk := range chunks {
		if strings.Contains(compareString, chunk) == false {
			if debug {
				fmt.Printf("DEBUG: Returning chunk '%s' at index %d\n", chunk, i)
			}
			return chunks[i]
		}
	}
	return ""
}

func getNextSharedPrefix(stringList []string) string {
	// need index of max profile name, maybe we want the smallest actually
	smallestProfileName := minimumProfileName(stringList)
	if debug {
		fmt.Printf("DEBUG: Smallest profile name: '%s'\n", smallestProfileName)
	}
	for i := 0; i < len(stringList)-1; i++ {
		for j := 0; j < len(smallestProfileName); j++ {
			if debug {
				fmt.Printf("DEBUG: Comparing character '%s' at index %d for string '%s' and '%s'\n", stringList[i][j], j, stringList[i], stringList[i+1])
			}
			if stringList[i][j] != stringList[i+1][j] {
				return stringList[0][:j]
			}
		}
	}
	return stringList[0]
}

func main() {
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}
	readProfiles()
	if os.Getenv("TEST") == "true" {
		test = true
	}
	if test {
		testUserInput()
	} else {
		selectedProfileName := selectProfileName()
		if selectedProfileName == "" {
			return
		}
		writeProfile(selectedProfileName)
	}
}

// Profile ops
func selectPrefixMatches(profileSelect string, currentProfileList []string) []string {
	if debug {
		fmt.Printf("DEBUG: Selecting prefix matches for input '%s' from profile list: %v\n", profileSelect, currentProfileList)
	}
	if profileSelect == "" {
		return currentProfileList
	}
	prefixMatches := make([]string, 0)
	for _,profileName := range currentProfileList {
		if strings.HasPrefix(profileName, profileSelect) {
			prefixMatches = append(prefixMatches, profileName)
		}
	}
	if debug {
		fmt.Printf("DEBUG: Found prefix matches: %v\n", prefixMatches)
	}
	return prefixMatches
}

func selectProfileName() string {
	if len(profiles) == 0 {
		fmt.Println("No profiles found")
		return ""
	}
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Failed to open keyboard: %s\n", err)
		return ""
	}
	defer func() {
		_ = keyboard.Close()
	}()
	profileSelect := make([]int32, 0)
	profileSelectStr := ""
	totalProfileList := make([]string, 0)
	currentProfileList := make([]string, 0)
	// init profile list
	for profileName := range profiles {
			totalProfileList = append(totalProfileList, profileName)
		}
	currentProfileList = totalProfileList
	if debug {
		fmt.Printf("DEBUG: Current profile list before input loop: %v\n", currentProfileList)
	}
	for{
		clearScreen()
		fmt.Println("Profiles:")
		for _,profileName := range currentProfileList {
			fmt.Printf("%s\n", profileName)
		}
		fmt.Print(profileSelectStr)
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Printf("Failed to read keyboard input: %s\n", err)
			return ""
		}
		if key == keyboard.KeyEnter {
			fmt.Println()
			break
		} else if key == keyboard.KeyEsc {
			fmt.Println("Input cancelled, exiting")
			return ""
		} else if key == keyboard.KeyTab {
			if len(currentProfileList) == 0 {
				fmt.Println("No profiles match current input, ignoring tab")
			} else if len(currentProfileList) == 1 {
				clearScreen()
				return currentProfileList[0]
			} else {
				largestListSharedPrefix := getNextSharedPrefix(currentProfileList)
				fmt.Printf("Next largest prefix chunk: '%s'\n", largestListSharedPrefix)
				// set the profile selection to the total up to the next chunk
				if len(largestListSharedPrefix) > len(profileSelectStr) {
					profileSelectStr = largestListSharedPrefix
					profileSelect = []rune(profileSelectStr)
					currentProfileList = selectPrefixMatches(profileSelectStr, currentProfileList)
				}
			}
		} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if len(profileSelect) > 0 {
				profileSelect = profileSelect[:len(profileSelect)-1]
				profileSelectStr = convertToString(profileSelect)
				// todo, the previx matches can be a linked list, even a prefix tree
				currentProfileList = selectPrefixMatches(profileSelectStr, totalProfileList)
			}
		} else if len(profileSelect)+1 > maximumProfileNameLength(currentProfileList) {
			fmt.Printf("Input '%s' exceeds maximum profile name length, ignoring input\n", string(char))
		} else {
			profileSelect = append(profileSelect, char)
			profileSelectStr = convertToString(profileSelect)
			currentProfileList = selectPrefixMatches(profileSelectStr, currentProfileList)
		}
	}
	clearScreen()
	return convertToString(profileSelect)
}

func printProfile(profile *Profile) {
	fmt.Printf("Profile Name: %s\n", profile.name)
	fmt.Printf("SSO Start URL: %s\n", profile.sso_start_url)
	fmt.Printf("SSO Region: %s\n", profile.sso_region)
	fmt.Printf("SSO Account ID: %s\n", profile.sso_account_id)
	fmt.Printf("SSO Role Name: %s\n", profile.sso_role_name)
	fmt.Printf("Region: %s\n", profile.region)
}

func readProfiles() {
	// get config file
	config_path := fmt.Sprintf("%s/.aws/config",getEnvVar("HOME","tmp"))
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		fmt.Printf("No config file found at %s\n",config_path)
		return
	}
	config_file, err := os.Open(config_path)
	if err != nil {
		fmt.Printf("Could not open config file at %s: %s\n",config_path,err)
	}
	if debug{
		fmt.Printf("Reading config file at %s\n",config_path)
	}
	defer config_file.Close()

	// parse config file
	line_reader := bufio.NewScanner(config_file)
	currProfile := Profile{}
	for line_reader.Scan() {
		curr := line_reader.Text()
		if strings.Contains(curr, "[profile ") {
			if currProfile.name != "" {
				if debug{
					fmt.Printf("Adding profile %s to profiles map\n", currProfile.name)
				}
				profiles[currProfile.name] = currProfile
			}
			currProfile.name = strings.TrimPrefix(strings.Trim(curr, "[]"), "profile ")
		} else if strings.Contains(curr, "sso_region") {
			currProfile.sso_region = strings.TrimSpace(strings.Split(curr, "=")[1])
		} else if strings.Contains(curr, "sso_account_id") {
		     currProfile.sso_account_id = strings.TrimSpace(strings.Split(curr, "=")[1])
		} else if strings.Contains(curr, "sso_role_name") {
		     currProfile.sso_role_name = strings.TrimSpace(strings.Split(curr, "=")[1])
		} else if strings.Contains(curr, "sso_start_url") {
			currProfile.sso_start_url = strings.TrimSpace(strings.Split(curr, "=")[1])
		} else if strings.Contains(curr, "region") {
			currProfile.region = strings.TrimSpace(strings.Split(curr, "=")[1])
		}
	}
	fmt.Printf("Found %d profiles:\n", len(profiles))
	if debug {
		for _, profile := range profiles {
			printProfile(&profile)
		}
	}
}
