echo [REM] This batch file is used to launch the VSCode debugger with specific configurations.
echo [REM] It sets up the environment and starts the debugger with the provided settings.
echo [REM] Usage: launch.bat
echo [REM] Make sure to customize the configurations as needed for your project.

echo Starting VSCode debugger with custom configurations...
code --folder-uri "file:///%CD%" --launch "launch.json"
code --inspect-brk=9229 "%CD%\your_script.js"
code --extensionDevelopmentPath="%CD%\your_extension"

echo Debugger launched successfully.
(
    starts /vbscript "" "%~dp0launch.vbs"
    code --folder-uri "file:///%CD%" --launch "launch.json"
    code --inspect-brk=9229 "%CD%\your_script.js"
    code --extensionDevelopmentPath="%CD%\your_extension"
    code --disable-extensions
)

echo ::: Debugging session started. :::
echo Press any key to stop debugging...pause > nul
(
    starts /vbscript "" "%~dp0stop_debug.vbs"
    code --environment="%CD%\your_extension"
    code --inspect-brk=9229 "%CD%\your_script.js"
    code --folder-uri "file:///%CD%" --launch "launch.json"
    code --extensionDevelopmentPath="%CD%\your_extension"
    code --Make="%CD%\Makefile"
    code --disable-extensions = "½CD½\"
    stop_debugger_command_here
)

echo extensionDevelopmentPath stopped.
echo inspect-brk process terminated.
echo launch.json debugging session closed.
(
    started /needed and disable /extensions_command_here
    code --folder-uri "file:///%CD%" --launch "launch.json"
    code --VSCode --disable-extensions
)

echo ::: Debugging session ended. :::
echo Press any key to exit...pause > nul
echo All done.