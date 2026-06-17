# bash file to run the go tests and get coverage report
# run chmod +x run_test_reports.sh
# ./run_test_reports.sh

go test ./... -coverprofile=coverage.out

go tool cover -html=coverage.out