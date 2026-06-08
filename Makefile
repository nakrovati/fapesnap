update-assets:
	wails3 update build-assets -name "fapesnap" -binaryname "fapesnap" -config build/config.yml
generate-icons:
	wails3 generate icons -input build/appicon.png -macfilename build/darwin/icons.icns -macassetdir build/darwin -windowsfilename build/windows/icon.ico
