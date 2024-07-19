# Compile the Go application
go build main.go

# Check if the build was successful
if ($?) {
    # Make a beep sound
    [console]::beep(1000, 500)
    
    # Run the executable
    ./main.exe
}
