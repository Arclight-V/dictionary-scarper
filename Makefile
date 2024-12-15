# Name of the executable file
APP_NAME=translator

# Output directory for the compiled binary
BIN_DIR=bin

# Main package (entry point)
MAIN=delivery/main.go

# Default paths to the CSV files (can be overridden when running make)
FRENCH_CSV?=data/french.csv
ENGLISH_CSV?=data/english.csv
RUSSIAN_CSV?=data/russian.csv

# Build target: compiles the main package into a binary
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN)

# Run target: builds the application and runs it with the specified CSV file paths
.PHONY: run
run: build
	$(BIN_DIR)/$(APP_NAME) -french=$(FRENCH_CSV) -english=$(ENGLISH_CSV) -russian=$(RUSSIAN_CSV)

# Clean target: removes the compiled binary and cleans up the bin directory
.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
