@echo off

@REM wails build -s -platform darwin/amd64 -o darwin-amd64 >NUL 2>NUL
@REM wails build -s -platform darwin/arm64 -o darwin-arm64 >NUL 2>NUL
@REM wails build -s -platform darwin/universal -o darwin-universal >NUL 2>NUL

@REM wails build -s -platform linux/amd64 -o linux-amd64 >NUL 2>NUL
@REM wails build -s -platform linux/arm64 -o linux-arm64 >NUL 2>NUL
@REM wails build -s -platform linux/arm -o linux-arm >NUL 2>NUL
set id=copilot-activator
set name=copilot激活器
set version=1.0.0

rd /s /q build\bin
wails build -s -platform windows/amd64 -o windows-amd64 >NUL 2>NUL
wails build -s -platform windows/arm64 -o windows-arm64 >NUL 2>NUL
wails build -s -platform windows/386 -o windows-386 >NUL 2>NUL
echo "生产更新包并上传到cos"
go-selfupdate build\bin\ %version% && coscli sync -r public\ cos://2go/pc/%id%/

@REM rsync 忽略 .gitignore 指定的文件
rsync -avz --exclude-from=.gitignore --delete . bmac:/opt/wailsbuild/