param(
    [switch]$SkipTests
)

# Generate Swagger docs
swag init -g server.go

if (-not $SkipTests) {
    # Start MongoDB container
    Write-Host "Starting MongoDB container..." -ForegroundColor Yellow
    docker run --name test-mongo -d -p 27017:27017 mongo:latest
    Start-Sleep -Seconds 5  # Give container time to start

    # Run all tests
    Write-Host "Running tests..." -ForegroundColor Green
    go test ./api/v1/test/... -v
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Tests failed! Build aborted." -ForegroundColor Red
        # Cleanup
        docker stop test-mongo
        docker rm test-mongo
        exit 1
    }

    # Cleanup after tests
    docker stop test-mongo
    docker rm test-mongo
}
else {
    Write-Host "Skipping tests..." -ForegroundColor Yellow
}

# Build and run
Write-Host "Building..." -ForegroundColor Green
go build server.go
./server.exe