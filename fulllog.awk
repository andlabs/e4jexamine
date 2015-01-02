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
	cmd = sprintf("./e4jexamine descdump.%s %s", $1, file)
	while (cmd | getline)
		printf "\t%s\n", $0
	close(cmd)
}
