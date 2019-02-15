# goconfiger
> .ini config file reader. Automatic hot loading while config file changed.



### Usage

* call `InitConfiger()` in your main func.
* use config file path as parameter to `InitConfiger(filePath)` or `./demo -c scm_config.ini`
* if you send path as  parameter. the program will not load -c in command line.
* use global var `ConfigerMap["key"]` or `ConfigerSection["section"]["key"]` to access the value.
* use global var `ConfigerInfo` to see configer setting.



### How to work?

when you call `InitConfiger()` , it will read config file to global vars. Then use one `go routine` to continue check config file md5, if  md5 is changed, it will automatic reread config file to global vars.