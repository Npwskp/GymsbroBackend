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
    $testOutput = go test ./api/v1/test/... -v -json | ForEach-Object { $_ | ConvertFrom-Json }
    
    # Initialize counters and tracking
    $testNames = @{}
    $passed = 0
    $failed = 0

    # Count test results
    $testOutput | ForEach-Object {
        if ($_.Action -eq "run") {
            $testNames[$_.Test] = $true
        }
        if ($_.Action -eq "pass") { $passed++ }
        if ($_.Action -eq "fail") { $failed++ }
    }

    # Get total unique tests
    $total = $testNames.Count

    # Display test summary
    Write-Host "`nTest Summary:" -ForegroundColor Cyan
    Write-Host "Total Tests: $total" -ForegroundColor White
    Write-Host "Passed: $passed" -ForegroundColor Green
    Write-Host "Failed: $failed" -ForegroundColor Red

    if ($failed -gt 0) {
        Write-Host "`nTests failed! Build aborted." -ForegroundColor Red
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