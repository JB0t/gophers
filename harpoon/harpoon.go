package main

import (
	"strings"
	"bufio"
	"strconv"
	"slices"
	"path/filepath"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Local utils
func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		os.Exit(1)
	}()
}

func getEnvVar(key string, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func checkInputIsIndex(possibleIndexStr string, debug bool) int {
	// If is index, return true
	possibleIndex, err := strconv.Atoi(possibleIndexStr)
	if err != nil {
		fmt.Printf("Could not convert arg %s to int: %s\n",possibleIndexStr,err)
		return 0
	}
	_, exists := harpoons[possibleIndex]
	if exists {
		return possibleIndex
	} else {
		if debug {
			fmt.Printf("Index %d not found. Here is the current list:\n",possibleIndex)
			printHarpoons()
		}
		return 0
	}
}

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

// Handle harpoons
func printHarpoons() {
	if len(harpoons) > 0 {
		for index,path := range harpoons {
			fmt.Printf("[%d] %s\n",index,path)
		}
	}
}
func checkHarpoonExists(path string)bool {
	for _,hPath := range harpoons {
		if hPath == path {
			return true
		}
	}
	return false
}

func loadHarpoons(debug bool)bool {
	if checkExists(h_file_path) {
		h_file, err := os.Open(h_file_path)
		if err != nil {
			fmt.Printf("Could not open %s: %s\n",h_file_path,err)
			return false
		} else {
			err := json.NewDecoder(h_file).Decode(&harpoons)
			if err != nil {
				fmt.Printf("Could not decode %s: %s\n",h_file_path,err)
				return false
			}
			if debug {
				printHarpoons()
			}
		}
	}
	return true
}
func saveHarpoons(debug bool) {
	removeFileIfExists(h_file_path)
	h_out, err := os.Create(h_file_path)
	if err != nil {
		fmt.Printf("Could not create file %s: %s\n",h_file_path,err)
		return
	}
	w_err := json.NewEncoder(h_out).Encode(harpoons)
	if w_err != nil {
		fmt.Printf("Could not write harpoons to file %s: %s\n",h_file_path,w_err)
		return
	}
}
func addHarpoon(path string, debug bool){
	var newIndex = 1
	if len(harpoons) > 0 {
		if checkHarpoonExists(path){
			fmt.Printf("Path exists at index %d\n",func(p string)int{for idx,pth := range harpoons{if p == pth{return idx}};return 0}(path))
		}
		keys := make([]int, len(harpoons))
		i := 0
		for k := range harpoons {
			keys[i] = k
			i++
		}
		newIndex = slices.Max(keys) + 1
	}
	harpoons[newIndex] = path
	saveHarpoons(debug)
}

// get/delete by index
func deletePathByIndex(index int, debug bool) {
	if len(harpoons) > 0 {
		path, exists := harpoons[index]
		if exists {
			delete(harpoons, index)
			saveHarpoons(debug)
			if debug {
				fmt.Printf("Removed [%d] %s\n",index,path)
			}
		} else {
			if debug {
				fmt.Printf("Index %d not found\n",index)
			}
		}
	}
}

func getPathByIndex(index int,debug bool) {
	if harpoons != nil {
		path, exists := harpoons[index]
		if exists {
			removeFileIfExists(h_out_path)
			h_out, err := os.Create(h_out_path)
			if err != nil {
					fmt.Printf("Could not create file %s: %s\n",h_out_path,err)
					return
			}
			_, w_err := h_out.WriteString(path)
			if w_err != nil {
				fmt.Printf("Could not write to file %s: %s\n",h_out_path,err)
				return
			}
			if debug {
				fmt.Printf("%s written to %s\n",path,h_out_path)
			}
		} else {
			if debug {
				fmt.Printf("Index %d not found\n",index)
			}
		}
	}
}

var harpoons = make(map[int]string)
var h_file_path = fmt.Sprintf("%s/.harpoons",getEnvVar("GOPATH",getEnvVar("HOME","tmp")))
var h_out_path = fmt.Sprintf("%s/.h_out",getEnvVar("GOPATH",getEnvVar("HOME","tmp")))

func main(){
	var addPath string
	var goToPath int
	var deletePath int
	var clearPaths bool
	var debug bool
	flag.StringVar(&addPath, "a", "", "Add a path")
	flag.IntVar(&goToPath, "i", 0, "Go to path at index")
	flag.IntVar(&deletePath, "r", 0, "Delete the path at index")
	flag.BoolVar(&clearPaths, "c", false, "Clear all indexes")
	flag.BoolVar(&debug, "d", false, "Set debug mode")

	flag.Parse()
	if loadHarpoons(debug) == false {
		return
	}

	if debug {
		fmt.Printf("addPath: %s\n",addPath)
		fmt.Printf("goToPath: %d\n",goToPath)
		fmt.Printf("deletePath: %d\n",deletePath)
		fmt.Printf("clearPaths: %s\n",func(b bool)string{if b{return"enabled"}else{return"disabled"}}(clearPaths))
	}

	// CLA options
	if clearPaths {
		if addPath == "" && goToPath + deletePath == 0 {
			if debug {
				fmt.Println("Clearing paths.")
			}
			removeFileIfExists(h_file_path)
		} else {
			fmt.Println("Ambiguous arguments. Please select 1")
		}
	} else if deletePath > 0 {
		if addPath == "" && goToPath == 0 && clearPaths == false {
			if debug {
				fmt.Printf("Removing path at index %d.\n",deletePath)
			}
			deletePathByIndex(deletePath, debug)
		} else {
			fmt.Println("Ambiguous arguments. Please select only one option.")
		}
	} else if goToPath > 0 {
		if addPath == "" && deletePath == 0 && clearPaths == false {
			if debug {
				fmt.Printf("Going to path at index %d.\n",goToPath)
			}
			getPathByIndex(goToPath, debug)
		} else {
			fmt.Println("Ambiguous arguments. Please select only one option.")
		}
	} else if addPath != "" {
		if goToPath + deletePath == 0 && clearPaths == false {
			absAddPath, err := filepath.Abs(addPath)
			if err != nil {
				fmt.Printf("Could not resolve filepath %s: %s",addPath,err)
				return
			}
			if debug {
				fmt.Printf("Adding path %s\n",absAddPath)
			}
			addHarpoon(absAddPath, debug)
		} else {
			fmt.Println("Ambiguous arguments. Please select only one option.")
		}
	} else {
		if len(os.Args) > 1 {
			if debug {
				fmt.Printf("Attempting to harpoon index passed by Arg[1]: %s\n",os.Args[1])
			}
			possibleIndexStr := os.Args[1]
			checkIndex := checkInputIsIndex(possibleIndexStr,debug)
			if checkIndex > 0 {
				getPathByIndex(checkIndex,debug)
			}
		} else {
			if len(harpoons) > 0 {
				fmt.Println("Harpoon Menu.")
				printHarpoons()
				// Loop until user selects a valid index or q or ctrl+c
				reader := bufio.NewReader(os.Stdin)
				for {
					fmt.Print("index> ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Printf("Error reading input: %v\n", err)
						break
					}
					input = strings.TrimSpace(input)
					if input == "" {
						continue // Skip empty lines
					} else if input == "q" {
						return
					} else {
						checkIndex := checkInputIsIndex(input,debug)
						if checkIndex > 0 {
							getPathByIndex(checkIndex,debug)
							break
						}
					}
				}
			} else {
				fmt.Println("Harpoon Empty.")
			}
		}
	}
}

