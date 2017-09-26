# Download counter based on logs and regex

#tail -n +0 -f access.log | grep "GET /\(.*\)/\(win32\|osx\|ubuntu.*\)/.*\([0-9]*\.[0-9]*\.[0-9]*\.[0-9]*\).* 200 .*\$" | sed -e "s#^.*\[\([0-9]*.*\)\] .*GET /\(.*\)/\(win32\|osx\|ubuntu.*\)/.*-\([0-9]*\.[0-9]*\.[0-9]*\.[0-9]*\).* 200 .*\$#\1|\2|\3|\4#

	"regex": "^.*\\[\\([0-9]*.*\\)\\] .*GET \/\\(.*\\)\/\\(win32\\|osx\\|ubuntu.*\\)\/.*-\\([0-9]*\\.[0-9]*\\.[0-9]*\\.[0-9]*\\).* 200 .*$",

	"regex": "^.*\\[([0-9]*.*)\\] .*GET \/(.*)/(win32|osx|ubuntu.*)\/.*([0-9]*\\.[0-9]*\\.[0-9]*).* 200 .*$",
	"fields": [ "time", "product", "platform", "version" ]

	"regex": "^(.*)$",
	"fields": [ "field" ]

	"regex": "^.*\\[([0-9]*.*)\\] .*GET (\/(.*)/(win32|osx|ubuntu.*)\/[^0-9]*([0-9]*\\.[0-9]*\\.[0-9]*\\.[0-9]*).*[^s][^i][^g]) HTTP.* 200 .*$",
	"fields": [ "time", "file", "product", "platform", "version" ]
