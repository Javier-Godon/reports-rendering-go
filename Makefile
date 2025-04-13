APP_PATH := app/main.go

# Define the protoc command (ensure protoc and plugins are installed)
PROTOC := protoc

# Add the grpc and gin dependencies to go.mod
deps:	
	@echo "Adding dependencies..."
	cd app && \
		go get -u google.golang.org/grpc && \
		go get -u github.com/gin-gonic/gin && \
		go get -u github.com/xuri/excelize/v2 && \
		go mod tidy
	@echo "Dependencies added."

# Default target
all: deps generate run

# Target to generate Go code from proto files
generate:
	@echo "Generating Go code from proto files..."
	mkdir -p app/proto/get_cpu_system_usage # Create the output directory if it doesn't exist inside app
	protoc --proto_path=proto \
		--go_out=app/proto/get_cpu_system_usage \
		--go-grpc_out=app/proto/get_cpu_system_usage \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		proto/get_cpu_system_usage.proto
	mkdir -p app/proto/get_cpu_user_usage     # Create the output directory if it doesn't exist  inside app
	protoc --proto_path=proto \
		--go_out=app/proto/get_cpu_user_usage \
		--go-grpc_out=app/proto/get_cpu_user_usage \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		proto/get_cpu_user_usage.proto
	@echo "Code generation complete."

# Target to run the Go application
run:
	@echo "Running the Go application..."
	cd app && go run main.go

# Target to clean generated files
clean:
	@echo "Cleaning generated files..."
	find . -name "*.pb.go" -type f -delete
	rm -rf app/proto/get_cpu_system_usage app/proto/get_cpu_user_usage #remove the directories
	@echo "Clean complete."

.PHONY: all generate run clean deps
