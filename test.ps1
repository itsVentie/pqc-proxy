Write-Host "Running comprehensive checks..." -ForegroundColor Cyan

Write-Host "--- Running go vet ---" -ForegroundColor Yellow
go vet ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "Code style/vet issues found!" -ForegroundColor Red
    exit 1
}

Write-Host "--- Running tests with race detector ---" -ForegroundColor Yellow
go test -v -race ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed!" -ForegroundColor Red
    exit 1
}

Write-Host "--- Building project ---" -ForegroundColor Green
go build -v -o pqc-proxy.exe ./cmd/pqc-proxy/main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "Everything is OK! Build successful." -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}