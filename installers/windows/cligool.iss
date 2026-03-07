; CliGool 安装程序脚本
; 支持普通用户和特权用户两种模式

; 版本号定义（可通过 /DAppVersion=x.x.x 参数覆盖）
#ifndef AppVersion
  #define AppVersion "1.0.0"
#endif

[Setup]
AppName=CliGool
AppVersion={#AppVersion}
AppPublisher=CliGool
AppPublisherURL=https://github.com/topcheer/cligool
AppSupportURL=https://github.com/topcheer/cligool/issues
AppUpdatesURL=https://github.com/topcheer/cligool/releases

; 默认安装目录（根据权限动态选择）
DefaultDirName={code:GetInstallDir}
DefaultGroupName=CliGool

; 输出文件名
OutputBaseFilename=cligool-setup

; 压缩和安装设置
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
WizardImageFile=installers\windows\cligool-setup.bmp
WizardSmallImageFile=installers\windows\cligool-small.bmp
SetupIconFile=installers\windows\cligool-icon.ico

; 权限相关
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=commandline
UninstallDisplayIcon={app}\cligool-windows-amd64.exe

; 版本信息
VersionInfoVersion={#AppVersion}.0
VersionInfoCompany=CliGool
VersionInfoDescription=CliGool 远程终端安装程序
VersionInfoCopyright=Copyright (C) 2024

[Languages]
Name: "chinesesimp"; MessagesFile: "compiler:Languages\ChineseSimp.isl"
Name: "english"; MessagesFile: "compiler:Languages\English.isl"

[Tasks]
Name: "desktopicon"; Description: "创建桌面快捷方式"; GroupDescription: "附加图标:"
Name: "quicklaunchicon"; Description: "创建快速启动图标"; GroupDescription: "附加图标:"
Name: "contextmenu"; Description: "添加右键菜单""在 CliGool 中打开"""; GroupDescription: "集成选项:"; Flags: unchecked
Name: "createconfig"; Description: "创建默认配置文件"; GroupDescription: "配置选项:"

[Files]
; 主程序（根据架构选择）
Source: "..\bin\cligool-windows-amd64.exe"; DestDir: "{app}"; DestName: "cligool.exe"; Flags: ignoreversion; Check: IsX64
Source: "..\bin\cligool-windows-arm64.exe"; DestDir: "{app}"; DestName: "cligool.exe"; Flags: ignoreversion; Check: IsARM64
Source: "..\bin\cligool-windows-amd64.exe"; DestDir: "{app}"; DestName: "cligool-windows-amd64.exe"; Flags: ignoreversion
Source: "..\bin\cligool-windows-arm64.exe"; DestDir: "{app}"; DestName: "cligool-windows-arm64.exe"; Flags: ignoreversion

; 文档
Source: "..\README.md"; DestDir: "{app}\docs"; Flags: ignoreversion
Source: "..\CONFIG_GUIDE.md"; DestDir: "{app}\docs"; Flags: ignoreversion
Source: "..\CMD_ARGS_USAGE.md"; DestDir: "{app}\docs"; Flags: ignoreversion

[Icons]
; 开始菜单
Name: "{group}\CliGool"; Filename: "{app}\cligool.exe"; Comment: "启动 CliGool 远程终端"
Name: "{group}\卸载 CliGool"; Filename: "{uninstallexe}"

; 桌面快捷方式
Name: "{autodesktop}\CliGool"; Filename: "{app}\cligool.exe"; Tasks: desktopicon; Comment: "启动 CliGool 远程终端"

; 快速启动
Name: "{userappdata}\Microsoft\Internet Explorer\Quick Launch\CliGool"; Filename: "{app}\cligool.exe"; Tasks: quicklaunchicon

[Registry]
; 系统级 PATH（管理员模式）
Root: "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; \
    ValueName: "Path"; \
    ValueData: "{app}"; \
    Check: IsAdminLoggedOn and NeedsAddToPath('{app}')

; 用户级 PATH（普通用户模式）
Root: "HKCU\Environment"; \
    ValueType: expandsz; \
    ValueName: "Path"; \
    ValueData: "{app}"; \
    Check: (not IsAdminLoggedOn) and NeedsAddToPath('{app}')

; 右键菜单集成（仅管理员模式）
Root: "HKCR\Directory\shell\CliGool"; \
    ValueType: string; \
    ValueData: "在 CliGool 中打开"; \
    Flags: uninsdeletekey; \
    Check: IsAdminLoggedOn and IsTaskSelected('contextmenu')
Root: "HKCR\Directory\shell\CliGool"; \
    ValueType: string; \
    ValueName: "Icon"; \
    ValueData: "{app}\cligool.exe,0"; \
    Check: IsAdminLoggedOn and IsTaskSelected('contextmenu')
Root: "HKCR\Directory\shell\CliGool\command"; \
    ValueType: string; \
    ValueData: """{app}\cligool.exe"" -cmd ""%1"""; \
    Check: IsAdminLoggedOn and IsTaskSelected('contextmenu')

; 文件夹右键菜单（仅管理员模式）
Root: "HKCR\Directory\Background\shell\CliGool"; \
    ValueType: string; \
    ValueData: "在 CliGool 中打开此目录"; \
    Flags: uninsdeletekey; \
    Check: IsAdminLoggedOn and IsTaskSelected('contextmenu')
Root: "HKCR\Directory\Background\shell\CliGool\command"; \
    ValueType: string; \
    ValueData: """{app}\cligool.exe"" -cmd ""%V"""; \
    Check: IsAdminLoggedOn and IsTaskSelected('contextmenu')

; 卸载信息
Root: "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\CliGool_is1"; \
    ValueType: string; \
    ValueData: "CliGool"; \
    Flags: uninsdeletekey; \
    Check: IsAdminLoggedOn
Root: "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\CliGool_is1"; \
    ValueType: string; \
    ValueName: "DisplayVersion"; \
    ValueData: "{#AppVersion}"; \
    Flags: uninsdeletekey; \
    Check: IsAdminLoggedOn

[Run]
; 安装完成后运行
Filename: "{app}\cligool.exe"; Description: "启动 CliGool"; Flags: postinstall nowait skipifsilent

[UninstallDelete]
; 删除配置文件选项（可选）
Type: files; Name: "{app}\docs"
Type: files; Name: "{userappdata}\cligool.json"

[Code]
var
  InstallModePage: TOutputOptionWizardPage;
  InstallMode: Integer;

const
  MODE_ADMIN = 1;
  MODE_USER = 2;

function GetInstallDir(Param: String): String;
begin
  if IsAdminLoggedOn then
    Result := '{pf}\CliGool'  // Program Files (管理员)
  else
    Result := '{localappdata}\CliGool';  // 用户目录 (普通用户)
end;

function NeedsAddToPath(PathToAdd: string): Boolean;
var
  OrigPath: string;
begin
  // 检查 PATH 中是否已包含该路径
  if IsAdminLoggedOn then
    RegQueryStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', OrigPath)
  else
    RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath);

  Result := (Pos(OrigPath, PathToAdd) = 0);
end;

procedure AddToPath(PathToAdd: string);
var
  OrigPath: string;
  NewPath: string;
begin
  if IsAdminLoggedOn then begin
    RegQueryStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', OrigPath);
    NewPath := OrigPath + ';' + PathToAdd;
    RegWriteStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', NewPath, REG_EXPAND_SZ);
  end else begin
    RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath);
    NewPath := OrigPath + ';' + PathToAdd;
    RegWriteStringValue(HKCU, 'Environment', 'Path', NewPath, REG_EXPAND_SZ);
  end;
end;

function IsX64: Boolean;
begin
  Result := Is64BitInstallMode and (ProcessorArchitecture = paX64);
end;

function IsARM64: Boolean;
begin
  Result := (ProcessorArchitecture = paARM64);
end;

procedure InitializeWizard;
begin
  // 创建自定义安装模式选择页面
  InstallModePage := CreateOutputOptionPage(wpWelcome,
    '选择安装模式', '请选择 CliGool 的安装模式',
    '安装模式决定了程序的安装位置和可用功能。',
    True, False);

  InstallModePage.AddOption('管理员模式（推荐）',
    '安装到 Program Files，添加到系统 PATH，' +
    '可以使用所有功能（包括右键菜单集成）。',
    MODE_ADMIN);

  InstallModePage.AddOption('用户模式',
    '安装到用户目录，仅添加到用户 PATH，' +
    '功能受限（需要管理员权限的功能不可用）。',
    MODE_USER);

  // 根据当前权限设置默认选项
  if IsAdminLoggedOn then
    InstallModePage.SelectedValue := MODE_ADMIN
  else
    InstallModePage.SelectedValue := MODE_USER;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;

  if CurPageID = InstallModePage.ID then begin
    InstallMode := InstallModePage.SelectedValue;

    // 如果选择了管理员模式但不是管理员，显示警告
    if (InstallMode = MODE_ADMIN) and (not IsAdminLoggedOn) then begin
      if MsgBox('您选择了管理员模式但当前不是管理员身份。' +
                '将以普通用户模式继续安装。' +
                '' + '' + '是否继续？',
                mbConfirmation, MB_YESNO) = IDNO then begin
        Result := False;
      end else begin
        InstallMode := MODE_USER;
      end;
    end;
  end;
end;

procedure CurPageChanged(CurPageID: Integer);
var
  ConfigPath: string;
  ConfigFile: TStringList;
begin
  if CurPageID = wpFinished then begin
    // 如果选择了创建配置文件
    if IsTaskSelected('createconfig') then begin
      ConfigPath := ExpandConstant('{userappdata}\cligool.json');

      if not FileExists(ConfigPath) then begin
        ConfigFile := TStringList.Create;
        try
          ConfigFile.Add('{');
          ConfigFile.Add('  "server": "https://cligool.zty8.cn",');
          ConfigFile.Add('  "cols": 0,');
          ConfigFile.Add('  "rows": 0');
          ConfigFile.Add('}');
          ConfigFile.SaveToFile(ConfigPath);
        finally
          ConfigFile.Free;
        end;
      end;
    end;

    // 添加到 PATH
    if NeedsAddToPath(ExpandConstant('{app}')) then begin
      AddToPath(ExpandConstant('{app}'));
    end;
  end;
end;

function ShouldInstallPage: Boolean;
begin
  // 跳过安装目录选择页面（使用默认位置）
  Result := (InstallMode = MODE_ADMIN);
end;

[Messages]
; 中文
chinesesimp.WelcomeLabel2=本程序将在您的计算机上安装 CliGool 远程终端工具。
chinesesimp.AdminPermissions=如需安装到系统目录或使用完整功能，请以管理员身份运行本程序。
chinesesimp.ClickNext=点击"下一步"继续。

; English
english.WelcomeLabel2=This program will install CliGool Remote Terminal on your computer.
english.AdminPermissions=Please run this program as administrator to install to system directories or use full functionality.
english.ClickNext=Click "Next" to continue.
