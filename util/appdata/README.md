An app data interface and implementation to facilitate storing app data
in the conventional location of the targeted OS.

where:

- ${env:NAME} represents environment variable resolved at runtime

path templates used are as follows:

global app data:

| OS      | Path template                | Reference                                                                                                                                                                              |
|:--------|:-----------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Windows | ${env:ProgramData}           | [msdn.microsoft.com](https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx) (FOLDERID_ProgramData)                                                           |
| OSX     | /Library/Application Support | [developer.apple.com](https://developer.apple.com/library/content/documentation/General/Conceptual/MOSXAppProgrammingGuide/AppRuntime/AppRuntime.html) (Application Support directory) |
| Linux   | /var/lib                     | [refspecs.linuxfoundation.org](http://refspecs.linuxfoundation.org/FHS_3.0/fhs-3.0.html#varlibVariableStateInformation)                                                                |

per user app data:

| OS      | Path template                       | Reference                                                                                                                                                                              |
|:--------|:------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Windows | ${LOCALAPPDATA}                     | [msdn.microsoft.com](https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx) (FOLDERID_LocalAppData)                                                          |
| OSX     | ${HOME}/Library/Application Support | [developer.apple.com](https://developer.apple.com/library/content/documentation/General/Conceptual/MOSXAppProgrammingGuide/AppRuntime/AppRuntime.html) (Application Support directory) |
| Linux   | ${HOME}                             | conventional                                                                                                                                                                           |


# Credits

Influenced by
[shibukawa/configdir](https://github.com/shibukawa/configdir)
