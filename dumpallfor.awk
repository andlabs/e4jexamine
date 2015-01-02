# 2 january 2015
# usage: awk -f dumpallfor.awk fulllogfile journalfile blockno
# fulllogfile should be the output of fulllog.awk
# TODO check for errors
# TODO find out why this causes heap corruption in mawk
BEGIN {
	# TODO make sure this is correct for overwriting ARGC/ARGV
	# TODO check ARGC/ARGV
	journal = ARGV[2]
	blockno = ARGV[3]
	ARGC -= 2
}
/^0/ { next }				# summary line
# commit record line
/commit time/ {
	times[timei] = $0
	timei++
	next
}
# what's left are descriptor lines
$4 != blockno { next }		# wrong block
{ flag = "" }
/escaped/ { flag = "-e" }
{
	cmd = sprintf("../e4jexamine %s blockdump.%s %s > %s", flag, $1, journal, $1)
	system(cmd)
	filetime[$1] = timei
}
END {
	for (i in filetime)
		printf "%s - %s\n", i, times[filetime[i]]
}
