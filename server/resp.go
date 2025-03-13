package server

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func ParseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Ensure it's a RESP array
	if line[0] != '*' {
		return nil, errors.New("invalid RESP array format")
	}

	bulkArg, err := RemoveRespPrefixAndSuffix(line)
	if err != nil {
		return nil, errors.New("invalid RESP array length format")
	}

	numArgs, err := strconv.Atoi(bulkArg)
	if err != nil {
		return nil, errors.New("invalid RESP array length value")
	}

	parsedArgs := make([]string, numArgs)

	for i := 0; i < numArgs; i++ {
		bufferLengthLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.New("invalid RESP bulk string length format")
		}

		if bufferLengthLine[0] != '$' {
			return nil, errors.New("invalid RESP bulk string identifier")
		}

		bufferLengthLineTrimmed, err := RemoveRespPrefixAndSuffix(bufferLengthLine)
		if err != nil {
			return nil, errors.New("invalid RESP bulk string length")
		}

		bufferLength, err := strconv.Atoi(bufferLengthLineTrimmed)
		if err != nil {
			return nil, errors.New("invalid RESP bulk string length value")
		}

		// Handle NULL Bulk Strings ($-1\r\n)
		if bufferLength < 0 {
			parsedArgs[i] = ""
			_, _ = io.ReadFull(reader, make([]byte, 2)) // Read and discard \r\n
			continue
		}

		// Read the actual string (+2 for \r\n)
		buf := make([]byte, bufferLength+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}

		// Validate RESP format (\r\n at the end)
		if buf[bufferLength] != '\r' || buf[bufferLength+1] != '\n' {
			return nil, errors.New("invalid RESP bulk string data")
		}

		parsedArgs[i] = string(buf[:bufferLength])
	}

	return parsedArgs, nil
}

func RemoveRespPrefixAndSuffix(respString string) (string, error) {
	trimmed := strings.TrimSuffix(respString, "\r\n")

	if len(trimmed) == len(respString) {
		return "", errors.New("invalid RESP format")
	}

	return trimmed[1:], nil // Remove RESP prefix (* or $)
}

/*
	Understanding RESP

	Redis uses RESP (Redis Serialization Protocol) to communicate with clients.
	It encodes different types of data as:

	Simple Strings → +OK\r\n
	Errors → -ERR message\r\n
	Integers → :1000\r\n
	Bulk Strings → $5\r\nhello\r\n
	Arrays → *2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n
*/
