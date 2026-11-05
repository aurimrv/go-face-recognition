# Set up STEP variavel 

```bash
STEP=01
```

# Measuring code coverage and generating textual and html reports
```bash
go test -coverprofile=${STEP}-coverage.out ./...
go tool cover -func=${STEP}-coverage.out &> ${STEP}-coverage.txt
go tool cover -html=${STEP}-coverage.out -o ${STEP}-coverage.html
```

# Measuring mutation score in a specific package and creating the html report
```bash
gomutants --html-output ${STEP}-mutation-report-entity.html ./internal/entity
```

# Measuring mutation score in all Go code and creating the html report
```bash
gomutants --html-output ${STEP}-mutation-report-full.html ./internal/entity ./internal/usecases
```