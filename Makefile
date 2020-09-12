.DEFAULT_GOAL := default

default:
	@echo "Fibonacci generator Setup, Test, Run Instructions"
	@echo "Only Tested on MacOS Catalina with Docker Desktop/Golang 1.15 installed."	
	@echo "Docker Desktop Installation REQUIRED"		
	@echo "Docker Compose Installation REQUIRED"			
	@echo "Golang 1.15 Installation REQUIRED"				
	@echo "Windows behavior unknown"		
	@echo "#1 Enter 'make desktop' to start your docker engine"		
	@echo "#2 Enter 'make setup' to setup your machine with docker running"	
	@echo "#3 Enter 'make test' to Run tests"
	@echo "#4 Enter 'make bench' to benchmark the fibonacci generator"	
	@echo "#5 Enter 'make run' to run the test demo"
	@echo "Extra-> Enter 'make cleanup' to wipe your docker containers"	
desktop:
	@echo "Starting docker desktop on MacOS"
	open /Applications/Docker.app &
	@echo "Wait around 10 seconds for docker to finish its initialization"	
cleanup:	
	@echo "Cleaning your docker environment"
	docker ps -a -q | xargs docker stop
	docker ps -a -q | xargs docker rm	
	docker rmi -f pg-fibo | true		
	docker rmi -f fibo | true			
setup:
	@echo "Setting up your machine"
	docker ps -a -q | xargs docker stop
	docker ps -a -q | xargs docker rm	
	docker build -t pg -f ./Dockerfile.postgres .	
	docker run -d --name pg-fibo -e POSTGRES_PASSWORD=mysecretpassword -t -p 5432:5432 pg
test:
	@echo "Running tests"
	go test 	
bench:
	@echo "Running Benchmarks"
	go test -run=Bench -bench=.
run:
	@echo "Running test demo.."
	#Stop all the containers
	docker ps -a -q | xargs docker stop
	#Remove all the containers	
	docker ps -a -q | xargs docker rm	
	#Remove any pg images
	docker rmi -f pg-fibo | true		
	#Build the container again
	docker build -t pg-fibo -f ./Dockerfile.postgres .	
	rm -f ./go-fibonacci
	rm -f ./go-modules	
	#Build our golang binary
	go clean

	GOARCH=amd64 GOOS=linux go build -o ./go-fibonacci --ldflags '-s -w -extldflags "-static"'		

	docker build -t fibo -f ./Dockerfile.fibo .		
	@echo "###################################################"		
	@echo "##############IMPORTANT README#####################################"	
	@echo "Now enter 'docker-compose up -d <enter>' to execute the demo"
	@echo "Enter 'docker-compose down<enter>' to stop the demo"	
	@echo "Open your web browser and navigate to: 'http://localhost:8081' to see the API in action"		
	@echo "###################################################"		
	


