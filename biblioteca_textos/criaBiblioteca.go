// Cria uma biblioteca com muitos pastas com muitos arquivos.
// Cada arquivo contém um número de palavras.
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const diferentWords int = 692000
const numberFolders int = 10
const wordsPerFile int = 68

// Cria uma biblioteca com muitos pastas com muitos arquivos.
// Cada arquivo contém um número de palavras.
func main() {

	startingTime := time.Now()
	numberOfFiles := diferentWords * numberFolders / wordsPerFile
	fmt.Println("Número de Arquivos: ", numberOfFiles)
	filesPerFolder := numberOfFiles / numberFolders
	fmt.Println("Números de Arquivos por Pasta: ", filesPerFolder)

	var wg sync.WaitGroup

	for i := 0; i < numberFolders; i++ {
		wg.Add(1)

		wordList := make([]string, diferentWords)
		file, err := os.Open("." + string(filepath.Separator) + "palavras_fonte.txt")

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// go = paralelo -> 2 minutos
		// sem go = série -> 4 minutos
		go func(i int) {

			wordCount := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				wordList[wordCount] = scanner.Text()
				wordCount++
			}

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(wordList), func(i, j int) { wordList[i], wordList[j] = wordList[j], wordList[i] })

			folderName := "pasta" + strconv.Itoa(i)
			dirFolder := "." + string(filepath.Separator) + folderName
			os.Mkdir(dirFolder, 0777)
			for j := 0; j < filesPerFolder; j++ {

				dir := dirFolder + string(filepath.Separator) + "livro-" + strconv.Itoa(j) + ".txt"
				//fmt.Println(dir)
				newFile, err2 := os.Create(dir)
				if err2 != nil {
					panic(err2)
				}

				w := bufio.NewWriter(newFile)
				for i := 0; i < wordsPerFile; i++ {
					newWord := wordList[len(wordList)-1]
					wordList = wordList[:len(wordList)-1]
					if newWord != "" {
						w.WriteString(newWord + "\n")
					}
				}
				w.Flush()
				newFile.Close()
				if j%(filesPerFolder/100) == 0 {
					percentage := 100*j/filesPerFolder + 1
					fmt.Println("Progresso: ", dirFolder, "/", numberFolders, ": ", percentage, "%")
				}
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println("Completado em ", time.Since(startingTime))
}
