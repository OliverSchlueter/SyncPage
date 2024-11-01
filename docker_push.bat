@echo off
setlocal

REM Define variables
set IMAGE_NAME=sync-page
set IMAGE_TAG=latest
set REPO_URL=oliverschlueter

REM Build the Docker image
echo Building Docker image...
docker build -t %IMAGE_NAME%:%IMAGE_TAG% .

REM Tag the image for the repository
echo Tagging Docker image for push...
docker tag %IMAGE_NAME%:%IMAGE_TAG% %REPO_URL%/%IMAGE_NAME%:%IMAGE_TAG%

REM Push the image to the repository
echo Pushing Docker image to repository...
docker push %REPO_URL%/%IMAGE_NAME%:%IMAGE_TAG%

REM Cleanup
endlocal
echo Done!