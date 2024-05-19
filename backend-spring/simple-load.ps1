# Define the number of requests
$numRequests = 50

# Define the API endpoint
$apiUrl = "http://localhost:8080/api/v1/markers"

# Create an array to store the jobs
$jobs = @()

# Send requests in parallel
for ($i = 0; $i -lt $numRequests; $i++) {
    $jobs += Start-Job -ScriptBlock {
        param ($url)
        try {
            $response = Invoke-RestMethod -Uri $url -Method Get
            Write-Output "Request $($url): $($response.StatusCode)"
        } catch {
            Write-Output "Request $($url) failed: $($_.Exception.Message)"
        }
    } -ArgumentList $apiUrl
}

# Wait for all jobs to complete
$jobs | ForEach-Object {
    Wait-Job $_
    Receive-Job $_
    Remove-Job $_
}
