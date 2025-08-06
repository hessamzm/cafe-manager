[Setup]
AppName=Cafe Manager
AppVersion=1.0
DefaultDirName={pf}\CafeManager
DefaultGroupName=Cafe Manager
UninstallDisplayIcon={app}\cafe-manager.exe
OutputBaseFilename=CafeSetup
Compression=lzma
SolidCompression=yes

[Files]
Source: "cafe-manager.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "icon.png"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Cafe Manager"; Filename: "{app}\cafe-manager.exe"
