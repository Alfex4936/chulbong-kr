# This is the equivalent of '@echo off' in batch
$ErrorActionPreference = "SilentlyContinue"

./gradlew clean test jacocoTestReport