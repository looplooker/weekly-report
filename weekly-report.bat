@echo off
cd /d %~dp0

:: 显示当前目录和可执行文件是否存在
echo Current directory: %~dp0
if exist "%~dp0weekly-report.exe" (
    echo Executable found
) else (
    echo Executable NOT found
)

:: 创建日志文件夹
if not exist logs mkdir logs

:: 设置日志文件名（使用当前日期时间）
set logfile=logs\weekly_report_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.log

:: 记录开始时间
echo [%date% %time%] Starting weekly report... >> %logfile%

:: 使用完整路径运行程序并记录输出
echo Attempting to run: "%~dp0weekly-report.exe" >> %logfile%
"%~dp0weekly-report.exe" -paths "C:/dev/FB/FB-ERP C:/dev/FB/new-erp-api" -command "log dev --since='1 week ago' --author='李志强' --no-merges --pretty=format:'%%s'" >> %logfile% 2>&1

:: 记录结束时间和状态
if %errorlevel% equ 0 (
    echo [%date% %time%] Weekly report completed successfully. >> %logfile%
) else (
    echo [%date% %time%] Weekly report failed with error code %errorlevel%. >> %logfile%
)

:: 保留最近30天的日志
forfiles /p "logs" /m *.log /d -30 /c "cmd /c del @path" 2>nul
