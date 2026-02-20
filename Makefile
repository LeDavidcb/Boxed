# Makefile for managing the project

# Variables
BUILD_DIR=build
CMD_DIR=cmd
ENV_FILE=.env
ENVCLI=$(BUILD_DIR)/cli
MAIN_EXECUTABLE=$(BUILD_DIR)/api

MIGRATIONS_DIR=migrations

# Dependencies
check-deps:
	@command -v node >/dev/null 2>&1 || { echo "Node.js is not installed."; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "npm is not installed."; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "Go compiler is not installed."; exit 1; }
	@command -v ffmpeg >/dev/null 2>&1 || { echo "FFmpeg is not installed."; exit 1; }
	@echo "All dependencies are installed."

# Build command
build: check-deps
	@mkdir -p $(BUILD_DIR)
	@for dir in $(CMD_DIR)/*; do \
	  if [ -d "$$dir" ]; then \
	    cmd_name=`basename $$dir`; \
	    go build -o $(BUILD_DIR)/$$cmd_name ./$$dir || exit 1; \
	    echo "Built: $$cmd_name"; \
	  fi; \
	done

# Run project
run: build
	@printf "Run migrations? [Y/n] "; \
	read answer; \
	case "$$answer" in \
	  [yY]*|"") $(MAKE) --no-print-directory migrate-up ;; \
	  [nN]*) echo "Skipping migrations." ;; \
	  *) echo "Invalid input. Aborting."; exit 1 ;; \
	esac
	@if [ ! -f $(ENV_FILE) ]; then \
	  echo "$(ENV_FILE) not found. Running envcli to generate it..."; \
	  $(ENVCLI) || { echo "Failed to run envcli."; exit 1; }; \
	fi;
	@echo "Running the project..."
	@./$(MAIN_EXECUTABLE)

# Help
default: help

migrate-%:
	@if [ ! -f $(ENV_FILE) ]; then echo "$(ENV_FILE) not found."; exit 1; fi
	@DB_URL=$$(grep '^DB_URL=' $(ENV_FILE) | cut -d'=' -f2-) && goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" $(subst migrate-,,$@)

help:
	@echo "Makefile Usage:"
	@echo "  make check-deps   Check if required dependencies are installed."
	@echo "  make build        Build all cmd main.go files to ./build."
	@echo "  make run          Build and run the project, auto-generate .env if missing."
	@echo "  make migrate-up   Run database migrations up."
	@echo "  make migrate-down Roll back database migrations."
