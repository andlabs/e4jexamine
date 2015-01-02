# 2 january 2015
# usage: ./e4jexamine summary file | awk -f fulllog.awk file
# TODO pass options to e4jexamine
# TODO check for errors
BEGIN {
	# TODO make sure this is sufficient to alter ARGC/ARGV
	file = ARGV[1]
	ARGC = 1
}
{ print }
/descriptor/ {
	subcmd("descdump", $1)
}
/commit/ {
	subcmd("commitdump", $1)
}
function subcmd(which, addr,		cmd) {
	cmd = sprintf("./e4jexamine %s.%s %s", which, addr, file)
	while (cmd | getline)
		printf "\t%s\n", $0
	close(cmd)
}
