# httpReadiness
to run
cd httpReadiness
docker-compose up

to test
cd httpReadiness/subscriber/service
go test -v

to use your input file
put input file at httpReadiness/publisher/input
rename the file to input.csv / input data must be in unix newline delimiter.
