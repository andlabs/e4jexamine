# 2 january 2015
# usage: awk -f dumpallfor.awk fulllogfile journalfile blockno
# fulllogfile should be the output of fulllog.awk
# TODO check for errors
BEGIN {
	# TODO make sure this is correct for overwriting ARGC/ARGV
	# TODO check ARGC/ARGV
	journal = ARGV[2]
	blockno = ARGV[3]
	ARGC -= 2
}
/^0/ { next }				# summary line
/commit time/ { next }		# commit record line
# what's left are descriptor lines
$4 != blockno { next }		# wrong block
{ flag = "" }
/escaped/ { flag = "-e" }
{
	cmd = sprintf("../e4jexamine %s blockdump.%s %s > %s", flag, $1, journal, $1)
	system(cmd)
}
