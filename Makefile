update-assets:
	wails3 update build-assets -name "fapesnap" -config build/config.yml -dir build
	# Wails 3 generates an empty CFBundleExecutable; fill it with the binary name
	sed -i '' '/<key>CFBundleExecutable<\/key>/ { N; s/<string\/>/<string>fapesnap<\/string>/; }' build/darwin/info.plist build/darwin/info.dev.plist build/ios/info.plist build/ios/info.dev.plist
	# If 'appicon' is used, Wails 3 uses the default icon instead of the user-defined one; change to 'icons'
	sed -i '' 's/<string>appicon<\/string>/<string>icons<\/string>/g' build/darwin/info.plist build/darwin/info.dev.plist
generate-icons:
	wails3 generate icons -input build/appicon.png -macassetdir darwin -macfilename build/darwin/icons.icns -windowsfilename build/windows/icon.ico
