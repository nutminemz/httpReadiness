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

architect
1xmaster(pub) -->Redis--> 2xslave(sub)
master read file, assign to slaves, collect report
slave fetch net/http checking for readiness then update report
master check job status every sec when finish submit report to LINE

scalable via compose file