$totalLines = 0
$lineCounts = git ls-files | Where-Object { $_ -match '\.java$' } | ForEach-Object {
    $filePath = $_
    if (Test-Path $filePath) {
        $lineCount = (Get-Content $filePath).Count
        $totalLines += $lineCount
        [PSCustomObject]@{Lines=$lineCount; File=$filePath}
    } else {
        Write-Warning "Could not find file at path $filePath"
    }
}

$lineCounts | Format-Table -AutoSize
Write-Host "Total Lines: $totalLines"