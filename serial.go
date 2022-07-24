package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	serial "github.com/tarm/goserial"
)

const MAX_LINE_SIZE = 4096
const STR_ENDL = "\r\n"
const ENDING_LINE = "END"
const ARDUINO_BAUD = 9600

type ArduinoConnection struct {
	Port io.ReadWriteCloser
}

var connection *ArduinoConnection

func ArduinoConnect(name string) (err error) {
	if connection != nil {
		connection.Port.Close()
	}
	config := serial.Config{
		Name: name,
		Baud: ARDUINO_BAUD,
	}
	port, err := serial.OpenPort(&config)
	if err != nil {
		return fmt.Errorf("could not connect to Arduino: %s", err.Error())
	}
	connection = &ArduinoConnection{Port: port}
	return nil
}

type LineFunc func(line string)

func ReadLoop(fn LineFunc) (err error) {
	if connection == nil {
		return errors.New("not connected")
	}
	bytes := make([]byte, MAX_LINE_SIZE)
	buffer := ""
	for {
		n, err := connection.Port.Read(bytes)
		if n > 0 {
			// fmt.Printf("bytes: %x, %q\n", bytes[:n], bytes[:n])
			buffer += string(bytes[:n])
			// fmt.Printf("buffer: %x, %q\n", buffer, buffer)
			if pos := strings.Index(buffer, STR_ENDL); pos != -1 {
				line := buffer[:pos]
				if line == ENDING_LINE {
					return nil
				}
				fn(line)
				buffer = buffer[pos+len(STR_ENDL):]
			}
		}
		/*
			Esto no parece funcionar, o sea, no sé si puedo detectar que el puerto
			ha sido cerrado desde el Arduino (con Serial.end())

			if n == 0 && err == io.EOF {
				// finished
				break
			}
		*/
		if err != nil {
			return err
		}
	}
}

func SaveArduinoDataToFile(filename string) (err error) {
	if connection == nil {
		return errors.New("not connected")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	aerr := ReadLoop(func(line string) {
		fmt.Printf("%q\n", line)
		fmt.Fprintf(file, "%s\n", line)
	})
	if aerr != nil {
		log.Fatal(err.Error())
	}

	return file.Close()
}
