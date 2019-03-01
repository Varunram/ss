package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"golang.org/x/crypto/sha3"
	"encoding/hex"

	"github.com/btcsuite/btcutil/bech32"
	"golang.org/x/crypto/ssh/terminal"
)

// ScanForInt scans for an integer
func ScanForInt() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, errors.New("Couldn't read user input")
	}
	num := scanner.Text()
	numI, err := strconv.Atoi(num)
	if err != nil {
		return -1, errors.New("Input not a number")
	}
	return numI, nil
}

func ScanRawPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}

func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	return hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
}

func Encrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "Error while opening new GCM block")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return data, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

var questions = []string{
	"What is your mother's first name?",
	"What is your father's first name?",
	"What is your dog's name?",
	"In which city did you meet your spouse in?",
	"Name the one food item that you like the most",
	"Which town was your grandma brought up in?",
}

// Decrypt decrypts a given data stream with a given passphrase
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 || len(passphrase) == 0 {
		return data, errors.New("Length of data is zero, can't decrpyt!")
	}
	key := []byte(SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, errors.Wrap(err, "Error while initializing new cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "failed to initialize new gcm block")
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return plaintext, errors.Wrap(err, "Error while opening gcm mode")
	}
	return plaintext, nil
}

func ReadQuestionsAndGetPassword() (string, error) {
	var answers [3]string
	// now print these questions and the user should choose three questions from these
	log.Println("Choose any three questions that you would like to answer by entering your choice below:")
	log.Println("Choose your first question: ")
	for i, elem := range questions {
		log.Println(i+1, ": ", elem)
	}
	input, err := ScanForInt()
	if err != nil {
		return "", err
	}
	if input > len(questions) {
		return "", fmt.Errorf("Invalid input")
	}
	log.Println(questions[input-1])
	log.Println("Enter your answer (warning: this prompt will not ask you to confirm this answer):")
	answers[0], err = ScanRawPassword()
	if err != nil {
		return "", err
	}

	log.Println("Choose your second question: ")
	for i, elem := range questions {
		log.Println(i+1, ": ", elem)
	}
	input, err = ScanForInt()
	if err != nil {
		return "", err
	}
	if input > len(questions) {
		return "", fmt.Errorf("Invalid input")
	}
	log.Println(questions[input-1])
	log.Println("Enter your answer (warning: this prompt will not ask you to confirm this answer):")
	answers[1], err = ScanRawPassword()
	if err != nil {
		return "", err
	}

	log.Println("Choose your third question: ")
	for i, elem := range questions {
		log.Println(i+1, ": ", elem)
	}
	input, err = ScanForInt()
	if err != nil {
		return "", err
	}
	if input > len(questions) {
		return "", fmt.Errorf("Invalid input")
	}
	log.Println(questions[input-1])
	log.Println("Enter your answer (warning: this prompt will not ask you to confirm this answer):")
	answers[2], err = ScanRawPassword()
	if err != nil {
		return "", err
	}

	password := answers[0] + answers[1] + answers[2]
	return password, nil
}
func EncryptQuestions(mnemonic string) (string, error) {
	// now encrypt the mnemonic with this particular password
	password, err := ReadQuestionsAndGetPassword()
	if err != nil {
		return "", err
	}

	byteData, err := Encrypt([]byte(mnemonic), password)
	if err != nil {
		return "", err
	}

	// convert this byte string into bech32
	conv, err := bech32.ConvertBits(byteData, 8, 5, true)
	if err != nil {
		return "", err
	}
	encoded, err := bech32.Encode("bithyve", conv)
	if err != nil {
		return "", err
	}

	log.Println("ENCODED bech32 string: ", encoded)
	return encoded, nil
}

func DecryptQuestions(encoded string) (string, error) {
	password, err := ReadQuestionsAndGetPassword()
	if err != nil {
		return "", err
	}

	hrp, decoded, err := bech32.Decode(encoded)
	if err != nil {
		log.Println("DECODED: ", decoded)
		return "", err
	}
	if hrp != "bithyve" {
		return "", fmt.Errorf("HRP doesn't match, quitting!")
	}

	// now we need to decrypt the decoded bytedata
	conv, err := bech32.ConvertBits(decoded, 5, 8, true)
	if err != nil {
		return "", err
	}
	// slice off the last bit since its zero
	byteString, err := Decrypt(conv[0:len(conv)-1], password)
	if err != nil {
		return "", err
	}

	return string(byteString), nil
}
