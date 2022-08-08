@echo off
echo Convert CUBE to SPHERE droplet
echo.

SET /A COUNT=0
FOR %%V in (%*) DO SET /A COUNT=COUNT+1

IF NOT %COUNT% == 6 GOTO ERROR
IF "%~1" == "" GOTO ERROR
IF "%~2" == "" GOTO ERROR
IF "%~3" == "" GOTO ERROR
IF "%~4" == "" GOTO ERROR
IF "%~5" == "" GOTO ERROR
IF "%~6" == "" GOTO ERROR
IF NOT EXIST "%~1" GOTO ERROR
IF NOT EXIST "%~2" GOTO ERROR
IF NOT EXIST "%~3" GOTO ERROR
IF NOT EXIST "%~4" GOTO ERROR
IF NOT EXIST "%~5" GOTO ERROR
IF NOT EXIST "%~6" GOTO ERROR

set KRPANOTOOLSEXE=krpanotools64.exe
if "%PROCESSOR_ARCHITECTURE%" == "x86" set KRPANOTOOLSEXE=krpanotools32.exe

"%~dp0\%KRPANOTOOLSEXE%" cube2sphere -config=convertdroplets.config "%~1" "%~2" "%~3" "%~4" "%~5" "%~6"
GOTO DONE

:ERROR
echo.
echo Usage:
echo.
echo - Drag and drop cubical panorama images on this droplet
echo   to convert them to spherical panorama images.
echo.
echo - The cubical images must have suffixes like "_l", "_f"
echo   or "left", "front" in their filename to allow an
echo   automatic identifcation of the cubeside.
echo.
echo.
echo Settings:
echo.
echo - See the "convertdroplets.config" file for settings
echo   regarding output imageformat or compression.
echo.
echo.

:DONE
echo.
